package query

import (
	"context"
	"time"

	"github.com/mistribe/subtracker/internal/domain/currency"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

// CurrencyRateQuery represents a query to convert an amount from one currency to another
type CurrencyRateQuery struct {
	ConversionDate time.Time
}

// CurrencyRateResponse represents the result of a currency conversion
type CurrencyRateResponse struct {
	Timestamp time.Time
	Rates     []currency.Rate
}

// CurrencyRateQueryHandler handles currency conversion queries
type CurrencyRateQueryHandler struct {
	currencyRepository currency.Repository
}

func NewCurrencyRateQueryHandler(repository currency.Repository) *CurrencyRateQueryHandler {
	return &CurrencyRateQueryHandler{
		currencyRepository: repository,
	}
}

// Handle processes a CurrencyRateQuery and returns a CurrencyRateResponse
func (h *CurrencyRateQueryHandler) Handle(
	ctx context.Context,
	query CurrencyRateQuery) result.Result[CurrencyRateResponse] {

	rates, err := h.currencyRepository.GetRatesByDate(ctx, query.ConversionDate)
	if err != nil {
		return result.Fail[CurrencyRateResponse](err)
	}

	return result.Success(CurrencyRateResponse{
		Timestamp: query.ConversionDate,
		Rates:     rates,
	})
}
