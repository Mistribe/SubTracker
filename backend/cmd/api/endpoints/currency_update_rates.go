package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/currency/command"
)

type CurrencyUpdateRatesEndpoint struct {
	handler core.CommandHandler[command.UpdateCurrencyRatesCommand, bool]
}

func NewCurrencyUpdateRatesEndpoint(
	handler core.CommandHandler[command.UpdateCurrencyRatesCommand, bool],
) *CurrencyUpdateRatesEndpoint {
	return &CurrencyUpdateRatesEndpoint{
		handler: handler,
	}
}

// Handle godoc
//
//	@Summary		Update Currency Rates
//	@Description	Update exchange rates for all currencies
//	@Produce		json
//	@Success		200		{object}	bool
//	@Failure		400		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/currency/rates/update [post]
func (e CurrencyUpdateRatesEndpoint) Handle(c *gin.Context) {
	cmd := command.UpdateCurrencyRatesCommand{}
	r := e.handler.Handle(c.Request.Context(), cmd)
	handleResponse(c, r, nil)
}

func (e CurrencyUpdateRatesEndpoint) Pattern() []string {
	return []string{
		"/rates/update",
	}
}

func (e CurrencyUpdateRatesEndpoint) Method() string {
	return http.MethodPost
}

func (e CurrencyUpdateRatesEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}
