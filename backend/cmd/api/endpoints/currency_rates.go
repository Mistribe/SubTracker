package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/currency/query"
	"github.com/oleexo/subtracker/internal/domain/currency"
	"github.com/oleexo/subtracker/pkg/slicesx"
)

// CurrencyRateResponse represents the response body for currency conversion
type CurrencyRateResponse struct {
	Timestamp time.Time          `json:"timestamp"`
	Rates     map[string]float64 `json:"rates"`
}

// CurrencyRateEndpoint handles currency conversion requests
type CurrencyRateEndpoint struct {
	convertHandler core.QueryHandler[query.CurrencyRateQuery, query.CurrencyRateResponse]
}

// NewCurrencyRateEndpoint creates a new CurrencyRateEndpoint
func NewCurrencyRateEndpoint(
	convertHandler core.QueryHandler[query.CurrencyRateQuery, query.CurrencyRateResponse],
) *CurrencyRateEndpoint {
	return &CurrencyRateEndpoint{
		convertHandler: convertHandler,
	}
}

// Handle godoc
//
//	@Summary		Get Currency Rates
//	@Description	Get exchange rates for all currencies at a specific date
//	@Produce		json
//	@Param			date	query		string	false	"Conversion date in RFC3339 format (default: current time)"
//	@Success		200		{object}	CurrencyRateResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/currency/rates [get]
func (e CurrencyRateEndpoint) Handle(c *gin.Context) {
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
			return CurrencyRateResponse{
				Timestamp: c.Timestamp,
				Rates: slicesx.ToMap(c.Rates,
					func(key currency.Rate) string {
						return key.FromCurrency().String()
					},
					func(value currency.Rate) float64 {
						return value.ExchangeRate()
					},
				),
			}
		}))
}

func (e CurrencyRateEndpoint) Pattern() []string {
	return []string{
		"/rates",
	}
}

func (e CurrencyRateEndpoint) Method() string {
	return http.MethodGet
}

func (e CurrencyRateEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}
