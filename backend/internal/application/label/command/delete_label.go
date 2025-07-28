package command

import (
	"context"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/label"
	"github.com/oleexo/subtracker/pkg/langext/result"
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
	existingLabel, err := h.repository.GetById(ctx, command.Id)
	if err != nil {
		return result.Fail[bool](err)
	}
	if existingLabel == nil {
		return result.Fail[bool](label.ErrLabelNotFound)
	}

	ok, err := h.repository.Delete(ctx, command.Id)
	if err != nil {
		return result.Fail[bool](err)
	}
	if !ok {
		return result.Fail[bool](label.ErrLabelNotFound)
	}

	return result.Success(true)
}
