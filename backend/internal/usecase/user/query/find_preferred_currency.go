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
	accountRepository ports.AccountRepository
	authentication    ports.Authentication
}

func NewFindPreferredCurrencyQueryHandler(accountRepository ports.AccountRepository,
	authentication ports.Authentication) *FindPreferredCurrencyQueryHandler {
	return &FindPreferredCurrencyQueryHandler{
		accountRepository: accountRepository,
		authentication:    authentication,
	}
}

func (h FindPreferredCurrencyQueryHandler) Handle(
	ctx context.Context,
	_ FindPreferredCurrencyQuery) result.Result[currency.Unit] {
	connectedAccount := h.authentication.MustGetConnectedAccount(ctx)
	account, err := h.accountRepository.GetById(ctx, connectedAccount.UserID())
	if err != nil {
		return result.Fail[currency.Unit](err)
	}
	return result.Success(account.Currency())
}
