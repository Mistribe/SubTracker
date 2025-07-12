package subscription

import (
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
	sortPayments(payments)
	s.payments = payments
}

func (s *Subscription) RemovePayment(paymentId uuid.UUID) {
	var payments []Payment
	for _, p := range s.payments {
		if p.id != paymentId {
			payments = append(payments, p)
		}
	}
	sortPayments(payments)
	s.payments = payments
}

func (s *Subscription) UpdatePayment(payment Payment) {
	var payments []Payment
	for i, p := range s.payments {
		if p.id == payment.id {
			payments[i] = payment
		} else {
			payments[i] = p
		}
	}
	sortPayments(payments)
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
		return ErrSubscriptionLabelExceeded
	}

	if len(s.familyMembers) > MaxFamilyMemberCount {
		return ErrSubscriptionFamilyMemberExceeded
	}

	for _, payment := range s.payments {
		if err := payment.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (s *Subscription) Equal(other Subscription) bool {
	if s.id != other.id || s.name != other.name {
		return false
	}

	if len(s.payments) != len(other.payments) {
		return false
	}
	for i := range s.payments {
		if s.payments[i] != other.payments[i] {
			return false
		}
	}

	if len(s.labels) != len(other.labels) {
		return false
	}
	for i := range s.labels {
		if s.labels[i] != other.labels[i] {
			return false
		}
	}

	if len(s.familyMembers) != len(other.familyMembers) {
		return false
	}
	for i := range s.familyMembers {
		if s.familyMembers[i] != other.familyMembers[i] {
			return false
		}
	}

	if !s.payer.Equal(other.payer) {
		return false
	}

	if !s.createdAt.Equal(other.createdAt) || !s.updatedAt.Equal(other.updatedAt) {
		return false
	}

	return true
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
