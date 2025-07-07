package persistence

import (
    "context"
    "github.com/google/uuid"
    "github.com/oleexo/subtracker/internal/domain/subscription"
)

type SubsriptionRepository struct {
}

func NewSubsriptionRepository() subscription.Repository {
    return &SubsriptionRepository{}
}

func (s SubsriptionRepository) Get(ctx context.Context, id uuid.UUID) (*subscription.Subscription, error) {
    //TODO implement me
    panic("implement me")
}

func (s SubsriptionRepository) GetAll(ctx context.Context) ([]subscription.Subscription, error) {
    //TODO implement me
    panic("implement me")
}

func (s SubsriptionRepository) Save(ctx context.Context, subscription subscription.Subscription) error {
    //TODO implement me
    panic("implement me")
}
