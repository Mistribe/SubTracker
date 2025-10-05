package ports

import (
	"context"
	"iter"
	"time"

	"github.com/mistribe/subtracker/internal/domain/subscription"
	"github.com/mistribe/subtracker/internal/domain/types"
)

type SubscriptionQueryParameters struct {
	QueryParameters

	SearchText   string
	Recurrencies []subscription.RecurrencyType
	FromDate     *time.Time
	ToDate       *time.Time
	Users        []types.UserID
	Providers    []types.ProviderID
	WithInactive bool
}

func NewSubscriptionQueryParameters(
	searchText string,
	recurrencies []subscription.RecurrencyType,
	fromDate *time.Time,
	toDate *time.Time,
	users []types.UserID,
	providers []types.ProviderID,
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
	Repository[types.SubscriptionID, subscription.Subscription]

	GetByIdForUser(ctx context.Context, userId types.UserID, id types.SubscriptionID) (subscription.Subscription, error)
	GetAll(ctx context.Context, parameters SubscriptionQueryParameters) ([]subscription.Subscription, int64, error)
	GetAllForUser(ctx context.Context, userId types.UserID,
		parameters SubscriptionQueryParameters) ([]subscription.Subscription, int64, error)
	GetAllIt(ctx context.Context, userId types.UserID, searchText string) iter.Seq[subscription.Subscription]
}
