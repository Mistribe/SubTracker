package subscription

import (
	"time"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/internal/domain/entity"
	"github.com/oleexo/subtracker/pkg/slicesx"
	"github.com/oleexo/subtracker/pkg/validationx"
	"github.com/oleexo/subtracker/pkg/x"
)

const (
	MaxLabelCount                       int = 10
	MaxFamilyMemberPerSubscriptionCount int = 10
)

type Subscription interface {
	entity.Entity
	entity.ETagEntity

	FriendlyName() *string
	FamilyId() *uuid.UUID
	FreeTrialDays() *uint
	ServiceProviderId() uuid.UUID
	PlanId() uuid.UUID
	PriceId() uuid.UUID
	Owner() auth.Owner
	Payer() Payer
	ServiceUsers() *slicesx.Tracked[uuid.UUID]
	StartDate() time.Time
	EndDate() *time.Time
	Recurrency() RecurrencyType
	CustomRecurrency() *uint

	SetFriendlyName(name *string)
	SetFamilyId(*uuid.UUID)
	SetFreeTrialDays(days *uint)
	SetOwner(owner auth.Owner)
	SetPayer(payer Payer)
	SetServiceUsers(familyMembers []uuid.UUID)
	SetStartDate(startDate time.Time)
	SetEndDate(endDate *time.Time)
	SetRecurrency(recurrency RecurrencyType)
	SetCustomRecurrency(customRecurrency *uint)

	Equal(other Subscription) bool
	GetValidationErrors() validationx.Errors
}

type subscription struct {
	*entity.Base

	friendlyName      *string
	familyId          *uuid.UUID
	freeTrialDays     *uint
	serviceProviderId uuid.UUID
	planId            uuid.UUID
	priceId           uuid.UUID
	owner             auth.Owner
	payer             Payer
	serviceUsers      *slicesx.Tracked[uuid.UUID]
	startDate         time.Time
	endDate           *time.Time
	recurrency        RecurrencyType
	customRecurrency  *uint
}

func NewSubscription(
	id uuid.UUID,
	friendlyName *string,
	freeTrialDays *uint,
	serviceProviderId uuid.UUID,
	planId uuid.UUID,
	priceId uuid.UUID,
	owner auth.Owner,
	payer Payer,
	serviceUsers []uuid.UUID,
	startDate time.Time,
	endDate *time.Time,
	recurrency RecurrencyType,
	customRecurrency *uint,
	createdAt time.Time,
	updatedAt time.Time,
) Subscription {
	return &subscription{
		Base:              entity.NewBase(id, createdAt, updatedAt, true, false),
		friendlyName:      friendlyName,
		freeTrialDays:     freeTrialDays,
		serviceProviderId: serviceProviderId,
		planId:            planId,
		priceId:           priceId,
		owner:             owner,
		payer:             payer,
		serviceUsers:      slicesx.NewTracked(serviceUsers, x.UuidUniqueComparer, x.UuidComparer),
		startDate:         startDate,
		endDate:           endDate,
		recurrency:        recurrency,
		customRecurrency:  customRecurrency,
	}
}

func (s *subscription) FriendlyName() *string {
	return s.friendlyName
}

func (s *subscription) FreeTrialDays() *uint {
	return s.freeTrialDays
}

func (s *subscription) FamilyId() *uuid.UUID {
	return s.familyId
}

func (s *subscription) SetFamilyId(id *uuid.UUID) {
	if s.familyId == id {
		return
	}
	s.familyId = id
	s.SetAsDirty()
}

func (s *subscription) ServiceProviderId() uuid.UUID {
	return s.serviceProviderId
}

func (s *subscription) PlanId() uuid.UUID {
	return s.planId
}

func (s *subscription) PriceId() uuid.UUID {
	return s.priceId
}

func (s *subscription) Owner() auth.Owner {
	return s.owner
}

func (s *subscription) Payer() Payer {
	return s.payer
}

func (s *subscription) ServiceUsers() *slicesx.Tracked[uuid.UUID] {
	return s.serviceUsers
}

func (s *subscription) StartDate() time.Time {
	return s.startDate
}

func (s *subscription) EndDate() *time.Time {
	return s.endDate
}

func (s *subscription) Recurrency() RecurrencyType {
	return s.recurrency
}

func (s *subscription) CustomRecurrency() *uint {
	return s.customRecurrency
}

func (s *subscription) SetFriendlyName(name *string) {
	s.friendlyName = name
	s.SetAsDirty()
}

func (s *subscription) SetFreeTrialDays(days *uint) {
	s.freeTrialDays = days
	s.SetAsDirty()
}

func (s *subscription) SetOwner(owner auth.Owner) {
	s.owner = owner
	s.SetAsDirty()
}

func (s *subscription) SetPayer(payer Payer) {
	s.payer = payer
	s.SetAsDirty()
}

func (s *subscription) SetServiceUsers(familyMembers []uuid.UUID) {
	s.serviceUsers = slicesx.NewTracked(familyMembers, x.UuidUniqueComparer, x.UuidComparer)
	s.SetAsDirty()
}

func (s *subscription) SetStartDate(startDate time.Time) {
	s.startDate = startDate
	s.SetAsDirty()
}

func (s *subscription) SetEndDate(endDate *time.Time) {
	s.endDate = endDate
	s.SetAsDirty()
}

func (s *subscription) SetRecurrency(recurrency RecurrencyType) {
	s.recurrency = recurrency
	s.SetAsDirty()
}

func (s *subscription) SetCustomRecurrency(customRecurrency *uint) {
	s.customRecurrency = customRecurrency
	s.SetAsDirty()
}

func (s *subscription) ETagFields() []interface{} {
	fields := []interface{}{
		s.friendlyName,
		s.freeTrialDays,
		s.serviceProviderId,
		s.planId,
		s.priceId,
		s.owner.ETag(),
		s.payer.ETag(),
		s.serviceUsers,
		s.startDate,
		s.endDate,
		s.recurrency.String(),
		s.customRecurrency,
	}

	if s.payer != nil {
		fields = append(fields, s.payer.ETag())
	}

	if s.serviceUsers != nil && s.serviceUsers.Len() > 0 {
		fields = append(fields, slicesx.Map(s.serviceUsers.Values(), func(su uuid.UUID) string {
			return su.String()
		}))
	}

	return fields
}
func (s *subscription) ETag() string {
	return entity.CalculateETag(s, s.Base)
}

func (s *subscription) Equal(other Subscription) bool {
	if other == nil {
		return false
	}

	return s.ETag() == other.ETag()
}

func (s *subscription) GetValidationErrors() validationx.Errors {
	var errors validationx.Errors
	// todo
	return errors
}
