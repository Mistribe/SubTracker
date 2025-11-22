package label

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
	createEndpoint *CreateEndpoint,
	updateEndpoint *UpdateEndpoint,
	deleteEndpoint *DeleteEndpoint,
	getEndpoint *GetEndpoint,
	getAllEndpoint *GetAllEndpoint,
	labelQuotaUsageEndpoint *GetQuotaUsageEndpoint,
	exportEndpoint *ExportEndpoint,
	authenticationMiddleware *middlewares.AuthenticationMiddleware) *EndpointGroup {
	return &EndpointGroup{
		routes: []ginfx.Endpoint{
			createEndpoint,
			updateEndpoint,
			deleteEndpoint,
			getEndpoint,
			getAllEndpoint,
			labelQuotaUsageEndpoint,
			exportEndpoint,
		},
		middlewares: []gin.HandlerFunc{
			authenticationMiddleware.Middleware(),
		},
	}
}

func (g EndpointGroup) Prefix() string {
	return "/labels"
}

func (g EndpointGroup) Routes() []ginfx.Endpoint {
	return g.routes
}

func (g EndpointGroup) Middlewares() []gin.HandlerFunc {
	return g.middlewares
}
