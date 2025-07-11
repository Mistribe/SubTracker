package query

import (
	"context"

	"github.com/oleexo/subtracker/internal/application/core/result"
	"github.com/oleexo/subtracker/internal/domain/label"
)

type FindAllQuery struct {
	WithDefault bool
}

func NewFindAllQuery(withDefault bool) FindAllQuery {
	return FindAllQuery{
		WithDefault: withDefault,
	}
}

type FindAllQueryHandler struct {
	repository label.Repository
}

func NewFindAllQueryHandler(repository label.Repository) *FindAllQueryHandler {
	return &FindAllQueryHandler{repository: repository}
}

func (h FindAllQueryHandler) Handle(ctx context.Context, query FindAllQuery) result.Result[[]label.Label] {
	lbs, err := h.repository.GetAll(ctx, query.WithDefault)
	if err != nil {
		return result.Fail[[]label.Label](err)
	}

	return result.Success(lbs)
}
