package subscription

import (
	"context"
	"iter"

	"github.com/oleexo/subtracker/internal/domain/entity"
)

type Repository interface {
	entity.Repository[Subscription]

	GetAll(ctx context.Context, parameters entity.QueryParameters) ([]Subscription, int64, error)
	GetAllIt(ctx context.Context) iter.Seq[Subscription]
}
