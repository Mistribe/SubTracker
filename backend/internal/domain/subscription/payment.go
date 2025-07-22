package subscription

import (
	"sort"
	"time"

	"github.com/google/uuid"
	"golang.org/x/text/currency"

	"github.com/oleexo/subtracker/internal/domain/entity"
	"github.com/oleexo/subtracker/pkg/langext/option"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type Payment struct {
	*entity.Base

	price          float64
	startDate      time.Time
	endDate        *time.Time
	months         int
	currency       currency.Unit
	subscriptionId uuid.UUID
}

func NewPayment(
	id uuid.UUID,
	price float64,
	startDate time.Time,
	endDate option.Option[time.Time],
	months int,
	currency currency.Unit,
	subscriptionId uuid.UUID,
	createdAt,
	updatedAt time.Time) result.Result[Payment] {
	payment := NewPaymentWithoutValidation(id, price, startDate, endDate, months, currency, subscriptionId, createdAt,
		updatedAt, false)
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
	currency currency.Unit,
	subscriptionId uuid.UUID,
	createdAt,
	updatedAt time.Time,
	isExists bool) Payment {
	return Payment{
		Base:      entity.NewBase(id, createdAt, updatedAt, true, isExists),
		price:     price,
		startDate: startDate.UTC().Truncate(24 * time.Hour),
		endDate: endDate.Transform(func(v time.Time) time.Time {
			return v.UTC().Truncate(24 * time.Hour)
		}).Value(),
		months:         months,
		currency:       currency,
		subscriptionId: subscriptionId,
	}
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

func (p *Payment) Currency() currency.Unit {
	return p.currency
}

func (p *Payment) Validate() error {
	if p.endDate != nil {
		if p.endDate.Before(p.startDate) {
			return ErrPaymentCannotEndBeforeStart
		}
	}

	return nil
}

func (p *Payment) SetPrice(price float64) {
	p.price = price
	p.SetAsDirty()
}

func (p *Payment) SetStartDate(date time.Time) {
	p.startDate = date.UTC().Truncate(24 * time.Hour)
	p.SetAsDirty()
}

func (p *Payment) SetEndDate(date option.Option[time.Time]) {
	p.endDate = date.Transform(func(v time.Time) time.Time {
		return v.UTC().Truncate(24 * time.Hour)
	}).Value()
	p.SetAsDirty()
}

func (p *Payment) SetMonths(months int) {
	p.months = months
	p.SetAsDirty()
}

func (p *Payment) SetCurrency(currency currency.Unit) {
	p.currency = currency
	p.SetAsDirty()
}

func (p *Payment) Equal(other Payment) bool {
	return p.Base.Equal(*other.Base) &&
		p.price == other.price &&
		p.startDate == other.startDate &&
		p.endDate == other.endDate &&
		p.months == other.months &&
		p.currency == other.currency
}

func (p *Payment) ETagFields() []interface{} {
	return []interface{}{
		p.price,
		p.startDate,
		p.endDate,
		p.months,
		p.currency,
	}
}
func (p *Payment) ETag() string {
	return entity.CalculateETag(p, p.Base)
}

func (p *Payment) SubscriptionId() uuid.UUID {
	return p.subscriptionId
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
