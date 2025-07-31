package command

import (
	"context"

	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/internal/domain/family"
	"github.com/oleexo/subtracker/internal/domain/label"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type CreateLabelCommand struct {
	Label label.Label
}

type CreateLabelCommandHandler struct {
	labelRepository  label.Repository
	familyRepository family.Repository
}

func NewCreateLabelCommandHandler(labelRepository label.Repository,
	familyRepository family.Repository) *CreateLabelCommandHandler {
	return &CreateLabelCommandHandler{
		labelRepository:  labelRepository,
		familyRepository: familyRepository,
	}
}

func (h CreateLabelCommandHandler) Handle(ctx context.Context, command CreateLabelCommand) result.Result[label.Label] {
	existingLabel, err := h.labelRepository.GetById(ctx, command.Label.Id())
	if err != nil {
		return result.Fail[label.Label](err)
	}

	err = auth.EnsureOwnership(ctx, h.familyRepository, command.Label.Owner())
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
		"",
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
