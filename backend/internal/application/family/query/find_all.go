package query

import (
	"context"

	"github.com/oleexo/subtracker/internal/domain/family"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type FindAllQuery struct {
}

type FindAllQueryHandler struct {
	repository family.Repository
}

func NewFindAllQueryHandler(repository family.Repository) *FindAllQueryHandler {
	return &FindAllQueryHandler{repository: repository}
}

func (h FindAllQueryHandler) Handle(ctx context.Context, query FindAllQuery) result.Result[[]family.Family] {
	families, err := h.repository.GetAll(ctx)
	if err != nil {
		return result.Fail[[]family.Family](err)
	}

	return result.Success(families)
}
