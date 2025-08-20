package query

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/internal/domain/currency"
	"github.com/oleexo/subtracker/internal/domain/subscription"
	"github.com/oleexo/subtracker/internal/domain/user"
	"github.com/oleexo/subtracker/pkg/langext/result"
	"github.com/oleexo/subtracker/pkg/types"
)

type FindAllQuery struct {
	SearchText   string
	Recurrencies []subscription.RecurrencyType
	FromDate     *time.Time
	ToDate       *time.Time
	Users        []uuid.UUID
	Providers    []uuid.UUID
	WithInactive bool
	SortBy       subscription.SortableField
	SortOrder    types.SortOrder
	Limit        int32
	Offset       int32
}

func NewFindAllQuery(searchText string,
	recurrencies []subscription.RecurrencyType,
	fromDate *time.Time,
	toDate *time.Time,
	users []uuid.UUID,
	providers []uuid.UUID,
	withInactive bool,
	sortBy subscription.SortableField,
	sortOrder types.SortOrder,
	size, page int32) FindAllQuery {
	return FindAllQuery{
		SearchText:   searchText,
		Limit:        size,
		Offset:       page,
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

type FindAllQueryHandler struct {
	subscriptionRepository subscription.Repository
	userService            user.Service
	currencyService        currency.Service
	authService            auth.Service
}

func NewFindAllQueryHandler(
	subscriptionRepository subscription.Repository,
	userService user.Service,
	currencyService currency.Service,
	authService auth.Service) *FindAllQueryHandler {
	return &FindAllQueryHandler{
		subscriptionRepository: subscriptionRepository,
		authService:            authService,
		userService:            userService,
		currencyService:        currencyService,
	}
}

func (h FindAllQueryHandler) Handle(
	ctx context.Context,
	query FindAllQuery) result.Result[core.PaginatedResponse[subscription.Subscription]] {
	userId := h.authService.MustGetUserId(ctx)
	parameters := subscription.NewQueryParameters(query.SearchText,
		query.Recurrencies,
		query.FromDate,
		query.ToDate,
		query.Users,
		query.Providers,
		query.WithInactive,
		query.SortBy,
		query.SortOrder,
		query.Limit,
		query.Offset)
	subs, count, err := h.subscriptionRepository.GetAllForUser(ctx, userId, parameters)
	if err != nil {
		return result.Fail[core.PaginatedResponse[subscription.Subscription]](err)
	}

	preferredCurrency := h.userService.GetPreferredCurrency(ctx, userId)

	for _, sub := range subs {
		if sub.CustomPrice() != nil && sub.CustomPrice().Currency() != preferredCurrency {
			convertedPrice, err := h.currencyService.ConvertTo(ctx, sub.CustomPrice(), preferredCurrency)
			if err != nil {
				return result.Fail[core.PaginatedResponse[subscription.Subscription]](err)
			}

			sub.CustomPrice().SetAmount(convertedPrice)
		}
	}

	return result.Success(core.NewPaginatedResponse(subs, count))
}
