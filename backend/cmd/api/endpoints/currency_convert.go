package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/currency/query"
)

// CurrencyConvertRequest represents the request body for currency conversion
type CurrencyConvertRequest struct {
	Amount         float64 `json:"amount" binding:"required"`
	FromCurrency   string  `json:"fromCurrency" binding:"required,len=3"`
	ToCurrency     string  `json:"toCurrency" binding:"required,len=3"`
	ConversionDate string  `json:"date,omitempty"`
}

// CurrencyConvertResponse represents the response body for currency conversion
type CurrencyConvertResponse struct {
	Amount       float64   `json:"amount"`
	FromCurrency string    `json:"fromCurrency"`
	ToCurrency   string    `json:"toCurrency"`
	Rate         float64   `json:"rate"`
	Date         time.Time `json:"date"`
}

// CurrencyConvertEndpoint handles currency conversion requests
type CurrencyConvertEndpoint struct {
	convertHandler core.QueryHandler[query.ConvertCurrencyQuery, query.ConvertCurrencyResult]
}

// NewCurrencyConvertEndpoint creates a new CurrencyConvertEndpoint
func NewCurrencyConvertEndpoint(
	convertHandler core.QueryHandler[query.ConvertCurrencyQuery, query.ConvertCurrencyResult],
) *CurrencyConvertEndpoint {
	return &CurrencyConvertEndpoint{
		convertHandler: convertHandler,
	}
}

// Handle godoc
//
//	@Summary		Convert Currency
//	@Description	Convert an amount from one currency to another at a specific date
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CurrencyConvertRequest	true	"Conversion Request"
//	@Success		200		{object}	CurrencyConvertResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/currencies/convert [post]
func (e CurrencyConvertEndpoint) Handle(c *gin.Context) {
	var request CurrencyConvertRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse date if provided, otherwise use current date
	var conversionDate time.Time
	if request.ConversionDate != "" {
		var err error
		conversionDate, err = time.Parse("2006-01-02", request.ConversionDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
			return
		}
	}

	// Create the query
	conversionQuery := query.ConvertCurrencyQuery{
		Amount:         request.Amount,
		FromCurrency:   request.FromCurrency,
		ToCurrency:     request.ToCurrency,
		ConversionDate: conversionDate,
	}

	// Execute the query
	r := e.convertHandler.Handle(c.Request.Context(), conversionQuery)
	handleResponse(c,
		r,
		withMapping[query.ConvertCurrencyResult](func(c query.ConvertCurrencyResult) any {
			return CurrencyConvertResponse{
				Amount:       c.Amount,
				FromCurrency: c.FromCurrency,
				ToCurrency:   c.ToCurrency,
				Rate:         c.Rate,
				Date:         c.Date,
			}
		}))
}

func (e CurrencyConvertEndpoint) Pattern() []string {
	return []string{
		"/convert",
	}
}

func (e CurrencyConvertEndpoint) Method() string {
	return http.MethodPost
}

func (e CurrencyConvertEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}
