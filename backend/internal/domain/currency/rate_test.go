package currency_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	xcur "golang.org/x/text/currency"

	dc "github.com/mistribe/subtracker/internal/domain/currency"
	"github.com/mistribe/subtracker/internal/domain/types"
)

func makeRateAt(
	id types.RateID,
	from, to xcur.Unit,
	rateDate time.Time,
	exchange float64,
	createdAt, updatedAt time.Time,
) dc.Rate {
	return dc.NewRate(id, from, to, rateDate, exchange, createdAt, updatedAt)
}

func TestNewRate(t *testing.T) {
	id := types.NewRateID()
	from := xcur.EUR
	to := xcur.USD
	rateDate := time.Date(2024, 12, 25, 10, 9, 8, 7, time.UTC)
	exchange := 1.2345
	createdAt := time.Date(2023, 1, 2, 3, 4, 5, 0, time.UTC)
	updatedAt := time.Date(2023, 6, 7, 8, 9, 10, 0, time.UTC)

	r := makeRateAt(id, from, to, rateDate, exchange, createdAt, updatedAt)

	t.Run("assigns from currency", func(t *testing.T) {
		assert.Equal(t, from, r.FromCurrency())
	})

	t.Run("assigns to currency", func(t *testing.T) {
		assert.Equal(t, to, r.ToCurrency())
	})

	t.Run("assigns rate date", func(t *testing.T) {
		assert.Equal(t, rateDate, r.RateDate())
	})

	t.Run("assigns exchange rate", func(t *testing.T) {
		assert.Equal(t, exchange, r.ExchangeRate())
	})

	t.Run("assigns entity fields", func(t *testing.T) {
		assert.Equal(t, id, r.Id())
		assert.Equal(t, createdAt, r.CreatedAt())
		assert.Equal(t, updatedAt, r.UpdatedAt())
		assert.True(t, r.IsDirty(), "newly created rate should be dirty")
		assert.False(t, r.IsExists(), "newly created rate should not be marked as exists")
	})
}

func TestRate_Equal_NilOtherReturnsFalse(t *testing.T) {
	r := makeRateAt(types.NewRateID(), xcur.EUR, xcur.USD, time.Now(), 1.1, time.Now(), time.Now())
	assert.False(t, r.Equal(nil))
}

func TestRate_Equal_SameFieldsReturnsTrue(t *testing.T) {
	rateDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	r1 := makeRateAt(types.NewRateID(), xcur.GBP, xcur.USD, rateDate, 1.5, time.Now(), time.Now())
	// Different id and timestamps but same fields used in ETag
	r2 := makeRateAt(types.NewRateID(), xcur.GBP, xcur.USD, rateDate, 1.5, time.Now().Add(1*time.Hour),
		time.Now().Add(2*time.Hour))

	assert.True(t, r1.Equal(r2))
}

func TestRate_Equal_DifferentExchangeRateReturnsFalse(t *testing.T) {
	rateDate := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	r1 := makeRateAt(types.NewRateID(), xcur.GBP, xcur.USD, rateDate, 1.5, time.Now(), time.Now())
	r2 := makeRateAt(types.NewRateID(), xcur.GBP, xcur.USD, rateDate, 1.6, time.Now(), time.Now())

	assert.False(t, r1.Equal(r2))
}

func TestRates_FindExchangeRate_FoundReturnsValueAndTrue(t *testing.T) {
	rate := makeRateAt(types.NewRateID(), xcur.EUR, xcur.USD, time.Now(), 1.2, time.Now(), time.Now())
	rates := dc.Rates{rate}

	value, ok := rates.FindExchangeRate(xcur.EUR, xcur.USD)

	assert.True(t, ok)
	assert.Equal(t, 1.2, value)
}

func TestRates_FindExchangeRate_NotFoundReturnsZeroAndFalse(t *testing.T) {
	rate := makeRateAt(types.NewRateID(), xcur.GBP, xcur.USD, time.Now(), 1.3, time.Now(), time.Now())
	rates := dc.Rates{rate}

	value, ok := rates.FindExchangeRate(xcur.EUR, xcur.USD)

	assert.False(t, ok)
	assert.Equal(t, 0.0, value)
}

func TestRates_WithReverse(t *testing.T) {
	id := types.NewRateID()
	rateDate := time.Date(2024, 5, 6, 7, 8, 9, 0, time.UTC)
	createdAt := time.Date(2023, 3, 3, 3, 3, 3, 0, time.UTC)
	updatedAt := time.Date(2023, 4, 4, 4, 4, 4, 0, time.UTC)
	original := makeRateAt(id, xcur.EUR, xcur.USD, rateDate, 2.0, createdAt, updatedAt)
	rates := dc.Rates{original}

	withReverse := rates.WithReverse()

	t.Run("doubles length", func(t *testing.T) {
		assert.Len(t, withReverse, 2)
	})

	t.Run("creates reversed with reciprocal rate", func(t *testing.T) {
		rev := withReverse[1]
		assert.Equal(t, xcur.USD, rev.FromCurrency())
		assert.Equal(t, xcur.EUR, rev.ToCurrency())
		assert.InDelta(t, 0.5, rev.ExchangeRate(), 1e-12)
	})

	t.Run("preserves dates and timestamps", func(t *testing.T) {
		rev := withReverse[1]
		assert.Equal(t, original.RateDate(), rev.RateDate())
		assert.Equal(t, original.CreatedAt(), rev.CreatedAt())
		assert.Equal(t, original.UpdatedAt(), rev.UpdatedAt())
	})

	t.Run("reverse has different LabelID from original", func(t *testing.T) {
		rev := withReverse[1]
		assert.NotEqual(t, original.Id(), rev.Id())
	})
}
