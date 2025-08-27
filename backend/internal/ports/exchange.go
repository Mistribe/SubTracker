package ports

import (
	"context"
	"time"

	"github.com/mistribe/subtracker/internal/domain/currency"
)

type Exchange interface {
	ToCurrency(
		ctx context.Context,
		initial currency.Amount,
		target currency.Unit) (currency.Amount, error)
	ToCurrencyAt(
		ctx context.Context,
		initial currency.Amount,
		target currency.Unit,
		at time.Time) (currency.Amount, error)
}
