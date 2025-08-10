package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/currency/command"
	"github.com/oleexo/subtracker/internal/domain/currency"
	"github.com/oleexo/subtracker/pkg/slicesx"
)

type CurrencyRefreshRatesEndpoint struct {
	handler core.CommandHandler[command.RefreshCurrencyRatesCommand, command.RefreshCurrencyRatesResponse]
}

func NewCurrencyRefreshRatesEndpoint(
	handler core.CommandHandler[command.RefreshCurrencyRatesCommand, command.RefreshCurrencyRatesResponse],
) *CurrencyRefreshRatesEndpoint {
	return &CurrencyRefreshRatesEndpoint{
		handler: handler,
	}
}

type CurrencyRefreshRatesResponse struct {
	Timestamp time.Time          `json:"timestamp"`
	Rates     map[string]float64 `json:"rates"`
}

// Handle godoc
//
//	@Summary		refresh Currency Rates
//	@Description	refresh exchange rates for all currencies
//	@Produce		json
//	@Success		200	{object}	CurrencyRefreshRatesResponse
//	@Failure		400	{object}	HttpErrorResponse
//	@Failure		500	{object}	HttpErrorResponse
//	@Router			/currencies/rates/refresh [post]
func (e CurrencyRefreshRatesEndpoint) Handle(c *gin.Context) {
	cmd := command.RefreshCurrencyRatesCommand{}
	r := e.handler.Handle(c.Request.Context(), cmd)
	handleResponse(c, r,
		withMapping[command.RefreshCurrencyRatesResponse](func(c command.RefreshCurrencyRatesResponse) any {
			return CurrencyRefreshRatesResponse{
				Timestamp: c.Timestamp,
				Rates: slicesx.ToMap(c.Rates,
					func(key currency.Rate) string {
						return key.ToCurrency().String()
					},
					func(value currency.Rate) float64 {
						return value.ExchangeRate()
					},
				),
			}
		}),
	)
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
