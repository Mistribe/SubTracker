package command

import (
	"context"

	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type DeleteUserCommand struct {
}

func NewDeleteUserCommand() DeleteUserCommand {
	return DeleteUserCommand{}
}

type DeleteUserCommandHandler struct {
	identityProvider ports.IdentityProvider
	userRepository   ports.UserRepository
	authService      ports.AuthenticationService
}

func NewDeleteUserCommandHandler(
	identityProvider ports.IdentityProvider,
	userRepository ports.UserRepository,
	authService ports.AuthenticationService) *DeleteUserCommandHandler {
	return &DeleteUserCommandHandler{
		identityProvider: identityProvider,
		userRepository:   userRepository,
		authService:      authService,
	}
}

func (h DeleteUserCommandHandler) Handle(ctx context.Context, cmd DeleteUserCommand) result.Result[bool] {
	userId := h.authService.MustGetUserId(ctx)

	if err := h.identityProvider.DeleteUser(ctx, userId); err != nil {
		return result.Fail[bool](err)
	}

	// todo remove all user data

	return result.Success(true)
}
