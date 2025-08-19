package subscription

import (
	"context"
	"iter"
	"time"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/entity"
	"github.com/oleexo/subtracker/pkg/types"
)

type QueryParameters struct {
	entity.QueryParameters

	SearchText   string
	WithInactive bool
	Recurrency   []RecurrencyType
	FromDate     *time.Time
	ToDate       *time.Time
	SortBy       SortableField
	SortOrder    types.SortOrder
}

func NewQueryParameters(
	searchText string,
	sortBy SortableField,
	sortOrder types.SortOrder,
	withInactive bool,
	limit, offset int32) QueryParameters {
	if sortBy == "" {
		sortBy = "friendly_name"
	}
	if sortOrder == "" {
		sortOrder = types.SortOrderAsc
	}
	return QueryParameters{
		QueryParameters: entity.QueryParameters{
			Limit:  limit,
			Offset: offset,
		},
		SearchText:   searchText,
		SortBy:       sortBy,
		SortOrder:    sortOrder,
		WithInactive: withInactive,
	}
}

type Repository interface {
	entity.Repository[Subscription]

	GetByIdForUser(ctx context.Context, userId string, id uuid.UUID) (Subscription, error)
	GetAll(ctx context.Context, parameters QueryParameters) ([]Subscription, int64, error)
	GetAllForUser(ctx context.Context, userId string, parameters QueryParameters) ([]Subscription, int64, error)
	GetAllIt(ctx context.Context, userId, searchText string) iter.Seq[Subscription]
}
