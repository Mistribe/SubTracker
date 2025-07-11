package persistence

import (
	"context"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/application/core/option"
	"github.com/oleexo/subtracker/internal/domain/subscription"
)

type SubscriptionRepository struct {
}

func NewSubscriptionRepository() *SubscriptionRepository {
	return &SubscriptionRepository{}
}

func (r SubscriptionRepository) Get(ctx context.Context, id uuid.UUID) (
	option.Option[subscription.Subscription],
	error) {
	//TODO implement me
	panic("implement me")
}

func (r SubscriptionRepository) GetAll(ctx context.Context) ([]subscription.Subscription, error) {
	//TODO implement me
	panic("implement me")
}

func (r SubscriptionRepository) Save(ctx context.Context, subscription subscription.Subscription) error {
	//TODO implement me
	panic("implement me")
}

func (r SubscriptionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}
