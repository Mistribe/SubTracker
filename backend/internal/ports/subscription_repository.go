package ports

import (
	"context"
	"iter"
	"time"

	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/domain/subscription"
)

type SubscriptionQueryParameters struct {
	QueryParameters

	SearchText   string
	Recurrencies []subscription.RecurrencyType
	FromDate     *time.Time
	ToDate       *time.Time
	Users        []uuid.UUID
	Providers    []uuid.UUID
	WithInactive bool
}

func NewSubscriptionQueryParameters(
	searchText string,
	recurrencies []subscription.RecurrencyType,
	fromDate *time.Time,
	toDate *time.Time,
	users []uuid.UUID,
	providers []uuid.UUID,
	withInactive bool,
	limit, offset int64) SubscriptionQueryParameters {
	return SubscriptionQueryParameters{
		QueryParameters: QueryParameters{
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

type SubscriptionRepository interface {
	Repository[subscription.Subscription]

	GetByIdForUser(ctx context.Context, userId string, id uuid.UUID) (subscription.Subscription, error)
	GetAll(ctx context.Context, parameters SubscriptionQueryParameters) ([]subscription.Subscription, int64, error)
	GetAllForUser(ctx context.Context, userId string,
		parameters SubscriptionQueryParameters) ([]subscription.Subscription, int64, error)
	GetAllIt(ctx context.Context, userId, searchText string) iter.Seq[subscription.Subscription]
}
