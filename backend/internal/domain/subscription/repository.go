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
	Recurrencies []RecurrencyType
	FromDate     *time.Time
	ToDate       *time.Time
	Users        []uuid.UUID
	Providers    []uuid.UUID
	WithInactive bool
	SortBy       SortableField
	SortOrder    types.SortOrder
}

func NewQueryParameters(
	searchText string,
	recurrencies []RecurrencyType,
	fromDate *time.Time,
	toDate *time.Time,
	users []uuid.UUID,
	providers []uuid.UUID,
	withInactive bool,
	sortBy SortableField,
	sortOrder types.SortOrder,
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
		Recurrencies: recurrencies,
		FromDate:     fromDate,
		ToDate:       toDate,
		Users:        users,
		Providers:    providers,
	}
}

type Repository interface {
	entity.Repository[Subscription]

	GetByIdForUser(ctx context.Context, userId string, id uuid.UUID) (Subscription, error)
	GetAll(ctx context.Context, parameters QueryParameters) ([]Subscription, int64, error)
	GetAllForUser(ctx context.Context, userId string, parameters QueryParameters) ([]Subscription, int64, error)
	GetAllIt(ctx context.Context, userId, searchText string) iter.Seq[Subscription]
}
