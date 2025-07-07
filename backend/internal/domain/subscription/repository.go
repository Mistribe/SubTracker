package subscription

import (
    "context"
    "github.com/google/uuid"
    "github.com/oleexo/subtracker/internal/application/core/option"
)

type Repository interface {
    Get(ctx context.Context, id uuid.UUID) (option.Option[Subscription], error)
    GetAll(ctx context.Context) ([]Subscription, error)
    Save(ctx context.Context, subscription Subscription) error
}
