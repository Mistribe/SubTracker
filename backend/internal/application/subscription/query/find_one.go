package query

import (
	"context"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/subscription"
	"github.com/oleexo/subtracker/pkg/langext/result"
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

func NewFindOneQueryHandler(repository subscription.Repository) *FindOneQueryHandler {
	return &FindOneQueryHandler{repository: repository}
}

func (h FindOneQueryHandler) Handle(ctx context.Context, query FindOneQuery) result.Result[subscription.Subscription] {
	sub, err := h.repository.GetById(ctx, query.id)
	if err != nil {
		return result.Fail[subscription.Subscription](err)
	}

	if sub == nil {
		return result.Fail[subscription.Subscription](subscription.ErrSubscriptionNotFound)
	}

	return result.Success(sub)
}
