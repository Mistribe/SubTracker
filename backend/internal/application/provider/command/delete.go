package command

import (
	"context"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/provider"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type DeleteCommand struct {
	ProviderId uuid.UUID
}

func NewDeleteCommand(providerId uuid.UUID) DeleteCommand {
	return DeleteCommand{
		ProviderId: providerId,
	}
}

type DeleteCommandHandler struct {
	repository provider.Repository
}

func NewDeleteCommandHandler(repository provider.Repository) *DeleteCommandHandler {
	return &DeleteCommandHandler{repository: repository}
}

func (h DeleteCommandHandler) Handle(ctx context.Context, cmd DeleteCommand) result.Result[bool] {
	ok, err := h.repository.Delete(ctx, cmd.ProviderId)
	if err != nil {
		return result.Fail[bool](err)
	}

	return result.Success(ok)
}
