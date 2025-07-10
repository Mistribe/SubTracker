package query

import (
	"context"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/application/core/option"
	"github.com/oleexo/subtracker/internal/application/core/result"
	"github.com/oleexo/subtracker/internal/domain/subscription"
)

type FindOneQuery struct {
	id uuid.UUID
}

func NewFindOneQuery(id uuid.UUID) FindOneQuery {
	return FindOneQuery{id: id}
}

type FindOneQueryHandler struct {
	repository subscription.Repository
}

func NewFindOneQueryHandler(respository subscription.Repository) *FindOneQueryHandler {
	return &FindOneQueryHandler{repository: respository}
}

func (h FindOneQueryHandler) Handle(ctx context.Context, query FindOneQuery) result.Result[subscription.Subscription] {
	subOpt, err := h.repository.Get(ctx, query.id)
	if err != nil {
		return result.Fail[subscription.Subscription](err)
	}

	return option.Match[subscription.Subscription, result.Result[subscription.Subscription]](subOpt,
		func(value subscription.Subscription) result.Result[subscription.Subscription] {
			return result.Success(value)
		},
		func() result.Result[subscription.Subscription] {
			return result.Fail[subscription.Subscription](subscription.ErrSubscriptionNotFound)
		})
}
