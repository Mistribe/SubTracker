package currency

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/mistribe/subtracker/internal/adapters/http/dto"
	"github.com/mistribe/subtracker/internal/domain/currency"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/currency/query"
	. "github.com/mistribe/subtracker/pkg/ginx"
	"github.com/mistribe/subtracker/pkg/x/herd"
)

// GetRateEndpoint handles currency conversion requests
type GetRateEndpoint struct {
	convertHandler ports.QueryHandler[query.CurrencyRateQuery, query.CurrencyRateResponse]
}

// NewGetRateEndpoint creates a new GetRateEndpoint
func NewGetRateEndpoint(
	convertHandler ports.QueryHandler[query.CurrencyRateQuery, query.CurrencyRateResponse],
) *GetRateEndpoint {
	return &GetRateEndpoint{
		convertHandler: convertHandler,
	}
}

// Handle godoc
//
//	@Summary		Get Currency Rates
//	@Description	Get exchange rates for all currencies at a specific date
//	@Tags			currencies
//	@Produce		json
//	@Param			date	query		string	false	"Conversion date in RFC3339 format (default: current time)"
//	@Success		200		{object}	dto.CurrencyRatesModel
//	@Failure		400		{object}	HttpErrorResponse
//	@Failure		500		{object}	HttpErrorResponse
//	@Router			/currencies/rates [get]
func (e GetRateEndpoint) Handle(c *gin.Context) {
	conversionDateParam := c.DefaultQuery("date", time.Now().Format(time.RFC3339))
	conversionDate, err := time.Parse(time.RFC3339, conversionDateParam)
	if err != nil {
		FromError(c, err)
		return
	}
	// Create the query
	conversionQuery := query.CurrencyRateQuery{
		ConversionDate: conversionDate,
	}

	// Execute the query
	r := e.convertHandler.Handle(c.Request.Context(), conversionQuery)
	FromResult(c,
		r,
		WithMapping[query.CurrencyRateResponse](func(c query.CurrencyRateResponse) any {
			return dto.CurrencyRatesModel{
				Timestamp: c.Timestamp,
				Rates: herd.Select(c.Rates,
					func(key currency.Rate) dto.CurrencyRateModel {
						return dto.CurrencyRateModel{
							Rate:     key.ExchangeRate(),
							Currency: key.ToCurrency().String(),
						}
					},
				),
			}
		}))
}

func (e GetRateEndpoint) Pattern() []string {
	return []string{
		"/rates",
	}
}

func (e GetRateEndpoint) Method() string {
	return http.MethodGet
}

func (e GetRateEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}
