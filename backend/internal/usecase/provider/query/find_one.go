package query

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/provider"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type FindOneQuery struct {
	ProviderID  *types.ProviderID
	ProviderKey *string
}

type FindOneQueryHandler struct {
	providerRepository ports.ProviderRepository
	authService        ports.Authentication
}

func NewFindOneQueryHandler(
	providerRepository ports.ProviderRepository,
	authService ports.Authentication) *FindOneQueryHandler {
	return &FindOneQueryHandler{
		providerRepository: providerRepository,
		authService:        authService,
	}
}

func (h FindOneQueryHandler) Handle(ctx context.Context, query FindOneQuery) result.Result[provider.Provider] {
	connectedAccount := h.authService.MustGetConnectedAccount(ctx)
	if query.ProviderID != nil {
		prvdr, err := h.providerRepository.GetByIdForUser(ctx, connectedAccount.UserID(), *query.ProviderID)
		if err != nil {
			return result.Fail[provider.Provider](err)
		}

		if prvdr == nil {
			return result.Fail[provider.Provider](provider.ErrProviderNotFound)
		}
		return result.Success(prvdr)
	} else if query.ProviderKey != nil {
		prvdr, err := h.providerRepository.GetByProviderKeyForUser(ctx, connectedAccount.UserID(), *query.ProviderKey)
		if err != nil {
			return result.Fail[provider.Provider](err)
		}

		if prvdr == nil {
			return result.Fail[provider.Provider](provider.ErrProviderNotFound)
		}
		return result.Success(prvdr)
	}

	return result.Fail[provider.Provider](provider.ErrProviderNotFound)
}
