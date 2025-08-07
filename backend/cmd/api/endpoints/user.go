package endpoints

import (
	"github.com/gin-gonic/gin"

	"github.com/oleexo/subtracker/cmd/api/ginfx"
	"github.com/oleexo/subtracker/cmd/api/middlewares"
	"github.com/oleexo/subtracker/internal/domain/user"
)

type UserEndpointGroup struct {
	routes      []ginfx.Route
	middlewares []gin.HandlerFunc
}

func NewUserEndpointGroup(
	profileEndpoint *UserGetPreferredCurrencyEndpoint,
	updateProfileEndpoint *UserUpdatePreferredCurrencyEndpoint,
	authenticationMiddleware *middlewares.AuthenticationMiddleware) *UserEndpointGroup {
	return &UserEndpointGroup{
		routes: []ginfx.Route{
			profileEndpoint,
			updateProfileEndpoint,
		},
		middlewares: []gin.HandlerFunc{
			authenticationMiddleware.Middleware(),
		},
	}
}

func (u UserEndpointGroup) Prefix() string {
	return "/users"
}

func (u UserEndpointGroup) Routes() []ginfx.Route {
	return u.routes
}

func (u UserEndpointGroup) Middlewares() []gin.HandlerFunc {
	return u.middlewares
}

type UserProfileModel struct {
	Currency string `json:"currency"`
}

func newUserProfileModel(source user.Profile) UserProfileModel {
	return UserProfileModel{
		Currency: source.Currency().String(),
	}
}
