package currency

import (
	"context"
	"time"

	"golang.org/x/text/currency"

	"github.com/oleexo/subtracker/internal/domain/entity"
)

// Repository defines the interface for currency rate operations
type Repository interface {
	entity.Repository[*Rate]

	// GetRateByDate retrieves a currency rate for a specific currency pair and date
	// If no exact match is found, it returns the closest rate before the specified date
	GetRateByDate(ctx context.Context, fromCurrency, toCurrency currency.Unit, date time.Time) (*Rate, error)

	// GetLatestRate retrieves the most recent currency rate for a specific currency pair
	GetLatestRate(ctx context.Context, fromCurrency, toCurrency currency.Unit) (*Rate, error)

	// ConvertAmount converts an amount from one currency to another at a specific date
	// If no rate is found for the exact date, it uses the closest rate before the specified date
	ConvertAmount(ctx context.Context, amount float64, fromCurrency, toCurrency currency.Unit, date time.Time) (float64, error)
}
