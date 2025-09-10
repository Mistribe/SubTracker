package command

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/domain/label"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type CreateLabelCommand struct {
	Label label.Label
}

type CreateLabelCommandHandler struct {
	labelRepository  ports.LabelRepository
	familyRepository ports.FamilyRepository
	authService      ports.AuthenticationService
}

func NewCreateLabelCommandHandler(labelRepository ports.LabelRepository,
	familyRepository ports.FamilyRepository,
	authService ports.AuthenticationService) *CreateLabelCommandHandler {
	return &CreateLabelCommandHandler{
		labelRepository:  labelRepository,
		familyRepository: familyRepository,
		authService:      authService,
	}
}

func (h CreateLabelCommandHandler) Handle(ctx context.Context, command CreateLabelCommand) result.Result[label.Label] {
	existingLabel, err := h.labelRepository.GetById(ctx, command.Label.Id())
	if err != nil {
		return result.Fail[label.Label](err)
	}

	ok, err := h.authService.IsOwner(ctx, command.Label.Owner())
	if err != nil {
		return result.Fail[label.Label](err)
	}
	if !ok {
		return result.Fail[label.Label](family.ErrFamilyNotFound)
	}
	if existingLabel != nil {
		if existingLabel.Equal(command.Label) {
			return result.Success(command.Label)
		}
		return result.Fail[label.Label](label.ErrLabelAlreadyExists)
	}

	return h.createLabel(ctx, command)
}

func (h CreateLabelCommandHandler) createLabel(
	ctx context.Context, command CreateLabelCommand) result.Result[label.Label] {
	lbl := label.NewLabel(command.Label.Id(),
		command.Label.Owner(),
		command.Label.Name(),
		nil,
		command.Label.Color(),
		command.Label.CreatedAt(),
		command.Label.CreatedAt(),
	)

	if err := lbl.GetValidationErrors(); err != nil {
		return result.Fail[label.Label](err)
	}

	if err := h.labelRepository.Save(ctx, lbl); err != nil {
		return result.Fail[label.Label](err)
	}
	return result.Success(lbl)
}
