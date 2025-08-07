package subscription

import (
	"golang.org/x/text/currency"

	"github.com/oleexo/subtracker/internal/domain/entity"
	"github.com/oleexo/subtracker/pkg/validationx"
)

type CustomPrice interface {
	entity.ETagEntity

	Currency() currency.Unit
	Amount() float64

	GetValidationErrors() validationx.Errors
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

func (c customSubscriptionPrice) GetValidationErrors() validationx.Errors {
	var errors validationx.Errors

	if c.amount <= 0 {
		errors = append(errors, validationx.NewError("amount", "Amount must be greater than 0"))
	}

	if errors.HasErrors() {
		return errors
	}

	return nil
}
