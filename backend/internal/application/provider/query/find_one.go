package query

import (
	"context"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/provider"
	"github.com/oleexo/subtracker/pkg/langext/result"
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
	repository provider.Repository
}

func NewFindOneQueryHandler(repository provider.Repository) *FindOneQueryHandler {
	return &FindOneQueryHandler{repository: repository}
}

func (h FindOneQueryHandler) Handle(ctx context.Context, query FindOneQuery) result.Result[provider.Provider] {
	prvdr, err := h.repository.GetById(ctx, query.ProviderId)
	if err != nil {
		return result.Fail[provider.Provider](err)
	}

	if prvdr == nil {
		return result.Fail[provider.Provider](provider.ErrProviderNotFound)
	}
	return result.Success(prvdr)
}
