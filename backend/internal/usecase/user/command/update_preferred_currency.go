package command

import (
	"context"

	"golang.org/x/text/currency"

	"github.com/mistribe/subtracker/internal/domain/auth"
	"github.com/mistribe/subtracker/internal/domain/user"
	"github.com/mistribe/subtracker/internal/ports"
	auth2 "github.com/mistribe/subtracker/internal/usecase/auth"
	"github.com/mistribe/subtracker/pkg/langext/result"
	"github.com/mistribe/subtracker/pkg/x"
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
	userRepository ports.UserRepository
}

func NewUpdatePreferredCurrencyCommandHandler(userRepository ports.UserRepository) *UpdatePreferredCurrencyCommandHandler {
	return &UpdatePreferredCurrencyCommandHandler{
		userRepository: userRepository,
	}
}

func (h UpdatePreferredCurrencyCommandHandler) Handle(
	ctx context.Context,
	cmd UpdatePreferredCurrencyCommand) result.Result[bool] {
	userId, ok := auth2.GetUserIdFromContext(ctx)
	if !ok {
		return result.Fail[bool](auth.ErrUnauthorized)
	}
	profile, err := h.userRepository.GetUser(ctx, userId)
	if err != nil {
		return result.Fail[bool](err)
	}

	if profile == nil {
		return h.createProfile(ctx, userId, cmd)
	}

	return h.updateProfile(ctx, cmd, profile)
}

func (h UpdatePreferredCurrencyCommandHandler) createProfile(
	ctx context.Context,
	userId string,
	cmd UpdatePreferredCurrencyCommand) result.Result[bool] {

	profile := user.New(userId, cmd.Currency, x.P(user.PlanFree))

	if err := h.userRepository.Save(ctx, profile); err != nil {
		return result.Fail[bool](err)
	}

	return result.Success(true)
}

func (h UpdatePreferredCurrencyCommandHandler) updateProfile(
	ctx context.Context,
	cmd UpdatePreferredCurrencyCommand,
	profile user.User) result.Result[bool] {
	profile.SetCurrency(cmd.Currency)

	if err := h.userRepository.Save(ctx, profile); err != nil {
		return result.Fail[bool](err)
	}

	return result.Success(true)
}
