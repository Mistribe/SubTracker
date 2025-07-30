package query

import (
	"context"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/internal/domain/label"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type FindAllQuery struct {
	Owners []auth.OwnerType
	Limit  int
	Offset int
}

func NewFindAllQuery(owners []auth.OwnerType, size int, offset int) FindAllQuery {
	return FindAllQuery{
		Owners: owners,
		Limit:  size,
		Offset: offset,
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
	params := label.NewQueryParameters(query.Limit, query.Offset, query.Owners)
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
