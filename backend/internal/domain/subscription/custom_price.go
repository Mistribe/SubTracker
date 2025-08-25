package subscription

import (
	"github.com/mistribe/subtracker/internal/domain/currency"
	"github.com/mistribe/subtracker/internal/domain/entity"
	"github.com/mistribe/subtracker/pkg/x/validation"
)

type CustomPrice interface {
	entity.ETagEntity
	currency.Amount

	SetAmount(amount currency.Amount)
	GetValidationErrors() validation.Errors
}

type customSubscriptionPrice struct {
	currency.Amount
}

func NewCustomPrice(
	amount float64,
	curreny currency.Unit) CustomPrice {
	return &customSubscriptionPrice{
		Amount: currency.NewAmount(amount, curreny),
	}
}

func (c *customSubscriptionPrice) SetAmount(amount currency.Amount) {
	c.Amount = amount
}

func (c *customSubscriptionPrice) ETag() string {
	return entity.CalculateETag(c)
}

func (c *customSubscriptionPrice) ETagFields() []interface{} {
	return []interface{}{
		c.Amount.Value(),
		c.Amount.Currency().String(),
	}
}

func (c *customSubscriptionPrice) GetValidationErrors() validation.Errors {
	var errors validation.Errors

	if c.Amount.Value() <= 0 {
		errors = append(errors, validation.NewError("amount", "Amount must be greater than 0"))
	}

	if errors.HasErrors() {
		return errors
	}

	return nil
}
