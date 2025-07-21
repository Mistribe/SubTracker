package subscription

import (
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/entity"
	"github.com/oleexo/subtracker/pkg/langext/option"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

const (
	MaxLabelCount        int = 10
	MaxFamilyMemberCount int = 10
)

type Subscription struct {
	*entity.Base

	familyId            *uuid.UUID
	name                string
	payments            []Payment
	labels              []uuid.UUID
	familyMembers       []uuid.UUID
	payerId             *uuid.UUID
	payedByJointAccount bool
}

func NewSubscription(
	id uuid.UUID,
	familyId option.Option[uuid.UUID],
	name string,
	payments []Payment,
	labels []uuid.UUID,
	familyMembers []uuid.UUID,
	payer option.Option[uuid.UUID],
	payedByJointAccount bool,
	createdAt time.Time,
	updatedAt time.Time) result.Result[Subscription] {
	sub := NewSubscriptionWithoutValidation(
		id,
		familyId.Value(),
		name,
		payments,
		labels,
		familyMembers,
		payer.Value(),
		payedByJointAccount,
		createdAt,
		updatedAt,
		false,
	)

	if err := sub.Validate(); err != nil {
		return result.Fail[Subscription](err)
	}

	return result.Success(sub)
}

func NewSubscriptionWithoutValidation(
	id uuid.UUID,
	familyId *uuid.UUID,
	name string,
	payments []Payment,
	labels []uuid.UUID,
	familyMembers []uuid.UUID,
	payerId *uuid.UUID,
	payedByJointAccount bool,
	createdAt time.Time,
	updatedAt time.Time,
	isExists bool) Subscription {
	return Subscription{
		Base:                entity.NewBase(id, createdAt, updatedAt, true, isExists),
		familyId:            familyId,
		name:                strings.TrimSpace(name),
		payments:            payments,
		labels:              labels,
		familyMembers:       familyMembers,
		payerId:             payerId,
		payedByJointAccount: payedByJointAccount,
	}
}

func (s *Subscription) Name() string {
	return s.name
}

func (s *Subscription) SetName(name string) {
	s.name = name
}

func (s *Subscription) PayedByJointAccount() bool {
	return s.payedByJointAccount
}

func (s *Subscription) FamilyId() option.Option[uuid.UUID] {
	return option.New(s.familyId)
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
		if p.Id() != paymentId {
			payments = append(payments, p)
		}
	}
	sortPayments(payments)
	s.payments = payments
}

func (s *Subscription) UpdatePayment(payment Payment) {
	payments := make([]Payment, len(s.payments))
	copy(payments, s.payments)
	for i, p := range s.payments {
		if p.Id() == payment.Id() {
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
		if p.Id() == paymentId {
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
	s.SetAsDirty()
}

func (s *Subscription) FamilyMembers() []uuid.UUID {
	return s.familyMembers
}

func (s *Subscription) SetFamilyMembers(familyMembers []uuid.UUID) {
	s.familyMembers = familyMembers
	s.SetAsDirty()
}

func (s *Subscription) Payer() option.Option[uuid.UUID] {
	return option.New(s.payerId)
}

func (s *Subscription) SetPayer(payer option.Option[uuid.UUID]) {
	s.payerId = payer.Value()
	s.SetAsDirty()
}

func (s *Subscription) SetFamilyId(familyId option.Option[uuid.UUID]) {
	s.familyId = familyId.Value()
	s.SetAsDirty()
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

	if s.familyId == nil && len(s.familyMembers) > 0 {
		return ErrCannotHaveFamilyMembersWithoutFamily
	}

	if s.payerId != nil && s.payedByJointAccount {
		return ErrPayerAndJointAccountConflict
	}

	if (s.familyId == nil && s.payedByJointAccount) ||
		(s.familyId == nil && s.payerId != nil) {
		return ErrNoFamilyDefined
	}

	return nil
}

func (s *Subscription) Equal(other Subscription) bool {
	if !s.Base.Equal(*other.Base) {
		return false
	}
	if s.name != other.name {
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

	if s.payerId != other.payerId {
		return false
	}

	return true
}

func (s *Subscription) SetPayedByJointAccount(payedByJointAccount bool) {
	s.payedByJointAccount = payedByJointAccount
	s.SetAsDirty()
}

func (s *Subscription) ETagFields() []interface{} {
	fields := []interface{}{
		s.familyId,
		s.name,
		s.payerId,
		s.payedByJointAccount,
	}

	for _, payment := range s.payments {
		fields = append(fields, payment.ETagFields()...)
	}

	for _, label := range s.labels {
		fields = append(fields, label.String())
	}

	for _, member := range s.familyMembers {
		fields = append(fields, member.String())
	}

	return fields
}

func (s *Subscription) ETag() string {
	return entity.CalculateETag(s, s.Base)
}
