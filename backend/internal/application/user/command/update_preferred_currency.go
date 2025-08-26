package command

import (
	"context"

	"golang.org/x/text/currency"

	auth2 "github.com/mistribe/subtracker/internal/application/auth"
	"github.com/mistribe/subtracker/internal/domain/auth"
	"github.com/mistribe/subtracker/internal/domain/user"
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
	userRepository user.Repository
}

func NewUpdatePreferredCurrencyCommandHandler(userRepository user.Repository) *UpdatePreferredCurrencyCommandHandler {
	return &UpdatePreferredCurrencyCommandHandler{
		userRepository: userRepository,
	}
}

func (h UpdatePreferredCurrencyCommandHandler) Handle(
	ctx context.Context,
	cmd UpdatePreferredCurrencyCommand) result.Result[bool] {
	userId, ok := auth2.GetUserIdFromContext(ctx)
	if !ok {
		return result.Fail[bool](auth.ErrNotAuthorized)
	}
	profile, err := h.userRepository.GetUserProfile(ctx, userId)
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

	profile := user.NewProfile(userId, cmd.Currency)

	if err := h.userRepository.SaveProfile(ctx, profile); err != nil {
		return result.Fail[bool](err)
	}

	return result.Success(true)
}

func (h UpdatePreferredCurrencyCommandHandler) updateProfile(
	ctx context.Context,
	cmd UpdatePreferredCurrencyCommand,
	profile user.Profile) result.Result[bool] {
	profile.SetCurrency(cmd.Currency)

	if err := h.userRepository.SaveProfile(ctx, profile); err != nil {
		return result.Fail[bool](err)
	}

	return result.Success(true)
}
