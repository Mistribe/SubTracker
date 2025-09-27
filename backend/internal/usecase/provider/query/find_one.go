package query

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/provider"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type FindOneQuery struct {
	ProviderID types.ProviderID
}

func NewFindOneQuery(providerID types.ProviderID) FindOneQuery {
	return FindOneQuery{
		ProviderID: providerID,
	}
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
	prvdr, err := h.providerRepository.GetByIdForUser(ctx, connectedAccount.UserID(), query.ProviderID)
	if err != nil {
		return result.Fail[provider.Provider](err)
	}

	if prvdr == nil {
		return result.Fail[provider.Provider](provider.ErrProviderNotFound)
	}
	return result.Success(prvdr)
}
