package query

import (
	"context"

	"golang.org/x/text/currency"

	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type FindPreferredCurrencyQuery struct {
}

func NewFindPreferredCurrencyQuery() FindPreferredCurrencyQuery {
	return FindPreferredCurrencyQuery{}
}

type FindPreferredCurrencyQueryHandler struct {
	accountService ports.AccountService
	authentication ports.Authentication
}

func NewFindPreferredCurrencyQueryHandler(
	accountService ports.AccountService,
	authentication ports.Authentication) *FindPreferredCurrencyQueryHandler {
	return &FindPreferredCurrencyQueryHandler{
		accountService: accountService,
		authentication: authentication,
	}
}

func (h FindPreferredCurrencyQueryHandler) Handle(
	ctx context.Context,
	_ FindPreferredCurrencyQuery) result.Result[currency.Unit] {
	connectedAccount := h.authentication.MustGetConnectedAccount(ctx)
	preferredCurrency := h.accountService.GetPreferredCurrency(ctx, connectedAccount.UserID())
	return result.Success(preferredCurrency)
}
