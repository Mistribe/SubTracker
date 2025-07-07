package query

import (
    "context"
    "github.com/google/uuid"
    "github.com/oleexo/subtracker/internal/application/core/option"
    "github.com/oleexo/subtracker/internal/application/core/result"
    "github.com/oleexo/subtracker/internal/domain/subscription"
)

type FindQuery struct {
    id uuid.UUID
}

func NewFindQuery(id uuid.UUID) FindQuery {
    return FindQuery{id: id}
}

type FindQueryHandler struct {
    respository subscription.Repository
}

func NewFindQueryHandler(respository subscription.Repository) *FindQueryHandler {
    return &FindQueryHandler{respository: respository}
}

func (h FindQueryHandler) Handle(ctx context.Context, query FindQuery) result.Result[subscription.Subscription] {
    subOpt, err := h.respository.Get(ctx, query.id)
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
