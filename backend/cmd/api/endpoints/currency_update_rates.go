package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/currency/command"
)

type CurrencyRefreshRatesEndpoint struct {
	handler core.CommandHandler[command.RefreshCurrencyRatesCommand, bool]
}

func NewCurrencyRefreshRatesEndpoint(
	handler core.CommandHandler[command.RefreshCurrencyRatesCommand, bool],
) *CurrencyRefreshRatesEndpoint {
	return &CurrencyRefreshRatesEndpoint{
		handler: handler,
	}
}

// Handle godoc
//
//	@Summary		refresh Currency Rates
//	@Description	refresh exchange rates for all currencies
//	@Produce		json
//	@Success		200	{object}	bool
//	@Failure		400	{object}	HttpErrorResponse
//	@Failure		500	{object}	HttpErrorResponse
//	@Router			/currency/rates/refresh [post]
func (e CurrencyRefreshRatesEndpoint) Handle(c *gin.Context) {
	cmd := command.RefreshCurrencyRatesCommand{}
	r := e.handler.Handle(c.Request.Context(), cmd)
	handleResponse(c, r, nil)
}

func (e CurrencyRefreshRatesEndpoint) Pattern() []string {
	return []string{
		"/rates/refresh",
	}
}

func (e CurrencyRefreshRatesEndpoint) Method() string {
	return http.MethodPost
}

func (e CurrencyRefreshRatesEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}
