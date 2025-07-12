package subscription

import (
	"sort"
	"time"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/pkg/langext/option"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type Payment struct {
	id        uuid.UUID
	price     float64
	startDate time.Time
	endDate   *time.Time
	months    int
	currency  string
	createdAt time.Time
	updatedAt time.Time
	isDirty   bool
}

func NewPayment(
	id uuid.UUID,
	price float64,
	startDate time.Time,
	endDate option.Option[time.Time],
	months int,
	currency string,
	createdAt,
	updatedAt time.Time) result.Result[Payment] {
	payment := NewPaymentWithoutValidation(id, price, startDate, endDate, months, currency, createdAt, updatedAt)
	if err := payment.Validate(); err != nil {
		return result.Fail[Payment](err)
	}
	return result.Success(payment)
}

func NewPaymentWithoutValidation(
	id uuid.UUID,
	price float64,
	startDate time.Time,
	endDate option.Option[time.Time],
	months int,
	currency string,
	createdAt,
	updatedAt time.Time) Payment {
	return Payment{
		id:        id,
		price:     price,
		startDate: startDate.UTC().Truncate(24 * time.Hour),
		endDate: endDate.Transform(func(v time.Time) time.Time {
			return v.UTC().Truncate(24 * time.Hour)
		}).Value(),
		months:    months,
		currency:  currency,
		createdAt: createdAt.UTC(),
		updatedAt: updatedAt.UTC(),
		isDirty:   true,
	}
}

func (p *Payment) Id() uuid.UUID {
	return p.id
}

func (p *Payment) Price() float64 {
	return p.price
}

func (p *Payment) StartDate() time.Time {
	return p.startDate
}

func (p *Payment) EndDate() option.Option[time.Time] {
	return option.New(p.endDate)
}

func (p *Payment) Months() int {
	return p.months
}

func (p *Payment) Currency() string {
	return p.currency
}

func (p *Payment) CreatedAt() time.Time {
	return p.createdAt
}

func (p *Payment) UpdatedAt() time.Time {
	return p.updatedAt
}

func (p *Payment) Validate() error {
	if p.endDate != nil {
		if p.endDate.Before(p.startDate) {
			return ErrPaymentCannotEndBeforeStart
		}
	}

	return nil
}

func (p *Payment) IsDirty() bool {
	return p.isDirty
}

func (p *Payment) Clean() {
	p.isDirty = false
}

func (p *Payment) SetPrice(price float64) {
	p.price = price
	p.isDirty = true
}

func (p *Payment) SetStartDate(date time.Time) {
	p.startDate = date.UTC().Truncate(24 * time.Hour)
	p.isDirty = true
}

func (p *Payment) SetEndDate(date option.Option[time.Time]) {
	p.endDate = date.Transform(func(v time.Time) time.Time {
		return v.UTC().Truncate(24 * time.Hour)
	}).Value()
	p.isDirty = true
}

func (p *Payment) SetMonths(months int) {
	p.months = months
	p.isDirty = true
}

func (p *Payment) SetCurrency(currency string) {
	p.currency = currency
	p.isDirty = true
}

func (p *Payment) SetUpdatedAt(updatedAt time.Time) {
	p.updatedAt = updatedAt
	p.isDirty = true
}

func (p *Payment) Equal(other Payment) bool {
	return p.id == other.id &&
		p.price == other.price &&
		p.startDate == other.startDate &&
		p.endDate == other.endDate &&
		p.months == other.months &&
		p.currency == other.currency &&
		p.createdAt == other.createdAt &&
		p.updatedAt == other.updatedAt
}

func sortPayments(payments []Payment) {
	sort.Slice(payments, func(i, j int) bool {
		return payments[i].startDate.After(payments[j].startDate)
	})

	var previous *Payment
	//var next *Payment
	for idx := 0; idx < len(payments); idx++ {
		current := &payments[idx]
		if idx == 0 {
			previous = nil
		} else {
			previous = &payments[idx-1]
		}
		//if idx+1 < len(payments) {
		//	next = &payments[idx+1]
		//}

		if previous != nil {
			if current.endDate == nil {
				current.endDate = &previous.startDate
			}
			if current.endDate.After(previous.startDate) {
				current.endDate = &previous.startDate
			}
		}

		payments[idx] = *current
	}

}
