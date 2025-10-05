package command

import (
	"context"

	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type DeleteAccountCommand struct {
}

func NewDeleteAccountCommand() DeleteAccountCommand {
	return DeleteAccountCommand{}
}

type DeleteUserCommandHandler struct {
	identityProvider ports.IdentityProvider
	userRepository   ports.AccountRepository
	authService      ports.Authentication
}

func NewDeleteUserCommandHandler(
	identityProvider ports.IdentityProvider,
	userRepository ports.AccountRepository,
	authService ports.Authentication) *DeleteUserCommandHandler {
	return &DeleteUserCommandHandler{
		identityProvider: identityProvider,
		userRepository:   userRepository,
		authService:      authService,
	}
}

func (h DeleteUserCommandHandler) Handle(ctx context.Context, cmd DeleteAccountCommand) result.Result[bool] {
	connectedAccount := h.authService.MustGetConnectedAccount(ctx)

	if err := h.identityProvider.DeleteUser(ctx, connectedAccount.UserID()); err != nil {
		return result.Fail[bool](err)
	}

	// todo remove all user data

	return result.Success(true)
}
