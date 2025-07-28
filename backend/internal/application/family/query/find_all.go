package query

import (
	"context"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/domain/entity"
	"github.com/oleexo/subtracker/internal/domain/family"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type FindAllQuery struct {
	Limit  int
	Offset int
}

func NewFindAllQuery(size int, page int) FindAllQuery {
	return FindAllQuery{Limit: size, Offset: page}
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
	parameters := entity.NewQueryParameters(query.Limit, query.Offset)
	families, err := h.repository.GetAll(ctx, parameters)
	if err != nil {
		return result.Fail[core.PaginatedResponse[family.Family]](err)
	}

	total, err := h.repository.GetAllCount(ctx)
	if err != nil {
		return result.Fail[core.PaginatedResponse[family.Family]](err)
	}
	return result.Success(core.NewPaginatedResponse(families, total))
}
