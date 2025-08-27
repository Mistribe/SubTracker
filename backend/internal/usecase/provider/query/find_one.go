package query

import (
	"context"

	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/domain/provider"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type FindOneQuery struct {
	ProviderId uuid.UUID
}

func NewFindOneQuery(providerId uuid.UUID) FindOneQuery {
	return FindOneQuery{
		ProviderId: providerId,
	}
}

type FindOneQueryHandler struct {
	providerRepository ports.ProviderRepository
	authService        ports.AuthService
}

func NewFindOneQueryHandler(
	providerRepository ports.ProviderRepository,
	authService ports.AuthService) *FindOneQueryHandler {
	return &FindOneQueryHandler{
		providerRepository: providerRepository,
		authService:        authService,
	}
}

func (h FindOneQueryHandler) Handle(ctx context.Context, query FindOneQuery) result.Result[provider.Provider] {
	userId := h.authService.MustGetUserId(ctx)
	prvdr, err := h.providerRepository.GetByIdForUser(ctx, userId, query.ProviderId)
	if err != nil {
		return result.Fail[provider.Provider](err)
	}

	if prvdr == nil {
		return result.Fail[provider.Provider](provider.ErrProviderNotFound)
	}
	return result.Success(prvdr)
}
