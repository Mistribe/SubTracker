package subscription

import (
	"github.com/mistribe/subtracker/internal/domain/currency"
	"github.com/mistribe/subtracker/internal/domain/entity"
	"github.com/mistribe/subtracker/pkg/x/validation"
)

type Price interface {
	entity.ETagEntity

	Amount() currency.Amount
	SetAmount(amount currency.Amount)
	GetValidationErrors() validation.Errors
}

type price struct {
	amount currency.Amount
}

func NewPrice(amount currency.Amount) Price {
	return &price{
		amount: amount,
	}
}

func (c *price) Amount() currency.Amount {
	return c.amount
}

func (c *price) SetAmount(amount currency.Amount) {
	c.amount = amount
}

func (c *price) ETag() string {
	return entity.CalculateETag(c)
}

func (c *price) ETagFields() []interface{} {
	return []interface{}{
		c.amount.Value(),
		c.amount.Currency().String(),
	}
}

func (c *price) GetValidationErrors() validation.Errors {
	var errors validation.Errors

	if c.Amount().Value() <= 0 {
		errors = append(errors, validation.NewError("amount", "Amount must be greater than 0"))
	}

	if errors.HasErrors() {
		return errors
	}

	return nil
}
