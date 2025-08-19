package subscription

import (
	"context"
	"iter"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/entity"
	"github.com/oleexo/subtracker/pkg/types"
)

type QueryParameters struct {
	entity.QueryParameters

	SearchText string
	SortBy     string
	SortOrder  types.SortOrder
}

func NewQueryParameters(
	searchText string,
	sortBy string,
	sortOrder types.SortOrder,
	limit, offset int32) QueryParameters {
	return QueryParameters{
		QueryParameters: entity.QueryParameters{
			Limit:  limit,
			Offset: offset,
		},
		SearchText: searchText,
		SortBy:     sortBy,
		SortOrder:  sortOrder,
	}
}

type Repository interface {
	entity.Repository[Subscription]

	GetByIdForUser(ctx context.Context, userId string, id uuid.UUID) (Subscription, error)
	GetAll(ctx context.Context, parameters QueryParameters) ([]Subscription, int64, error)
	GetAllForUser(ctx context.Context, userId string, parameters QueryParameters) ([]Subscription, int64, error)
	GetAllIt(ctx context.Context, userId, searchText string) iter.Seq[Subscription]
}
