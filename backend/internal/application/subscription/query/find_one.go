package query

import (
	"context"

	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/domain/auth"
	"github.com/mistribe/subtracker/internal/domain/subscription"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type FindOneQuery struct {
	id uuid.UUID
}

func NewFindOneQuery(id uuid.UUID) FindOneQuery {
	return FindOneQuery{id: id}
}

type FindOneQueryHandler struct {
	subscriptionRepository subscription.Repository
	authService            auth.Service
}

func NewFindOneQueryHandler(
	subscriptionRepository subscription.Repository,
	authService auth.Service) *FindOneQueryHandler {
	return &FindOneQueryHandler{
		subscriptionRepository: subscriptionRepository,
		authService:            authService,
	}
}

func (h FindOneQueryHandler) Handle(ctx context.Context, query FindOneQuery) result.Result[subscription.Subscription] {
	userId := h.authService.MustGetUserId(ctx)
	sub, err := h.subscriptionRepository.GetByIdForUser(ctx, userId, query.id)
	if err != nil {
		return result.Fail[subscription.Subscription](err)
	}

	if sub == nil {
		return result.Fail[subscription.Subscription](subscription.ErrSubscriptionNotFound)
	}

	return result.Success(sub)
}
