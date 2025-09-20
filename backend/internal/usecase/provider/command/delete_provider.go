package command

import (
	"context"

	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/domain/provider"
	"github.com/mistribe/subtracker/internal/domain/user"
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

type DeleteProviderCommandHandler struct {
	providerRepository ports.ProviderRepository
	authorization      ports.Authorization
}

func NewDeleteProviderCommandHandler(
	providerRepository ports.ProviderRepository,
	authorization ports.Authorization) *DeleteProviderCommandHandler {
	return &DeleteProviderCommandHandler{
		providerRepository: providerRepository,
		authorization:      authorization,
	}
}

func (h DeleteProviderCommandHandler) Handle(ctx context.Context, cmd DeleteProviderCommand) result.Result[bool] {
	prov, err := h.providerRepository.GetById(ctx, cmd.ProviderId)
	if err != nil {
		return result.Fail[bool](err)
	}

	if prov == nil {
		return result.Fail[bool](provider.ErrProviderNotFound)
	}

	if err = h.authorization.Can(ctx, user.PermissionDelete).For(prov); err != nil {
		return result.Fail[bool](err)
	}

	ok, err := h.providerRepository.Delete(ctx, cmd.ProviderId)
	if err != nil {
		return result.Fail[bool](err)
	}

	return result.Success(ok)
}
