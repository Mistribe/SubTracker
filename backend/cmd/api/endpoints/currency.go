package endpoints

import (
	"github.com/gin-gonic/gin"
	"github.com/oleexo/subtracker/cmd/api/ginfx"
	"github.com/oleexo/subtracker/cmd/api/middlewares"
)

type CurrencyEndpointGroup struct {
	routes      []ginfx.Route
	middlewares []gin.HandlerFunc
}

func NewCurrencyGroupEndpointGroup(
	authenticationMiddleware *middlewares.AuthenticationMiddleware) *CurrencyEndpointGroup {
	return &CurrencyEndpointGroup{
		routes: nil,
		middlewares: []gin.HandlerFunc{
			authenticationMiddleware.Middleware(),
		},
	}
}

func (g CurrencyEndpointGroup) Prefix() string {
	return "/currencies"
}

func (g CurrencyEndpointGroup) Routes() []ginfx.Route {
	return g.routes
}

func (g CurrencyEndpointGroup) Middlewares() []gin.HandlerFunc {
	return g.middlewares
}
