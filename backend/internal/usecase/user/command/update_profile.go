package command

import (
	"context"
	"net/http"

	"github.com/mistribe/subtracker/internal/adapters/auth/kinde"
	openapi2 "github.com/mistribe/subtracker/internal/adapters/auth/kinde/gen"
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
	authService    ports.AuthService
	kindeGenerator kinde.TokenGenerator
}

func NewUpdateProfileCommandHandler(
	kindeGenerator kinde.TokenGenerator,
	authService ports.AuthService) *UpdateProfileCommandHandler {
	return &UpdateProfileCommandHandler{
		kindeGenerator: kindeGenerator,
		authService:    authService,
	}
}

func (h UpdateProfileCommandHandler) Handle(ctx context.Context, cmd UpdateProfileCommand) result.Result[bool] {
	userId := h.authService.MustGetUserId(ctx)
	//	currentProfile, err := kinde.MakeRequest(h.kindeGenerator,
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

	_, err := kinde.MakeRequest(h.kindeGenerator,
		func(ctx context.Context, client *openapi2.APIClient) (*openapi2.UpdateUserResponse, *http.Response, error) {
			return client.UsersAPI.UpdateUser(ctx).Id(userId).UpdateUserRequest(openapi2.UpdateUserRequest{
				GivenName:                &cmd.GivenName,
				FamilyName:               &cmd.FamilyName,
				Picture:                  nil,
				IsSuspended:              nil,
				IsPasswordResetRequested: nil,
				ProvidedId:               nil,
			}).Execute()
		})

	if err != nil {
		return result.Fail[bool](err)
	}

	return result.Success(true)
}
