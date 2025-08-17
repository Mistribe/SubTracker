package subscription

import (
	"github.com/oleexo/subtracker/internal/domain/currency"
	"github.com/oleexo/subtracker/internal/domain/entity"
	"github.com/oleexo/subtracker/pkg/validationx"
)

type CustomPrice interface {
	entity.ETagEntity
	currency.Amount

	SetAmount(amount currency.Amount)
	GetValidationErrors() validationx.Errors
}

type customSubscriptionPrice struct {
	currency.Amount
}

func NewCustomPrice(amount float64,
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

func (c *customSubscriptionPrice) GetValidationErrors() validationx.Errors {
	var errors validationx.Errors

	if c.Amount.Value() <= 0 {
		errors = append(errors, validationx.NewError("amount", "Amount must be greater than 0"))
	}

	if errors.HasErrors() {
		return errors
	}

	return nil
}
