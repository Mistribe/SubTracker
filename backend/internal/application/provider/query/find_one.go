package query

import (
	"context"

	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/domain/auth"
	"github.com/mistribe/subtracker/internal/domain/provider"
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
	providerRepository provider.Repository
	authService        auth.Service
}

func NewFindOneQueryHandler(
	providerRepository provider.Repository,
	authService auth.Service) *FindOneQueryHandler {
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
