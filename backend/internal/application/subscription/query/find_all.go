package query

import (
	"context"

	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/internal/domain/entity"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/domain/subscription"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type FindAllQuery struct {
	Limit  int32
	Offset int32
}

func NewFindAllQuery(size, page int32) FindAllQuery {
	return FindAllQuery{
		Limit:  size,
		Offset: page,
	}
}

type FindAllQueryHandler struct {
	subscriptionRepository subscription.Repository
	authService            auth.Service
}

func NewFindAllQueryHandler(
	subscriptionRepository subscription.Repository,
	authService auth.Service) *FindAllQueryHandler {
	return &FindAllQueryHandler{
		subscriptionRepository: subscriptionRepository,
		authService:            authService,
	}
}

func (h FindAllQueryHandler) Handle(
	ctx context.Context,
	query FindAllQuery) result.Result[core.PaginatedResponse[subscription.Subscription]] {
	userId := h.authService.MustGetUserId(ctx)
	parameters := entity.NewQueryParameters(query.Limit, query.Offset)
	subs, count, err := h.subscriptionRepository.GetAllForUser(ctx, userId, parameters)
	if err != nil {
		return result.Fail[core.PaginatedResponse[subscription.Subscription]](err)
	}

	return result.Success(core.NewPaginatedResponse(subs, count))
}
