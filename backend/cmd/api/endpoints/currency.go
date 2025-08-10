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
	supportedEndpoint *CurrencySupportedEndpoint,
	convertEndpoint *CurrencyGetRateEndpoint,
	updateRatesEndpoint *CurrencyRefreshRatesEndpoint,
	authenticationMiddleware *middlewares.AuthenticationMiddleware) *CurrencyEndpointGroup {
	return &CurrencyEndpointGroup{
		routes: []ginfx.Route{
			supportedEndpoint,
			convertEndpoint,
			updateRatesEndpoint,
		},
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
