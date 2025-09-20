package currency

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
	supportedEndpoint *SupportedEndpoint,
	convertEndpoint *GetRateEndpoint,
	authenticationMiddleware *middlewares.AuthenticationMiddleware) *EndpointGroup {
	return &EndpointGroup{
		routes: []ginfx.Endpoint{
			supportedEndpoint,
			convertEndpoint,
		},
		middlewares: []gin.HandlerFunc{
			authenticationMiddleware.Middleware(),
		},
	}
}

func (g EndpointGroup) Prefix() string {
	return "/currencies"
}

func (g EndpointGroup) Routes() []ginfx.Endpoint {
	return g.routes
}

func (g EndpointGroup) Middlewares() []gin.HandlerFunc {
	return g.middlewares
}
