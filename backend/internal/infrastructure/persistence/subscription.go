package persistence

import (
	"context"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/application/core/option"
	"github.com/oleexo/subtracker/internal/domain/subscription"
)

type SubscriptionRepository struct {
	subscriptions map[uuid.UUID]subscription.Subscription
}

func NewSubscriptionRepository() *SubscriptionRepository {
	return &SubscriptionRepository{
		subscriptions: make(map[uuid.UUID]subscription.Subscription),
	}
}

func (r SubscriptionRepository) Get(ctx context.Context, id uuid.UUID) (
	option.Option[subscription.Subscription],
	error) {
	if sub, ok := r.subscriptions[id]; ok {
		return option.Some(sub), nil
	}
	return option.None[subscription.Subscription](), nil
}

func (r SubscriptionRepository) GetAll(ctx context.Context) ([]subscription.Subscription, error) {
	subscriptions := make([]subscription.Subscription, 0, len(r.subscriptions))
	for _, sub := range r.subscriptions {
		subscriptions = append(subscriptions, sub)
	}
	return subscriptions, nil
}

func (r SubscriptionRepository) Save(ctx context.Context, subscription subscription.Subscription) error {
	r.subscriptions[subscription.Id()] = subscription
	return nil
}

func (r SubscriptionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	delete(r.subscriptions, id)
	return nil
}
