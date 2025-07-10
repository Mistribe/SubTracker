package command

import (
	"context"

	"github.com/oleexo/subtracker/internal/application/core/result"
	"github.com/oleexo/subtracker/internal/domain/label"
)

type CreateLabelCommand struct {
	Label label.Label
}

type CreateLabelCommandHandler struct {
	repository label.Repository
}

func NewCreateLabelCommandHandler(repository label.Repository) *CreateLabelCommandHandler {
	return &CreateLabelCommandHandler{repository: repository}
}

func (h CreateLabelCommandHandler) Handle(ctx context.Context, command CreateLabelCommand) result.Result[label.Label] {
	existingLabel, err := h.repository.Get(ctx, command.Label.Id())
	if err != nil {
		return result.Fail[label.Label](err)
	}
	if existingLabel != nil {
		return result.Fail[label.Label](label.ErrLabelAlreadyExists)
	}

	if err := command.Label.Validate(); err != nil {
		return result.Fail[label.Label](err)
	}
	if err := h.repository.Save(ctx, command.Label); err != nil {
		return result.Fail[label.Label](err)
	}

	return result.Success(command.Label)
}
