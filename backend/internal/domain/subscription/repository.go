package subscription

import (
	"context"

	"github.com/oleexo/subtracker/internal/domain/entity"
)

type Repository interface {
	entity.Repository[Subscription]

	GetAll(ctx context.Context, parameters entity.QueryParameters) ([]Subscription, int64, error)
}
