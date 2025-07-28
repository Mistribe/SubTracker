package query

import (
	"context"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/domain/label"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type FindAllQuery struct {
	WithDefault bool
	Limit       int
	Offset      int
}

func NewFindAllQuery(withDefault bool, size int, offset int) FindAllQuery {
	return FindAllQuery{
		WithDefault: withDefault,
		Limit:       size,
		Offset:      offset,
	}
}

type FindAllQueryHandler struct {
	repository label.Repository
}

func NewFindAllQueryHandler(repository label.Repository) *FindAllQueryHandler {
	return &FindAllQueryHandler{repository: repository}
}

func (h FindAllQueryHandler) Handle(
	ctx context.Context,
	query FindAllQuery) result.Result[core.PaginatedResponse[label.Label]] {
	params := label.NewQueryParameters(query.Limit, query.Offset, query.WithDefault)
	lbs, err := h.repository.GetAll(ctx, params)
	if err != nil {
		return result.Fail[core.PaginatedResponse[label.Label]](err)
	}
	total, err := h.repository.GetAllCount(ctx, params)
	if err != nil {
		return result.Fail[core.PaginatedResponse[label.Label]](err)
	}
	return result.Success(core.NewPaginatedResponse(lbs, total))
}
