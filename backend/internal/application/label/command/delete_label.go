package command

import (
	"context"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/application/core/result"
	"github.com/oleexo/subtracker/internal/domain/label"
)

type DeleteLabelCommand struct {
	Id uuid.UUID
}

type DeleteLabelCommandHandler struct {
	repository label.Repository
}

func NewDeleteLabelCommandHandler(repository label.Repository) *DeleteLabelCommandHandler {
	return &DeleteLabelCommandHandler{repository: repository}
}

func (h DeleteLabelCommandHandler) Handle(ctx context.Context, command DeleteLabelCommand) result.Result[bool] {
	existingLabel, err := h.repository.Get(ctx, command.Id)
	if err != nil {
		return result.Fail[bool](err)
	}
	if existingLabel == nil {
		return result.Fail[bool](label.ErrLabelNotFound)
	}

	if err := h.repository.Delete(ctx, command.Id); err != nil {
		return result.Fail[bool](err)
	}

	return result.Success(true)
}
