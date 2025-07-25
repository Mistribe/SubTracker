package subscription

import (
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/entity"
	"github.com/oleexo/subtracker/pkg/langext/option"
	"github.com/oleexo/subtracker/pkg/langext/result"
	"github.com/oleexo/subtracker/pkg/slicesx"
)

const (
	MaxLabelCount        int = 10
	MaxFamilyMemberCount int = 10
)

type Subscription struct {
	*entity.Base

	name                string
	payments            *slicesx.Tracked[Payment]
	labels              *slicesx.Tracked[uuid.UUID]
	familyMembers       *slicesx.Tracked[uuid.UUID]
	familyId            *uuid.UUID
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

func paymentUniqueComparer(a Payment, b Payment) bool {
	return a.Id() == b.Id()
}

func paymentComparer(a Payment, b Payment) bool {
	return a.Equal(b)
}

func paymentSorter(a Payment, b Payment) bool {
	return a.startDate.After(b.startDate)
}

func uuidUniqueComparer(a uuid.UUID, b uuid.UUID) bool {
	return a == b
}

func uuidComparer(a uuid.UUID, b uuid.UUID) bool {
	return a == b
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
		payments:            slicesx.NewTracked(payments, paymentUniqueComparer, paymentComparer),
		labels:              slicesx.NewTracked(labels, uuidUniqueComparer, uuidComparer),
		familyMembers:       slicesx.NewTracked(familyMembers, uuidUniqueComparer, uuidComparer),
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

func (s *Subscription) Payments() *slicesx.Tracked[Payment] {
	return s.payments
}

func (s *Subscription) AddPayment(newPayment Payment) {
	s.payments.Add(newPayment)
	s.payments.Sort(paymentSorter)
	s.SetAsDirty()
}

func (s *Subscription) RemovePayment(paymentId uuid.UUID) {
	for p := range s.payments.It() {
		if p.Id() != paymentId {
			s.payments.Remove(p)
			break
		}
	}
	s.payments.Sort(paymentSorter)
	s.SetAsDirty()
}

func (s *Subscription) UpdatePayment(payment Payment) {
	s.payments.Update(payment)
	s.payments.Sort(paymentSorter)
	s.SetAsDirty()
}

func (s *Subscription) GetPaymentById(paymentId uuid.UUID) option.Option[Payment] {
	for p := range s.payments.It() {
		if p.Id() == paymentId {
			return option.Some(p)
		}
	}

	return option.None[Payment]()
}

func (s *Subscription) Labels() *slicesx.Tracked[uuid.UUID] {
	return s.labels
}

func (s *Subscription) SetLabels(labels []uuid.UUID) {
	s.labels.Set(labels)
	s.SetAsDirty()
}

func (s *Subscription) FamilyMembers() *slicesx.Tracked[uuid.UUID] {
	return s.familyMembers
}

func (s *Subscription) SetFamilyMembers(familyMembers []uuid.UUID) {
	s.familyMembers.Set(familyMembers)
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

	if s.payments.Len() == 0 {
		return ErrSubscriptionPaymentMissing
	}

	if s.labels.Len() > MaxLabelCount {
		return ErrSubscriptionLabelExceeded
	}

	if s.familyMembers.Len() > MaxFamilyMemberCount {
		return ErrSubscriptionFamilyMemberExceeded
	}

	for payment := range s.payments.It() {
		if err := payment.Validate(); err != nil {
			return err
		}
	}

	if s.familyId == nil && s.familyMembers.Len() > 0 {
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

	if !s.payments.Equals(other.payments.Values()) {
		return false
	}

	if !s.labels.Equals(other.labels.Values()) {
		return false
	}

	if s.familyMembers.Equals(other.familyMembers.Values()) {
		return false
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

	for payment := range s.payments.It() {
		fields = append(fields, payment.ETagFields()...)
	}

	for label := range s.labels.It() {
		fields = append(fields, label.String())
	}

	for member := range s.familyMembers.It() {
		fields = append(fields, member.String())
	}

	return fields
}

func (s *Subscription) ETag() string {
	return entity.CalculateETag(s, s.Base)
}
