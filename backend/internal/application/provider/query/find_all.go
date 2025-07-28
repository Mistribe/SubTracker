package query

import (
	"context"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/domain/entity"
	"github.com/oleexo/subtracker/internal/domain/provider"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type FindAllQuery struct {
	Limit  int
	Offset int
}

func NewFindAllQuery(limit, offset int) FindAllQuery {
	return FindAllQuery{
		Limit:  limit,
		Offset: offset,
	}
}

type FindAllQueryHandler struct {
	repository provider.Repository
}

func NewFindAllQueryHandler(repository provider.Repository) *FindAllQueryHandler {
	return &FindAllQueryHandler{repository: repository}
}

func (h FindAllQueryHandler) Handle(
	ctx context.Context,
	query FindAllQuery) result.Result[core.PaginatedResponse[provider.Provider]] {
	providers, err := h.repository.GetAll(ctx, entity.NewQueryParameters(query.Limit, query.Offset))
	if err != nil {
		return result.Fail[core.PaginatedResponse[provider.Provider]](err)
	}

	total, err := h.repository.GetAllCount(ctx)
	if err != nil {
		return result.Fail[core.PaginatedResponse[provider.Provider]](err)
	}
	return result.Success(core.NewPaginatedResponse(providers, total))
}
