package currency_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	xcur "golang.org/x/text/currency"

	dc "github.com/oleexo/subtracker/internal/domain/currency"
)

func makeRate(from, to xcur.Unit, exchange float64) dc.Rate {
	return dc.NewRate(uuid.New(), from, to, time.Now(), exchange, time.Now(), time.Now())
}

func makeRates(rates ...dc.Rate) dc.Rates {
	return dc.Rates(rates)
}

func TestNewAmount(t *testing.T) {
	a := dc.NewAmount(123.45, xcur.USD)

	assert.True(t, a.IsValid(), "new amount should be valid")
	assert.Equal(t, 123.45, a.Value(), "value should match the provided one")
	assert.Equal(t, xcur.USD, a.Currency(), "currency should match the provided unit")
}

func TestNewInvalidAmount(t *testing.T) {
	a := dc.NewInvalidAmount()

	assert.False(t, a.IsValid(), "invalid amount should not be valid")
}

func TestAmount_ToCurrency(t *testing.T) {
	t.Run("invalid amount returns itself invalid", func(t *testing.T) {
		a := dc.NewInvalidAmount()
		out := a.ToCurrency(xcur.EUR, nil)

		assert.False(t, out.IsValid(), "ToCurrency on invalid amount should remain invalid")
	})

	t.Run("already in target currency returns same amount", func(t *testing.T) {
		a := dc.NewAmount(7.5, xcur.EUR)
		out := a.ToCurrency(xcur.EUR, nil)

		assert.True(t, out.IsValid(), "result should be valid")
		assert.Equal(t, xcur.EUR, out.Currency(), "currency should remain target")
		assert.Equal(t, 7.5, out.Value(), "value should be unchanged")
	})

	t.Run("missing USD->target rate returns invalid", func(t *testing.T) {
		// Even when amount is already USD, without USD->EUR conversion rate, it should be invalid
		a := dc.NewAmount(10.0, xcur.USD)
		out := a.ToCurrency(xcur.EUR, makeRates())

		assert.False(t, out.IsValid(), "missing USD->target conversion rate should yield invalid amount")
	})

	t.Run("converts from USD to target when USD->target rate exists", func(t *testing.T) {
		a := dc.NewAmount(10.0, xcur.USD)
		rates := makeRates(
			makeRate(xcur.USD, xcur.EUR, 0.8), // In implementation: usdValue / 0.8
		)

		out := a.ToCurrency(xcur.EUR, rates)

		assert.True(t, out.IsValid(), "result should be valid")
		assert.Equal(t, xcur.EUR, out.Currency(), "currency should be target")
		assert.InDelta(t, 12.5, out.Value(), 1e-9, "value should be usdValue divided by USD->target rate")
	})

	t.Run("converts from non-USD via USD when both rates exist", func(t *testing.T) {
		a := dc.NewAmount(10.0, xcur.EUR)
		rates := makeRates(
			makeRate(xcur.EUR, xcur.USD, 1.2), // 10 EUR -> 12 USD
			makeRate(xcur.USD, xcur.GBP, 0.5), // In implementation: usdValue / 0.5
		)

		out := a.ToCurrency(xcur.GBP, rates)

		assert.True(t, out.IsValid(), "result should be valid")
		assert.Equal(t, xcur.GBP, out.Currency(), "currency should be target")
		assert.InDelta(t, 24.0, out.Value(), 1e-9, "value should be (EUR->USD) then divided by (USD->GBP)")
	})
}
