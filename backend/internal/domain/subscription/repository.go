package subscription

import (
    "context"
    "github.com/google/uuid"
)

type Repository interface {
    Get(ctx context.Context, id uuid.UUID) (*Subscription, error)
    GetAll(ctx context.Context) ([]Subscription, error)
    Save(ctx context.Context, subscription Subscription) error
}
