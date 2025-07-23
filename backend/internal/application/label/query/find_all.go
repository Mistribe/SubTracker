package query

import (
	"context"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/domain/label"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type FindAllQuery struct {
	WithDefault bool
	Size        int
	Page        int
}

func NewFindAllQuery(withDefault bool, size int, page int) FindAllQuery {
	if size < 0 {
		size = 10
	}
	if page < 1 {
		page = 1
	}
	return FindAllQuery{
		WithDefault: withDefault,
		Size:        size,
		Page:        page,
	}
}

type FindAllQueryHandler struct {
	repository label.Repository
}

func NewFindAllQueryHandler(repository label.Repository) *FindAllQueryHandler {
	return &FindAllQueryHandler{repository: repository}
}

func (h FindAllQueryHandler) Handle(ctx context.Context,
	query FindAllQuery) result.Result[core.PaginatedResponse[label.Label]] {
	lbs, err := h.repository.GetAll(ctx, query.Size, query.Page, query.WithDefault)
	if err != nil {
		return result.Fail[core.PaginatedResponse[label.Label]](err)
	}
	total, err := h.repository.GetAllCount(ctx, query.WithDefault)
	if err != nil {
		return result.Fail[core.PaginatedResponse[label.Label]](err)
	}
	return result.Success(core.NewPaginatedResponse(lbs, total))
}
