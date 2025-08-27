package command

import (
	"context"
	"net/http"

	"github.com/mistribe/subtracker/internal/adapters/auth/kinde"
	openapi2 "github.com/mistribe/subtracker/internal/adapters/auth/kinde/gen"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type DeleteUserCommand struct {
}

func NewDeleteUserCommand() DeleteUserCommand {
	return DeleteUserCommand{}
}

type DeleteUserCommandHandler struct {
	kindeTokenGenerator kinde.TokenGenerator
	userRepository      ports.UserRepository
	authService         ports.AuthService
}

func NewDeleteUserCommandHandler(
	kindeTokenGenerator kinde.TokenGenerator,
	userRepository ports.UserRepository,
	authService ports.AuthService) *DeleteUserCommandHandler {
	return &DeleteUserCommandHandler{
		kindeTokenGenerator: kindeTokenGenerator,
		userRepository:      userRepository,
		authService:         authService,
	}
}

func (h DeleteUserCommandHandler) Handle(ctx context.Context, cmd DeleteUserCommand) result.Result[bool] {
	userId := h.authService.MustGetUserId(ctx)

	_, err := kinde.MakeRequest(h.kindeTokenGenerator,
		func(ctx context.Context, client *openapi2.APIClient) (*openapi2.SuccessResponse, *http.Response, error) {
			return client.UsersAPI.DeleteUser(ctx).Id(userId).Execute()
		})
	if err != nil {
		return result.Fail[bool](err)
	}

	// remove all user data

	return result.Success(true)
}
