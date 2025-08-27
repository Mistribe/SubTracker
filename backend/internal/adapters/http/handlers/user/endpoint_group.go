package user

import (
	"github.com/gin-gonic/gin"

	"github.com/mistribe/subtracker/internal/adapters/http/router/ginfx"
	"github.com/mistribe/subtracker/internal/adapters/http/router/middlewares"
)

type EndpointGroup struct {
	routes      []ginfx.Endpoint
	middlewares []gin.HandlerFunc
}

func NewEndpointGroup(
	preferredCurrencyEndpoint *UserGetPreferredCurrencyEndpoint,
	updatePreferredCurrencyEndpoint *UserUpdatePreferredCurrencyEndpoint,
	updateProfileEndpoint *UserUpdateProfileEndpoint,
	deleteEndpoint *UserDeleteEndpoint,
	authenticationMiddleware *middlewares.AuthenticationMiddleware) *EndpointGroup {
	return &EndpointGroup{
		routes: []ginfx.Endpoint{
			preferredCurrencyEndpoint,
			updatePreferredCurrencyEndpoint,
			updateProfileEndpoint,
			deleteEndpoint,
		},
		middlewares: []gin.HandlerFunc{
			authenticationMiddleware.Middleware(),
		},
	}
}

func (u EndpointGroup) Prefix() string {
	return "/users"
}

func (u EndpointGroup) Routes() []ginfx.Endpoint {
	return u.routes
}

func (u EndpointGroup) Middlewares() []gin.HandlerFunc {
	return u.middlewares
}
