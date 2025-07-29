package query

import (
	"context"
	"github.com/oleexo/subtracker/internal/domain/entity"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/domain/subscription"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type FindAllQuery struct {
	Limit  int
	Offset int
}

func NewFindAllQuery(size, page int) FindAllQuery {
	return FindAllQuery{
		Limit:  size,
		Offset: page,
	}
}

type FindAllQueryHandler struct {
	repository subscription.Repository
}

func NewFindAllQueryHandler(repository subscription.Repository) *FindAllQueryHandler {
	return &FindAllQueryHandler{repository: repository}
}

func (h FindAllQueryHandler) Handle(
	ctx context.Context,
	query FindAllQuery) result.Result[core.PaginatedResponse[subscription.Subscription]] {
	parameters := entity.NewQueryParameters(query.Limit, query.Offset)
	subs, err := h.repository.GetAll(ctx, parameters)
	if err != nil {
		return result.Fail[core.PaginatedResponse[subscription.Subscription]](err)
	}

	total, err := h.repository.GetAllCount(ctx)
	if err != nil {
		return result.Fail[core.PaginatedResponse[subscription.Subscription]](err)
	}
	return result.Success(core.NewPaginatedResponse(subs, total))
}
