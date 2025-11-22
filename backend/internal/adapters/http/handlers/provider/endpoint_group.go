package provider

import (
	"github.com/gin-gonic/gin"

	"github.com/mistribe/subtracker/internal/adapters/http/router/ginfx"
	"github.com/mistribe/subtracker/internal/adapters/http/router/middlewares"
)

type EndpointGroup struct {
	routes      []ginfx.Endpoint
	middlewares []gin.HandlerFunc
}

func (s EndpointGroup) Prefix() string {
	return "/providers"
}

func (s EndpointGroup) Routes() []ginfx.Endpoint {
	return s.routes
}

func (s EndpointGroup) Middlewares() []gin.HandlerFunc {
	return s.middlewares
}

func NewEndpointGroup(
	getEndpoint *GetEndpoint,
	getAllEndpoint *GetAllEndpoint,
	createEndpoint *CreateEndpoint,
	updateEndpoint *UpdateEndpoint,
	deleteEndpoint *DeleteEndpoint,
	providerQuotaUsageEndpoint *GetQuotaUsageEndpoint,
	exportEndpoint *ExportEndpoint,
	authenticationMiddleware *middlewares.AuthenticationMiddleware) *EndpointGroup {
	return &EndpointGroup{
		routes: []ginfx.Endpoint{
			getEndpoint,
			getAllEndpoint,
			createEndpoint,
			updateEndpoint,
			deleteEndpoint,
			providerQuotaUsageEndpoint,
			exportEndpoint,
		},
		middlewares: []gin.HandlerFunc{
			authenticationMiddleware.Middleware(),
		},
	}
}
