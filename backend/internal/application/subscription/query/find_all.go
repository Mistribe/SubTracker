package query

import (
	"context"
	"github.com/oleexo/subtracker/internal/application/core/result"
	"github.com/oleexo/subtracker/internal/domain/subscription"
)

type FindAllQuery struct {
}

func NewFindAllQuery() FindAllQuery {
	return FindAllQuery{}
}

type FindAllQueryHandler struct {
	repository subscription.Repository
}

func NewFindAllQueryHandler(repository subscription.Repository) *FindAllQueryHandler {
	return &FindAllQueryHandler{repository: repository}
}

func (h FindAllQueryHandler) Handle(ctx context.Context, query FindAllQuery) result.Result[[]subscription.Repository] {
	
}
