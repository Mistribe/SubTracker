package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/mistribe/subtracker/internal/application/core"
	"github.com/mistribe/subtracker/internal/application/currency/query"
	"github.com/mistribe/subtracker/internal/domain/currency"
	"github.com/mistribe/subtracker/pkg/slicesx"
)

// CurrencyGetRateResponse represents the response body for currency conversion
type CurrencyGetRateResponse struct {
	Timestamp time.Time          `json:"timestamp"`
	Rates     map[string]float64 `json:"rates"`
}

// CurrencyGetRateEndpoint handles currency conversion requests
type CurrencyGetRateEndpoint struct {
	convertHandler core.QueryHandler[query.CurrencyRateQuery, query.CurrencyRateResponse]
}

// NewCurrencyGetRateEndpoint creates a new CurrencyGetRateEndpoint
func NewCurrencyGetRateEndpoint(
	convertHandler core.QueryHandler[query.CurrencyRateQuery, query.CurrencyRateResponse],
) *CurrencyGetRateEndpoint {
	return &CurrencyGetRateEndpoint{
		convertHandler: convertHandler,
	}
}

// Handle godoc
//
//	@Summary		Get Currency Rates
//	@Description	Get exchange rates for all currencies at a specific date
//	@Produce		json
//	@Param			date	query		string	false	"Conversion date in RFC3339 format (default: current time)"
//	@Success		200		{object}	CurrencyRatesModel
//	@Failure		400		{object}	HttpErrorResponse
//	@Failure		500		{object}	HttpErrorResponse
//	@Router			/currencies/rates [get]
func (e CurrencyGetRateEndpoint) Handle(c *gin.Context) {
	conversionDateParam := c.DefaultQuery("date", time.Now().Format(time.RFC3339))
	conversionDate, err := time.Parse(time.RFC3339, conversionDateParam)
	if err != nil {
		handleErrorResponse(c, err)
		return
	}
	// Create the query
	conversionQuery := query.CurrencyRateQuery{
		ConversionDate: conversionDate,
	}

	// Execute the query
	r := e.convertHandler.Handle(c.Request.Context(), conversionQuery)
	handleResponse(c,
		r,
		withMapping[query.CurrencyRateResponse](func(c query.CurrencyRateResponse) any {
			return CurrencyRatesModel{
				Timestamp: c.Timestamp,
				Rates: slicesx.Select(c.Rates,
					func(key currency.Rate) CurrencyRateModel {
						return CurrencyRateModel{
							Rate:     key.ExchangeRate(),
							Currency: key.ToCurrency().String(),
						}
					},
				),
			}
		}))
}

func (e CurrencyGetRateEndpoint) Pattern() []string {
	return []string{
		"/rates",
	}
}

func (e CurrencyGetRateEndpoint) Method() string {
	return http.MethodGet
}

func (e CurrencyGetRateEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}
