package query

import (
	"context"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/provider"
	"github.com/oleexo/subtracker/pkg/langext/option"
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
	providerOpt, err := h.repository.GetById(ctx, query.ProviderId)
	if err != nil {
		return result.Fail[provider.Provider](err)
	}

	return option.Match(providerOpt, func(in provider.Provider) result.Result[provider.Provider] {
		return result.Success(in)
	}, func() result.Result[provider.Provider] {
		return result.Fail[provider.Provider](provider.ErrProviderNotFound)
	})
}
