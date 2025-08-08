package query

import (
	"context"
	"time"

	currencyunit "golang.org/x/text/currency"

	"github.com/oleexo/subtracker/internal/domain/currency"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

// ConvertCurrencyQuery represents a query to convert an amount from one currency to another
type ConvertCurrencyQuery struct {
	Amount         float64
	FromCurrency   string
	ToCurrency     string
	ConversionDate time.Time
}

// ConvertCurrencyResult represents the result of a currency conversion
type ConvertCurrencyResult struct {
	Amount       float64
	FromCurrency string
	ToCurrency   string
	Rate         float64
	Date         time.Time
}

// ConvertCurrencyQueryHandler handles currency conversion queries
type ConvertCurrencyQueryHandler struct {
	repository currency.Repository
}

// NewConvertCurrencyHandler creates a new ConvertCurrencyQueryHandler
func NewConvertCurrencyQueryHandler(repository currency.Repository) *ConvertCurrencyQueryHandler {
	return &ConvertCurrencyQueryHandler{
		repository: repository,
	}
}

// Handle processes a ConvertCurrencyQuery and returns a ConvertCurrencyResult
func (h *ConvertCurrencyQueryHandler) Handle(
	ctx context.Context,
	query ConvertCurrencyQuery) result.Result[ConvertCurrencyResult] {
	// Parse currency units
	fromCurrency, err := currencyunit.ParseISO(query.FromCurrency)
	if err != nil {
		return result.Fail[ConvertCurrencyResult](err)
	}

	toCurrency, err := currencyunit.ParseISO(query.ToCurrency)
	if err != nil {
		return result.Fail[ConvertCurrencyResult](err)
	}

	// Use current date if not specified
	conversionDate := query.ConversionDate
	if conversionDate.IsZero() {
		conversionDate = time.Now()
	}

	// Convert the amount
	convertedAmount, err := h.repository.ConvertAmount(ctx, query.Amount, fromCurrency, toCurrency, conversionDate)
	if err != nil {
		return result.Fail[ConvertCurrencyResult](err)
	}

	// Get the rate used for conversion
	var rate float64
	if fromCurrency == toCurrency {
		rate = 1.0
	} else {
		rateObj, err := h.repository.GetRateByDate(ctx, fromCurrency, toCurrency, conversionDate)
		if err == nil {
			rate = rateObj.ExchangeRate()
		} else {
			// Try reverse rate
			rateObj, err = h.repository.GetRateByDate(ctx, toCurrency, fromCurrency, conversionDate)
			if err == nil {
				rate = 1.0 / rateObj.ExchangeRate()
			}
		}
	}

	return result.Success(ConvertCurrencyResult{
		Amount:       convertedAmount,
		FromCurrency: query.FromCurrency,
		ToCurrency:   query.ToCurrency,
		Rate:         rate,
		Date:         conversionDate,
	})
}
