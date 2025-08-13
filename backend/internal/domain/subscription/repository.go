package subscription

import (
	"context"
	"iter"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/entity"
)

type QueryParameters struct {
	entity.QueryParameters

	SearchText string
}

func NewQueryParameters(
	searchText string,
	limit, offset int32) QueryParameters {
	return QueryParameters{
		QueryParameters: entity.QueryParameters{
			Limit:  limit,
			Offset: offset,
		},
		SearchText: searchText,
	}
}

type Repository interface {
	entity.Repository[Subscription]

	GetByIdForUser(ctx context.Context, userId string, id uuid.UUID) (Subscription, error)
	GetAll(ctx context.Context, parameters QueryParameters) ([]Subscription, int64, error)
	GetAllForUser(ctx context.Context, userId string, parameters QueryParameters) ([]Subscription, int64, error)
	GetAllIt(ctx context.Context, userId, searchText string) iter.Seq[Subscription]
}
