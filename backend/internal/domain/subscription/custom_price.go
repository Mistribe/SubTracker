package subscription

import (
	"github.com/mistribe/subtracker/internal/domain/currency"
	"github.com/mistribe/subtracker/internal/domain/entity"
	"github.com/mistribe/subtracker/pkg/x/validation"
)

type CustomPrice interface {
	entity.ETagEntity

	Amount() currency.Amount
	SetAmount(amount currency.Amount)
	GetValidationErrors() validation.Errors
}

type customSubscriptionPrice struct {
	amount currency.Amount
}

func NewCustomPrice(
	amount float64,
	unit currency.Unit) CustomPrice {
	return &customSubscriptionPrice{
		amount: currency.NewAmount(amount, unit),
	}
}

func (c *customSubscriptionPrice) Amount() currency.Amount {
	return c.amount
}

func (c *customSubscriptionPrice) SetAmount(amount currency.Amount) {
	c.amount = amount
}

func (c *customSubscriptionPrice) ETag() string {
	return entity.CalculateETag(c)
}

func (c *customSubscriptionPrice) ETagFields() []interface{} {
	return []interface{}{
		c.amount.Value(),
		c.amount.Currency().String(),
	}
}

func (c *customSubscriptionPrice) GetValidationErrors() validation.Errors {
	var errors validation.Errors

	if c.Amount().Value() <= 0 {
		errors = append(errors, validation.NewError("amount", "Amount must be greater than 0"))
	}

	if errors.HasErrors() {
		return errors
	}

	return nil
}
