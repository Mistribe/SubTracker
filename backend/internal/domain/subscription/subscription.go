package subscription

import (
	"time"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/pkg/langext/option"
	"github.com/oleexo/subtracker/pkg/langext/result"
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

func NewSubscription(
	id uuid.UUID,
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

func NewSubscriptionWithoutValidation(
	id uuid.UUID,
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
