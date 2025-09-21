package command

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/label"
	"github.com/mistribe/subtracker/internal/domain/user"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type CreateLabelCommand struct {
	Label label.Label
}

type CreateLabelCommandHandler struct {
	labelRepository  ports.LabelRepository
	familyRepository ports.FamilyRepository
	authorization    ports.Authorization
}

func NewCreateLabelCommandHandler(labelRepository ports.LabelRepository,
	familyRepository ports.FamilyRepository,
	authorization ports.Authorization) *CreateLabelCommandHandler {
	return &CreateLabelCommandHandler{
		labelRepository:  labelRepository,
		familyRepository: familyRepository,
		authorization:    authorization,
	}
}

func (h CreateLabelCommandHandler) Handle(ctx context.Context, command CreateLabelCommand) result.Result[label.Label] {
	if err := h.authorization.Can(ctx, user.PermissionWrite).For(command.Label); err != nil {
		return result.Fail[label.Label](err)
	}

	// Ensure the operation respects limits regardless of label existence state
	if err := h.authorization.EnsureWithinLimit(ctx, user.FeatureCustomLabels, 1); err != nil {
		return result.Fail[label.Label](err)
	}

	existingLabel, err := h.labelRepository.GetById(ctx, command.Label.Id())
	if err != nil {
		return result.Fail[label.Label](err)
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
