package subscription

import (
	"context"
	"iter"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/entity"
)

type Repository interface {
	entity.Repository[Subscription]

	GetByIdForUser(ctx context.Context, userId string, id uuid.UUID) (Subscription, error)
	GetAll(ctx context.Context, parameters entity.QueryParameters) ([]Subscription, int64, error)
	GetAllForUser(ctx context.Context, userId string, parameters entity.QueryParameters) ([]Subscription, int64, error)
	GetAllIt(ctx context.Context, userId string) iter.Seq[Subscription]
}
