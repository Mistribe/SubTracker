package command

import (
	"context"

	"golang.org/x/text/currency"

	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type UpdateProfileCommand struct {
	Currency currency.Unit
}

type UpdateProfileCommandHandler struct {
	authRepository auth.Repository
}

func NewUpdateProfileCommandHandler(authRepository auth.Repository) *UpdateProfileCommandHandler {
	return &UpdateProfileCommandHandler{
		authRepository: authRepository,
	}
}

func (h UpdateProfileCommandHandler) Handle(ctx context.Context,
	cmd UpdateProfileCommand) result.Result[auth.UserProfile] {
	userId, ok := auth.GetUserIdFromContext(ctx)
	if !ok {
		return result.Fail[auth.UserProfile](auth.ErrNotAuthorized)
	}
	profile, err := h.authRepository.GetUserProfile(ctx, userId)
	if err != nil {
		return result.Fail[auth.UserProfile](err)
	}

	if profile == nil {
		return h.createProfile(ctx, userId, cmd)
	}

	return h.updateProfile(ctx, cmd, profile)
}

func (h UpdateProfileCommandHandler) createProfile(ctx context.Context,
	userId string,
	cmd UpdateProfileCommand) result.Result[auth.UserProfile] {

	profile := auth.NewUserProfile(userId, cmd.Currency)

	if err := h.authRepository.SaveProfile(ctx, profile); err != nil {
		return result.Fail[auth.UserProfile](err)
	}

	return result.Success(profile)
}

func (h UpdateProfileCommandHandler) updateProfile(ctx context.Context,
	cmd UpdateProfileCommand,
	profile auth.UserProfile) result.Result[auth.UserProfile] {
	profile.SetCurrency(cmd.Currency)

	if err := h.authRepository.SaveProfile(ctx, profile); err != nil {
		return result.Fail[auth.UserProfile](err)
	}

	return result.Success(profile)
}
