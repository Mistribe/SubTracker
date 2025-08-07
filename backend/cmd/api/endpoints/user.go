package endpoints

import (
	"github.com/gin-gonic/gin"

	"github.com/oleexo/subtracker/cmd/api/ginfx"
	"github.com/oleexo/subtracker/cmd/api/middlewares"
)

type UserEndpointGroup struct {
	routes      []ginfx.Route
	middlewares []gin.HandlerFunc
}

func NewUserEndpointGroup(
	preferredCurrencyEndpoint *UserGetPreferredCurrencyEndpoint,
	updatePreferredCurrencyEndpoint *UserUpdatePreferredCurrencyEndpoint,
	updateProfileEndpoint *UserUpdateProfileEndpoint,
	authenticationMiddleware *middlewares.AuthenticationMiddleware) *UserEndpointGroup {
	return &UserEndpointGroup{
		routes: []ginfx.Route{
			preferredCurrencyEndpoint,
			updatePreferredCurrencyEndpoint,
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
