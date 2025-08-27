package query

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/domain/subscription"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/shared"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type FindAllQuery struct {
	SearchText   string
	Recurrencies []subscription.RecurrencyType
	FromDate     *time.Time
	ToDate       *time.Time
	Users        []uuid.UUID
	Providers    []uuid.UUID
	WithInactive bool
	Limit        int64
	Offset       int64
}

func NewFindAllQuery(
	searchText string,
	recurrencies []subscription.RecurrencyType,
	fromDate *time.Time,
	toDate *time.Time,
	users []uuid.UUID,
	providers []uuid.UUID,
	withInactive bool,
	size, page int64) FindAllQuery {
	return FindAllQuery{
		SearchText:   searchText,
		Limit:        size,
		Offset:       page,
		WithInactive: withInactive,
		Recurrencies: recurrencies,
		FromDate:     fromDate,
		ToDate:       toDate,
		Users:        users,
		Providers:    providers,
	}
}

type FindAllQueryHandler struct {
	subscriptionRepository ports.SubscriptionRepository
	userService            ports.UserService
	exchange               ports.Exchange
	authService            ports.AuthService
}

func NewFindAllQueryHandler(
	subscriptionRepository ports.SubscriptionRepository,
	userService ports.UserService,
	exchange ports.Exchange,
	authService ports.AuthService) *FindAllQueryHandler {
	return &FindAllQueryHandler{
		subscriptionRepository: subscriptionRepository,
		authService:            authService,
		userService:            userService,
		exchange:               exchange,
	}
}

func (h FindAllQueryHandler) Handle(
	ctx context.Context,
	query FindAllQuery) result.Result[shared.PaginatedResponse[subscription.Subscription]] {
	userId := h.authService.MustGetUserId(ctx)
	parameters := ports.NewSubscriptionQueryParameters(query.SearchText,
		query.Recurrencies,
		query.FromDate,
		query.ToDate,
		query.Users,
		query.Providers,
		query.WithInactive,
		query.Limit,
		query.Offset)
	subs, count, err := h.subscriptionRepository.GetAllForUser(ctx, userId, parameters)
	if err != nil {
		return result.Fail[shared.PaginatedResponse[subscription.Subscription]](err)
	}

	preferredCurrency := h.userService.GetPreferredCurrency(ctx, userId)

	for _, sub := range subs {
		if sub.CustomPrice() != nil && sub.CustomPrice().Currency() != preferredCurrency {
			convertedPrice, err := h.exchange.ToCurrencyAt(ctx, sub.CustomPrice(), preferredCurrency, sub.StartDate())
			if err != nil {
				return result.Fail[shared.PaginatedResponse[subscription.Subscription]](err)
			}

			sub.CustomPrice().SetAmount(convertedPrice)
		}
	}

	return result.Success(shared.NewPaginatedResponse(subs, count))
}
