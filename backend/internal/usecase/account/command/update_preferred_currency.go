package command

import (
	"context"

	"golang.org/x/text/currency"

	"github.com/mistribe/subtracker/internal/domain/account"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type UpdatePreferredCurrencyCommand struct {
	Currency currency.Unit
}

func NewUpdatePreferredCurrencyCommand(currency currency.Unit) UpdatePreferredCurrencyCommand {
	return UpdatePreferredCurrencyCommand{
		Currency: currency,
	}
}

type UpdatePreferredCurrencyCommandHandler struct {
	accountRepository ports.AccountRepository
	authentication    ports.Authentication
}

func NewUpdatePreferredCurrencyCommandHandler(accountRepository ports.AccountRepository,
	authentication ports.Authentication) *UpdatePreferredCurrencyCommandHandler {
	return &UpdatePreferredCurrencyCommandHandler{
		accountRepository: accountRepository,
		authentication:    authentication,
	}
}

func (h UpdatePreferredCurrencyCommandHandler) Handle(ctx context.Context,
	cmd UpdatePreferredCurrencyCommand) result.Result[bool] {
	connectedAccount := h.authentication.MustGetConnectedAccount(ctx)
	acc, err := h.accountRepository.GetById(ctx, connectedAccount.UserID())
	if err != nil {
		return result.Fail[bool](err)
	}

	if acc == nil {
		return h.createAccount(ctx, connectedAccount.UserID(), cmd)
	}

	return h.updateAccount(ctx, cmd, acc)
}

func (h UpdatePreferredCurrencyCommandHandler) createAccount(
	ctx context.Context,
	userID types.UserID,
	cmd UpdatePreferredCurrencyCommand) result.Result[bool] {

	acc := account.New(userID,
		cmd.Currency,
		types.PlanFree,
		types.RoleUser,
		nil,
	)

	if err := h.accountRepository.Save(ctx, acc); err != nil {
		return result.Fail[bool](err)
	}

	return result.Success(true)
}

func (h UpdatePreferredCurrencyCommandHandler) updateAccount(
	ctx context.Context,
	cmd UpdatePreferredCurrencyCommand,
	acc account.Account) result.Result[bool] {
	acc.SetCurrency(cmd.Currency)

	if err := h.accountRepository.Save(ctx, acc); err != nil {
		return result.Fail[bool](err)
	}

	return result.Success(true)
}
