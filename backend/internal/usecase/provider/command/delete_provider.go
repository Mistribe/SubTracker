package command

import (
	"context"

	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type DeleteProviderCommand struct {
	ProviderId uuid.UUID
}

func NewDeleteProviderCommand(providerId uuid.UUID) DeleteProviderCommand {
	return DeleteProviderCommand{
		ProviderId: providerId,
	}
}

type DeleteCommandHandler struct {
	repository ports.ProviderRepository
}

func NewDeleteProviderCommandHandler(repository ports.ProviderRepository) *DeleteCommandHandler {
	return &DeleteCommandHandler{repository: repository}
}

func (h DeleteCommandHandler) Handle(ctx context.Context, cmd DeleteProviderCommand) result.Result[bool] {
	ok, err := h.repository.Delete(ctx, cmd.ProviderId)
	if err != nil {
		return result.Fail[bool](err)
	}

	return result.Success(ok)
}
