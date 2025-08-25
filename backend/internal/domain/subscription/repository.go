package subscription

import (
	"context"
	"iter"
	"time"

	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/domain/entity"
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
}

func NewQueryParameters(
	searchText string,
	recurrencies []RecurrencyType,
	fromDate *time.Time,
	toDate *time.Time,
	users []uuid.UUID,
	providers []uuid.UUID,
	withInactive bool,
	limit, offset int64) QueryParameters {
	return QueryParameters{
		QueryParameters: entity.QueryParameters{
			Limit:  limit,
			Offset: offset,
		},
		SearchText:   searchText,
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
