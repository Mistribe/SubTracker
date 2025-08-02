package subscription

import (
	"golang.org/x/text/currency"

	"github.com/oleexo/subtracker/internal/domain/entity"
)

type CustomPrice interface {
	entity.ETagEntity

	Currency() currency.Unit
	Amount() float64
}

type customSubscriptionPrice struct {
	amount   float64
	currency currency.Unit
}

func NewCustomPrice(amount float64,
	curreny currency.Unit) CustomPrice {
	return &customSubscriptionPrice{
		amount:   amount,
		currency: curreny,
	}
}

func (c customSubscriptionPrice) Amount() float64 {
	return c.amount
}

func (c customSubscriptionPrice) Currency() currency.Unit {
	return c.currency
}

func (c customSubscriptionPrice) ETag() string {
	return entity.CalculateETag(c)
}

func (c customSubscriptionPrice) ETagFields() []interface{} {
	return []interface{}{
		c.amount,
		c.currency.String(),
	}
}
