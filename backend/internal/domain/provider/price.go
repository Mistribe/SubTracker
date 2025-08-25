package provider

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/text/currency"

	"github.com/oleexo/subtracker/internal/domain/entity"
	"github.com/oleexo/subtracker/pkg/x/validation"
)

func priceUniqueComparer(p1, p2 Price) bool {
	return p1.Id() == p2.Id()
}

func priceComparer(p1, p2 Price) bool {
	return p1.Equal(p2)
}

type Price interface {
	entity.Entity
	entity.ETagEntity

	StartDate() time.Time
	SetStartDate(startDate time.Time)
	EndDate() *time.Time
	SetEndDate(endDate *time.Time)
	Currency() currency.Unit
	SetCurrency(currency currency.Unit)
	Amount() float64
	SetAmount(amount float64)
	Equal(other Price) bool
	GetValidationErrors() validation.Errors
}

type price struct {
	*entity.Base

	startDate time.Time
	endDate   *time.Time
	currency  currency.Unit
	amount    float64
}

func NewPrice(
	id uuid.UUID,
	startDate time.Time,
	endDate *time.Time,
	currency currency.Unit,
	amount float64,
	createdAt time.Time,
	updatedAt time.Time) Price {
	return &price{
		Base:      entity.NewBase(id, createdAt, updatedAt, true, false),
		startDate: startDate,
		endDate:   endDate,
		currency:  currency,
		amount:    amount,
	}
}

func (p *price) StartDate() time.Time {
	return p.startDate
}

func (p *price) SetStartDate(startDate time.Time) {
	p.startDate = startDate
	p.SetAsDirty()
}

func (p *price) EndDate() *time.Time {
	return p.endDate
}

func (p *price) SetEndDate(endDate *time.Time) {
	p.endDate = endDate
	p.SetAsDirty()
}

func (p *price) Currency() currency.Unit {
	return p.currency
}

func (p *price) SetCurrency(currency currency.Unit) {
	p.currency = currency
	p.SetAsDirty()
}

func (p *price) Amount() float64 {
	return p.amount
}

func (p *price) SetAmount(amount float64) {
	p.amount = amount
	p.SetAsDirty()
}

func (p *price) ETagFields() []interface{} {
	return []interface{}{
		p.startDate,
		p.endDate,
		p.currency.String(),
		p.amount,
	}
}

func (p *price) ETag() string {
	return entity.CalculateETag(p)
}

func (p *price) Equal(other Price) bool {
	if other == nil {
		return false
	}

	return p.ETag() == other.ETag()
}

func (p *price) GetValidationErrors() validation.Errors {
	var errors validation.Errors

	if p.startDate.IsZero() {
		errors = append(errors, validation.NewError(
			"startDate",
			"start date is required",
		))
	}

	if p.amount <= 0 {
		errors = append(errors, validation.NewError(
			"amount",
			"amount must be greater than 0",
		))
	}

	if p.currency.String() == "" {
		errors = append(errors, validation.NewError(
			"currency",
			"currency is required",
		))
	}

	if errors.HasErrors() {
		return errors
	}

	return nil
}
