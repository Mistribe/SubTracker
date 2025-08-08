package currency

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/text/currency"

	"github.com/oleexo/subtracker/internal/domain/entity"
)

// Rate represents a currency exchange rate at a specific date
type Rate struct {
	*entity.Base
	fromCurrency currency.Unit
	toCurrency   currency.Unit
	rateDate     time.Time
	exchangeRate float64
	etag         string
}

// NewRate creates a new currency rate
func NewRate(
	fromCurrency currency.Unit,
	toCurrency currency.Unit,
	rateDate time.Time,
	exchangeRate float64) *Rate {
	now := time.Now().UTC()
	return &Rate{
		Base:         entity.NewBase(uuid.New(), now, now, true, false),
		fromCurrency: fromCurrency,
		toCurrency:   toCurrency,
		rateDate:     rateDate,
		exchangeRate: exchangeRate,
		etag:         uuid.NewString(),
	}
}

// FromExisting creates a Rate from existing data
func FromExisting(
	id uuid.UUID,
	fromCurrency currency.Unit,
	toCurrency currency.Unit,
	rateDate time.Time,
	exchangeRate float64,
	createdAt time.Time,
	updatedAt time.Time,
	etag string) *Rate {
	return &Rate{
		Base:         entity.NewBase(id, createdAt, updatedAt, false, true),
		fromCurrency: fromCurrency,
		toCurrency:   toCurrency,
		rateDate:     rateDate,
		exchangeRate: exchangeRate,
		etag:         etag,
	}
}

// FromCurrency returns the source currency
func (r *Rate) FromCurrency() currency.Unit {
	return r.fromCurrency
}

// ToCurrency returns the target currency
func (r *Rate) ToCurrency() currency.Unit {
	return r.toCurrency
}

// RateDate returns the date of the exchange rate
func (r *Rate) RateDate() time.Time {
	return r.rateDate
}

// ExchangeRate returns the exchange rate value
func (r *Rate) ExchangeRate() float64 {
	return r.exchangeRate
}

// ETag returns the entity tag
func (r *Rate) ETag() string {
	return r.etag
}

// SetExchangeRate updates the exchange rate value
func (r *Rate) SetExchangeRate(value float64) {
	if r.exchangeRate != value {
		r.exchangeRate = value
		r.updateModified()
	}
}

// updateModified updates the modification tracking fields
func (r *Rate) updateModified() {
	r.SetUpdatedAt(time.Now().UTC())
	r.etag = uuid.NewString()
}

// Ensure Rate implements the Entity interface
var _ entity.Entity = (*Rate)(nil)
