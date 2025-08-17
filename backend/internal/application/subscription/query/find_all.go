package query

import (
	"context"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/internal/domain/currency"
	"github.com/oleexo/subtracker/internal/domain/subscription"
	"github.com/oleexo/subtracker/internal/domain/user"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type FindAllQuery struct {
	SearchText string
	Limit      int32
	Offset     int32
}

func NewFindAllQuery(searchText string, size, page int32) FindAllQuery {
	return FindAllQuery{
		SearchText: searchText,
		Limit:      size,
		Offset:     page,
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
	parameters := subscription.NewQueryParameters(query.SearchText, query.Limit, query.Offset)
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
