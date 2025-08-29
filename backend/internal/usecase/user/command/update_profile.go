package command

import (
	"context"

	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type UpdateProfileCommand struct {
	GivenName  string
	FamilyName string
}

func NewUpdateProfileCommand(givenName, familyName string) UpdateProfileCommand {
	return UpdateProfileCommand{GivenName: givenName, FamilyName: familyName}
}

type UpdateProfileCommandHandler struct {
	authService      ports.AuthService
	identityProvider ports.IdentityProvider
}

func NewUpdateProfileCommandHandler(
	identityProvider ports.IdentityProvider,
	authService ports.AuthService) *UpdateProfileCommandHandler {
	return &UpdateProfileCommandHandler{
		identityProvider: identityProvider,
		authService:      authService,
	}
}

func (h UpdateProfileCommandHandler) Handle(ctx context.Context, cmd UpdateProfileCommand) result.Result[bool] {
	// userId := h.authService.MustGetUserId(ctx)
	//	currentProfile, err := kinde.MakeRequest(h.identityProvider,
	//		func(ctx context.Context, client *openapi.APIClient) (*openapi.User, *http.Response, error) {
	//			return client.UsersAPI.GetUserData(ctx).Id(userId).Execute()
	//		},
	//	)
	//	if err != nil {
	//		return result.Fail[bool](err)
	//	}
	//
	//	if currentProfile == nil {
	//		return result.Fail[bool](auth.ErrNotAuthorized)
	//	}

	return result.Success(true)
}
