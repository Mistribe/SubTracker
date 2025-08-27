package currency

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/mistribe/subtracker/internal/adapters/http/router/ginfx"
	"github.com/mistribe/subtracker/internal/adapters/http/router/middlewares"
)

type EndpointGroup struct {
	routes      []ginfx.Endpoint
	middlewares []gin.HandlerFunc
}

func NewEndpointGroup(
	supportedEndpoint *CurrencySupportedEndpoint,
	convertEndpoint *CurrencyGetRateEndpoint,
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

type CurrencyRateModel struct {
	Rate     float64 `json:"rate" binding:"required"`
	Currency string  `json:"currency" binding:"required"`
}
type CurrencyRatesModel struct {
	Rates     []CurrencyRateModel `json:"rates" binding:"required"`
	Timestamp time.Time           `json:"timestamp" binding:"required" format:"date-time"`
}
