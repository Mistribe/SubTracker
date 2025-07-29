package endpoints

import (
	"github.com/gin-gonic/gin"
	"github.com/oleexo/subtracker/cmd/api/ginfx"
	"github.com/oleexo/subtracker/cmd/api/middlewares"
)

type ProviderEndpointGroup struct {
	routes      []ginfx.Route
	middlewares []gin.HandlerFunc
}

func (s ProviderEndpointGroup) Prefix() string {
	return "/providers"
}

func (s ProviderEndpointGroup) Routes() []ginfx.Route {
	return s.routes
}

func (s ProviderEndpointGroup) Middlewares() []gin.HandlerFunc {
	return s.middlewares
}

func NewProviderEndpointGroup(
	getEndpoint *ProviderGetEndpoint,
	authenticationMiddleware *middlewares.AuthenticationMiddleware) *SubscriptionEndpointGroup {
	return &SubscriptionEndpointGroup{
		routes: []ginfx.Route{
			getEndpoint,
		},
		middlewares: []gin.HandlerFunc{
			authenticationMiddleware.Middleware(),
		},
	}
}
