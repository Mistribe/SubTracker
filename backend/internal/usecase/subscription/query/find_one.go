package query

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/subscription"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type FindOneQuery struct {
	SubscriptionID types.SubscriptionID
}

func NewFindOneQuery(subscriptionID types.SubscriptionID) FindOneQuery {
	return FindOneQuery{SubscriptionID: subscriptionID}
}

type FindOneQueryHandler struct {
	subscriptionRepository ports.SubscriptionRepository
	authService            ports.Authentication
}

func NewFindOneQueryHandler(
	subscriptionRepository ports.SubscriptionRepository,
	authService ports.Authentication) *FindOneQueryHandler {
	return &FindOneQueryHandler{
		subscriptionRepository: subscriptionRepository,
		authService:            authService,
	}
}

func (h FindOneQueryHandler) Handle(ctx context.Context, query FindOneQuery) result.Result[subscription.Subscription] {
	connectedAccount := h.authService.MustGetConnectedAccount(ctx)
	sub, err := h.subscriptionRepository.GetByIdForUser(ctx, connectedAccount.UserID(), query.SubscriptionID)
	if err != nil {
		return result.Fail[subscription.Subscription](err)
	}

	if sub == nil {
		return result.Fail[subscription.Subscription](subscription.ErrSubscriptionNotFound)
	}

	return result.Success(sub)
}
