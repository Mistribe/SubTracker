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

	// FriendlyName is an override name to the provider name
	FriendlyName() *string
	//FreeTrial is the configuration of the free trial for the subscription
	FreeTrial() FreeTrial
	// ProviderId is the provider of the subscription
	ProviderId() uuid.UUID
	// PlanId reserved not used for now
	PlanId() *uuid.UUID
	// PriceId reserved not used for now
	PriceId() *uuid.UUID
	// CustomPrice is an override to the price Id
	CustomPrice() CustomPrice
	// Owner is who own the subscription
	Owner() auth.Owner
	// Payer is who is paying the subscription
	Payer() Payer
	// Service users is the member of the family using the service
	ServiceUsers() *slicesx.Tracked[uuid.UUID]
	// StartDate is when the subscription as start
	StartDate() time.Time
	// EndDate is when the subscription as stop
	EndDate() *time.Time
	// Recurrency defined the type of recurrency for the subscription
	Recurrency() RecurrencyType
	// CustomRecurrency When the recurrency is custom, this is the number of month between each recurrence
	CustomRecurrency() *int32

	SetFriendlyName(name *string)
	SetFreeTrial(trial FreeTrial)
	SetCustomPrice(price CustomPrice)
	SetOwner(owner auth.Owner)
	SetPayer(payer Payer)
	SetServiceUsers(familyMembers []uuid.UUID)
	SetStartDate(startDate time.Time)
	SetEndDate(endDate *time.Time)
	SetRecurrency(recurrency RecurrencyType)
	SetCustomRecurrency(customRecurrency *int32)

	Equal(other Subscription) bool
	GetValidationErrors() validationx.Errors
}

type subscription struct {
	*entity.Base

	friendlyName     *string
	freeTrial        FreeTrial
	providerId       uuid.UUID
	planId           *uuid.UUID
	priceId          *uuid.UUID
	customPrice      CustomPrice
	owner            auth.Owner
	payer            Payer
	serviceUsers     *slicesx.Tracked[uuid.UUID]
	startDate        time.Time
	endDate          *time.Time
	recurrency       RecurrencyType
	customRecurrency *int32
}

func NewSubscription(
	id uuid.UUID,
	friendlyName *string,
	freeTrial FreeTrial,
	providerId uuid.UUID,
	planId *uuid.UUID,
	priceId *uuid.UUID,
	customPrice CustomPrice,
	owner auth.Owner,
	payer Payer,
	serviceUsers []uuid.UUID,
	startDate time.Time,
	endDate *time.Time,
	recurrency RecurrencyType,
	customRecurrency *int32,
	createdAt time.Time,
	updatedAt time.Time,
) Subscription {
	return &subscription{
		Base:             entity.NewBase(id, createdAt, updatedAt, true, false),
		friendlyName:     friendlyName,
		freeTrial:        freeTrial,
		providerId:       providerId,
		planId:           planId,
		priceId:          priceId,
		customPrice:      customPrice,
		owner:            owner,
		payer:            payer,
		serviceUsers:     slicesx.NewTracked(serviceUsers, x.UuidUniqueComparer, x.UuidComparer),
		startDate:        startDate,
		endDate:          endDate,
		recurrency:       recurrency,
		customRecurrency: customRecurrency,
	}
}

func (s *subscription) FriendlyName() *string {
	return s.friendlyName
}

func (s *subscription) CustomPrice() CustomPrice {
	return s.customPrice
}

func (s *subscription) SetCustomPrice(price CustomPrice) {
	s.customPrice = price
	s.SetAsDirty()
}

func (s *subscription) FreeTrial() FreeTrial {
	return s.freeTrial
}

func (s *subscription) ProviderId() uuid.UUID {
	return s.providerId
}

func (s *subscription) PlanId() *uuid.UUID {
	return s.planId
}

func (s *subscription) PriceId() *uuid.UUID {
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

func (s *subscription) CustomRecurrency() *int32 {
	return s.customRecurrency
}

func (s *subscription) SetFriendlyName(name *string) {
	s.friendlyName = name
	s.SetAsDirty()
}

func (s *subscription) SetFreeTrial(trial FreeTrial) {
	s.freeTrial = trial
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

func (s *subscription) SetCustomRecurrency(customRecurrency *int32) {
	s.customRecurrency = customRecurrency
	s.SetAsDirty()
}

func (s *subscription) ETagFields() []interface{} {
	fields := []interface{}{
		s.friendlyName,
		s.providerId,
		s.planId,
		s.priceId,
		s.owner.ETag(),
		s.startDate,
		s.endDate,
		s.recurrency.String(),
	}

	if s.freeTrial != nil {
		fields = append(fields, s.freeTrial.ETag())
	}

	if s.customRecurrency != nil {
		fields = append(fields, *s.customRecurrency)
	}

	if s.customPrice != nil {
		fields = append(fields, s.customPrice.ETag())
	}

	if s.payer != nil {
		fields = append(fields, s.payer.ETag())
	}

	if s.serviceUsers != nil && s.serviceUsers.Len() > 0 {
		fields = append(fields, slicesx.Select(s.serviceUsers.Values(), func(su uuid.UUID) string {
			return su.String()
		}))
	}

	return fields
}
func (s *subscription) ETag() string {
	return entity.CalculateETag(s)
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
