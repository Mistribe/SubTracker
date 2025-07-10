package subscription

import (
	"sort"
	"time"

	"github.com/oleexo/subtracker/internal/application/core/result"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/application/core/option"
)

const (
	MaxLabelCount        int = 10
	MaxFamilyMemberCount int = 10
)

type Subscription struct {
	id            uuid.UUID
	name          string
	payments      []Payment
	labels        []uuid.UUID
	familyMembers []uuid.UUID
	payer         option.Option[uuid.UUID]
	createdAt     time.Time
	updatedAt     time.Time
	isDirty       bool
}

func NewSubscription(id uuid.UUID,
	name string,
	payments []Payment,
	labels []uuid.UUID,
	familyMembers []uuid.UUID,
	payer option.Option[uuid.UUID],
	createdAt time.Time,
	updatedAt time.Time) result.Result[Subscription] {
	sub := NewSubscriptionWithoutValidation(
		id,
		name,
		payments,
		labels,
		familyMembers,
		payer,
		createdAt,
		updatedAt,
	)

	if err := sub.Validate(); err != nil {
		return result.Fail[Subscription](err)
	}

	return result.Success(sub)
}

func NewSubscriptionWithoutValidation(id uuid.UUID,
	name string,
	payments []Payment,
	labels []uuid.UUID,
	familyMembers []uuid.UUID,
	payer option.Option[uuid.UUID],
	createdAt time.Time,
	updatedAt time.Time) Subscription {
	return Subscription{
		id:            id,
		name:          name,
		payments:      payments,
		labels:        labels,
		familyMembers: familyMembers,
		payer:         payer,
		createdAt:     createdAt,
		updatedAt:     updatedAt,
		isDirty:       true,
	}
}

func (s *Subscription) Id() uuid.UUID {
	return s.id
}

func (s *Subscription) Name() string {
	return s.name
}

func (s *Subscription) SetName(name string) {
	s.name = name
}

func (s *Subscription) Payments() []Payment {
	return s.payments
}

func (s *Subscription) AddPayment(newPayment Payment) {
	payments := append(s.payments, newPayment)
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

	s.payments = payments
}

func (s *Subscription) GetPaymentById(paymentId uuid.UUID) option.Option[Payment] {
	for _, p := range s.payments {
		if p.id == paymentId {
			return option.Some(p)
		}
	}

	return option.None[Payment]()
}

func (s *Subscription) Labels() []uuid.UUID {
	return s.labels
}

func (s *Subscription) SetLabels(labels []uuid.UUID) {
	s.labels = labels
	s.isDirty = true
}

func (s *Subscription) FamilyMembers() []uuid.UUID {
	return s.familyMembers
}

func (s *Subscription) SetFamilyMembers(familyMembers []uuid.UUID) {
	s.familyMembers = familyMembers
	s.isDirty = true
}

func (s *Subscription) Payer() option.Option[uuid.UUID] {
	return s.payer
}

func (s *Subscription) SetPayer(payer option.Option[uuid.UUID]) {
	s.payer = payer
	s.isDirty = true
}
func (s *Subscription) CreatedAt() time.Time {
	return s.createdAt
}

func (s *Subscription) UpdatedAt() time.Time {
	return s.updatedAt
}

func (s *Subscription) SetUpdatedAt(newValue time.Time) {
	s.updatedAt = newValue
	s.isDirty = true
}

func (s *Subscription) IsDirty() bool {
	return s.isDirty
}

func (s *Subscription) Clean() {
	s.isDirty = false
}

func (s *Subscription) Validate() error {
	if len(s.name) < 1 {
		return ErrSubscriptionNameTooShort
	}

	if len(s.payments) == 0 {
		return ErrSubscriptionPaymentMissing
	}

	if len(s.labels) > MaxLabelCount {
		return ErrSubscriptionLabelExcedeed
	}

	if len(s.familyMembers) > MaxFamilyMemberCount {
		return ErrSubscriptionFamilyMemberExcedeed
	}

	for _, payment := range s.payments {
		if err := payment.Validate(); err != nil {
			return err
		}
	}

	return nil
}

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

func NewPayment(id uuid.UUID,
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

func NewPaymentWithoutValidation(id uuid.UUID,
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
		createdAt: createdAt.UTC().Truncate(24 * time.Hour),
		updatedAt: updatedAt,
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
