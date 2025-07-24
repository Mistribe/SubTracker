package query

import (
	"context"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/domain/family"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type FindAllQuery struct {
	Size int
	Page int
}

func NewFindAllQuery(size int, page int) FindAllQuery {
	return FindAllQuery{Size: size, Page: page}
}

type FindAllQueryHandler struct {
	repository family.Repository
}

func NewFindAllQueryHandler(repository family.Repository) *FindAllQueryHandler {
	return &FindAllQueryHandler{repository: repository}
}

func (h FindAllQueryHandler) Handle(
	ctx context.Context,
	query FindAllQuery) result.Result[core.PaginatedResponse[family.Family]] {
	families, err := h.repository.GetAll(ctx, query.Size, query.Page)
	if err != nil {
		return result.Fail[core.PaginatedResponse[family.Family]](err)
	}

	total, err := h.repository.GetAllCount(ctx)
	if err != nil {
		return result.Fail[core.PaginatedResponse[family.Family]](err)
	}
	return result.Success(core.NewPaginatedResponse(families, total))
}
