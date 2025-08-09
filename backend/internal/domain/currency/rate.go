package currency

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/text/currency"

	"github.com/oleexo/subtracker/internal/domain/entity"
)

type Rate interface {
	entity.Entity
	entity.ETagEntity

	FromCurrency() currency.Unit
	ToCurrency() currency.Unit
	RateDate() time.Time
	ExchangeRate() float64
	Equal(other Rate) bool
}

// Rate represents a currency exchange rate at a specific date
type rate struct {
	*entity.Base

	fromCurrency currency.Unit
	toCurrency   currency.Unit
	rateDate     time.Time
	exchangeRate float64
}

// NewRate creates a new currency rate
func NewRate(
	fromCurrency currency.Unit,
	toCurrency currency.Unit,
	rateDate time.Time,
	exchangeRate float64,
	createdAt time.Time,
	updatedAt time.Time) Rate {
	return &rate{
		Base:         entity.NewBase(uuid.New(), createdAt, updatedAt, true, false),
		fromCurrency: fromCurrency,
		toCurrency:   toCurrency,
		rateDate:     rateDate,
		exchangeRate: exchangeRate,
	}
}

// FromCurrency returns the source currency
func (r *rate) FromCurrency() currency.Unit {
	return r.fromCurrency
}

// ToCurrency returns the target currency
func (r *rate) ToCurrency() currency.Unit {
	return r.toCurrency
}

// RateDate returns the date of the exchange rate
func (r *rate) RateDate() time.Time {
	return r.rateDate
}

// ExchangeRate returns the exchange rate value
func (r *rate) ExchangeRate() float64 {
	return r.exchangeRate
}

func (r *rate) Equal(other Rate) bool {
	if other == nil {
		return false
	}

	return r.ETag() == other.ETag()
}

func (r *rate) ETagFields() []interface{} {
	return []interface{}{
		r.fromCurrency.String(),
		r.toCurrency.String(),
		r.rateDate,
		r.exchangeRate,
	}
}
func (r *rate) ETag() string {
	return entity.CalculateETag(r)
}
