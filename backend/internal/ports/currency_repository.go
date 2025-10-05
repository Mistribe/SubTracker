package ports

import (
	"context"
	"time"

	"github.com/mistribe/subtracker/internal/domain/currency"
	"github.com/mistribe/subtracker/internal/domain/types"
)

type CurrencyRepository interface {
	Repository[types.RateID, currency.Rate]

	GetRateAt(ctx context.Context, from, to currency.Unit, at time.Time) (currency.Rate, error)
	GetRatesByDate(ctx context.Context, date time.Time) (currency.Rates, error)
	GetLatestUpdateDate(ctx context.Context) (time.Time, error)
}
