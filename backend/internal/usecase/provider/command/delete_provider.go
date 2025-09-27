package command

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/authorization"
	"github.com/mistribe/subtracker/internal/domain/provider"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type DeleteProviderCommand struct {
	ProviderID types.ProviderID
}

func NewDeleteProviderCommand(providerID types.ProviderID) DeleteProviderCommand {
	return DeleteProviderCommand{
		ProviderID: providerID,
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
	prov, err := h.providerRepository.GetById(ctx, cmd.ProviderID)
	if err != nil {
		return result.Fail[bool](err)
	}

	if prov == nil {
		return result.Fail[bool](provider.ErrProviderNotFound)
	}

	if err = h.authorization.Can(ctx, authorization.PermissionDelete).For(prov); err != nil {
		return result.Fail[bool](err)
	}

	ok, err := h.providerRepository.Delete(ctx, cmd.ProviderID)
	if err != nil {
		return result.Fail[bool](err)
	}

	return result.Success(ok)
}
