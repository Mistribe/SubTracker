package query

import (
	"context"
	"time"

	"github.com/mistribe/subtracker/internal/domain/currency"
	"github.com/mistribe/subtracker/internal/domain/subscription"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/shared"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type FindAllQuery struct {
	SearchText   string
	Recurrencies []subscription.RecurrencyType
	FromDate     *time.Time
	ToDate       *time.Time
	Users        []types.UserID
	Providers    []types.ProviderID
	WithInactive bool
	Limit        int64
	Offset       int64
}

func NewFindAllQuery(
	searchText string,
	recurrencies []subscription.RecurrencyType,
	fromDate *time.Time,
	toDate *time.Time,
	users []types.UserID,
	providers []types.ProviderID,
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
	userService            ports.AccountService
	exchange               ports.Exchange
	authentication         ports.Authentication
	authorization          ports.Authorization
}

func NewFindAllQueryHandler(
	subscriptionRepository ports.SubscriptionRepository,
	userService ports.AccountService,
	exchange ports.Exchange,
	authService ports.Authentication,
	authorization ports.Authorization) *FindAllQueryHandler {
	return &FindAllQueryHandler{
		subscriptionRepository: subscriptionRepository,
		authentication:         authService,
		userService:            userService,
		exchange:               exchange,
		authorization:          authorization,
	}
}

func (h FindAllQueryHandler) Handle(
	ctx context.Context,
	query FindAllQuery) result.Result[shared.PaginatedResponse[subscription.Subscription]] {
	connectedAccount := h.authentication.MustGetConnectedAccount(ctx)
	parameters := ports.NewSubscriptionQueryParameters(query.SearchText,
		query.Recurrencies,
		query.FromDate,
		query.ToDate,
		query.Users,
		query.Providers,
		query.WithInactive,
		query.Limit,
		query.Offset)
	userId := connectedAccount.UserID()
	subs, count, err := h.subscriptionRepository.GetAllForUser(ctx, userId, parameters)
	if err != nil {
		return result.Fail[shared.PaginatedResponse[subscription.Subscription]](err)
	}

	preferredCurrency := h.userService.GetPreferredCurrency(ctx, userId)

	for _, sub := range subs {
		if sub.Price() != nil {
			var convertedPrice currency.Amount
			convertedPrice, err = h.exchange.ToCurrencyAt(ctx,
				sub.Price().Amount(),
				preferredCurrency,
				sub.StartDate())
			if err != nil {
				return result.Fail[shared.PaginatedResponse[subscription.Subscription]](err)
			}

			sub.Price().SetAmount(convertedPrice)
		}
	}

	return result.Success(shared.NewPaginatedResponse(subs, count))
}
