package query

import (
	"context"

	"golang.org/x/text/currency"

	"github.com/mistribe/subtracker/internal/domain/auth"
	"github.com/mistribe/subtracker/internal/ports"
	auth2 "github.com/mistribe/subtracker/internal/usecase/auth"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type FindPreferredCurrencyQuery struct {
}

func NewFindPreferredCurrencyQuery() FindPreferredCurrencyQuery {
	return FindPreferredCurrencyQuery{}
}

type FindPreferredCurrencyQueryHandler struct {
	userService ports.UserService
}

func NewFindPreferredCurrencyQueryHandler(userService ports.UserService) *FindPreferredCurrencyQueryHandler {
	return &FindPreferredCurrencyQueryHandler{
		userService: userService,
	}
}

func (h FindPreferredCurrencyQueryHandler) Handle(
	ctx context.Context,
	_ FindPreferredCurrencyQuery) result.Result[currency.Unit] {
	userId, ok := auth2.GetUserIdFromContext(ctx)
	if !ok {
		return result.Fail[currency.Unit](auth.ErrNotAuthorized)
	}
	return result.Success(h.userService.GetPreferredCurrency(ctx, userId))
}
