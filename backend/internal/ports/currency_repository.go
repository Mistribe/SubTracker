package ports

import (
	"context"
	"time"

	"github.com/mistribe/subtracker/internal/domain/currency"
)

// UserRepository defines the interface for currency rate operations
type CurrencyRepository interface {
	Repository[currency.Rate]

	GetRateAt(ctx context.Context, from, to currency.Unit, at time.Time) (currency.Rate, error)
	GetRatesByDate(ctx context.Context, date time.Time) (currency.Rates, error)
	GetLatestUpdateDate(ctx context.Context) (time.Time, error)
}
