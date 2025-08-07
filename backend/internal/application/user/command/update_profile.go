package command

import (
	"context"
	"net/http"

	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/internal/infrastructure/kinde"
	openapi "github.com/oleexo/subtracker/internal/infrastructure/kinde/gen"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type UpdateProfileCommand struct {
	GivenName  string
	FamilyName string
}

func NewUpdateProfileCommand(givenName, familyName string) UpdateProfileCommand {
	return UpdateProfileCommand{GivenName: givenName, FamilyName: familyName}
}

type UpdateProfileCommandHandler struct {
	authService    auth.Service
	kindeGenerator kinde.TokenGenerator
}

func NewUpdateProfileCommandHandler(
	kindeGenerator kinde.TokenGenerator,
	authService auth.Service) *UpdateProfileCommandHandler {
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
		func(ctx context.Context, client *openapi.APIClient) (*openapi.UpdateUserResponse, *http.Response, error) {
			return client.UsersAPI.UpdateUser(ctx).Id(userId).UpdateUserRequest(openapi.UpdateUserRequest{
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
