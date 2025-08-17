package subscription

import (
	"time"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/internal/domain/currency"
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

	// GetRecurrencyAmount returns the amount of the subscription for the given recurrency type
	GetRecurrencyAmount(to RecurrencyType) currency.Amount
	// GetNextRenewalDate returns the next renewal date of the subscription
	GetNextRenewalDate() *time.Time
	// GetPrice returns the price of the subscription from custom or price Id
	GetPrice() currency.Amount
	// GetTotalSpent returns the total spent of the subscription
	GetTotalSpent() currency.Amount
	// IsActive returns true if the subscription is active
	IsActive() bool

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

func (s *subscription) IsActive() bool {
	if s.startDate.After(time.Now()) {
		return false
	}
	if s.endDate != nil {
		return s.endDate.After(time.Now())
	}

	return true
}

func (s *subscription) GetTotalSpent() currency.Amount {
	if s.customPrice == nil && s.priceId == nil {
		return currency.NewInvalidAmount()
	}

	price := s.GetRecurrencyAmount(MonthlyRecurrency)
	if !price.IsValid() {
		return currency.NewInvalidAmount()
	}

	now := time.Now()
	endDate := now
	if s.endDate != nil {
		endDate = *s.endDate
	}

	if s.startDate.After(now) {
		return currency.NewAmount(0, price.Currency())
	}

	months := monthsBetweenCalendar(s.startDate, endDate)
	if months <= 0 {
		return currency.NewAmount(0, price.Currency())
	}

	return currency.NewAmount(price.Value()*float64(months), price.Currency())
}

func (s *subscription) GetPrice() currency.Amount {
	if s.customPrice != nil {
		return currency.NewAmount(s.customPrice.Value(), s.customPrice.Currency())
	}

	if s.priceId != nil {
		panic("not implemented yet")
	}

	return currency.NewInvalidAmount()
}

func (s *subscription) getMonths() int {
	switch s.recurrency {
	case CustomRecurrency:
		if s.customRecurrency != nil {
			return int(*s.customRecurrency)
		}
		return 0
	case OneTimeRecurrency:
		var endDate time.Time
		startDate := s.startDate
		if s.endDate != nil {
			endDate = *s.endDate
		} else {
			endDate = time.Now()
		}
		return monthsBetweenCalendar(startDate, endDate)
	default:
		return s.recurrency.TotalMonths()
	}
}

func (s *subscription) getMonthlyPrice() currency.Amount {
	var price float64
	numberOfMonths := s.getMonths()
	if s.customPrice != nil {
		price = s.customPrice.Value()
	} else if s.priceId != nil {
		panic("not implemented yet")
	}

	return currency.NewAmount(price/float64(numberOfMonths), s.customPrice.Currency())
}

func (s *subscription) GetRecurrencyAmount(to RecurrencyType) currency.Amount {
	if s.customPrice != nil {
		monthlyPrice := s.getMonthlyPrice()

		return currency.NewAmount(
			monthlyPrice.Value()*float64(to.TotalMonths()),
			s.customPrice.Currency(),
		)
	}

	if s.priceId != nil {
		panic("not implemented yet")
	}

	return currency.NewInvalidAmount()
}

func (s *subscription) GetNextRenewalDate() *time.Time {
	if s.endDate != nil {
		if s.endDate.Before(time.Now()) {
			return nil
		}
		return s.endDate
	}

	months := s.getMonthsUntilNextRenewal()
	if months == 0 {
		return nil
	}

	now := time.Now()

	// If the subscription starts in the future, next renewal is the start date
	if s.startDate.After(now) {
		next := s.startDate
		return &next
	}

	// Total whole months between start and now (may be adjusted below)
	diffMonths := (now.Year()-s.startDate.Year())*12 + int(now.Month()-s.startDate.Month())

	// Align start by diffMonths and adjust to ensure aligned date is <= now.
	// This avoids partial-month overcounting.
	aligned := s.startDate.AddDate(0, diffMonths, 0)
	if aligned.After(now) {
		diffMonths--
		aligned = s.startDate.AddDate(0, diffMonths, 0)
	}

	// Compute the minimal k such that start + k*months >= now
	// This is a ceiling division: k = ceil(diffMonths / months)
	k := diffMonths / months
	if diffMonths%months != 0 {
		k++
	}

	next := s.startDate.AddDate(0, k*months, 0)
	// Safety: in rare edge cases (e.g., time-of-day differences), ensure next is not in the past
	if next.Before(now) {
		next = next.AddDate(0, months, 0)
	}
	return &next
}

func (s *subscription) getMonthsUntilNextRenewal() int {
	switch s.recurrency {
	case MonthlyRecurrency:
		return 1
	case QuarterlyRecurrency:
		return 3
	case HalfYearlyRecurrency:
		return 6
	case YearlyRecurrency:
		return 12
	case CustomRecurrency:
		if s.customRecurrency != nil {
			return int(*s.customRecurrency)
		}
		return 0
	default:
		return 0
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

	if s.friendlyName != nil && len(*s.friendlyName) > 100 {
		errors = append(errors,
			validationx.NewError("friendlyName", "FriendlyName cannot be longer than 100 characters"))
	}

	if s.freeTrial != nil {
		if err := s.freeTrial.GetValidationErrors(); err != nil {
			errors = append(errors, err...)
		}
	}

	if s.customPrice != nil {
		if err := s.customPrice.GetValidationErrors(); err != nil {
			errors = append(errors, err...)
		}
	}

	if s.serviceUsers != nil && s.serviceUsers.Len() > MaxFamilyMemberPerSubscriptionCount {
		errors = append(errors, validationx.NewError("serviceUsers", "Number of service users cannot exceed 10"))
	}

	if s.recurrency == CustomRecurrency && s.customRecurrency == nil {
		errors = append(errors,
			validationx.NewError("customRecurrency",
				"CustomRecurrency is required when Recurrency is CustomRecurrency"))
	}

	if s.customRecurrency != nil && *s.customRecurrency <= 0 {
		errors = append(errors, validationx.NewError("customRecurrency", "CustomRecurrency must be greater than 0"))
	}

	if s.endDate != nil && !s.endDate.IsZero() && s.endDate.Before(s.startDate) {
		errors = append(errors, validationx.NewError("endDate", "EndDate must be after StartDate"))
	}

	if errors.HasErrors() {
		return errors
	}

	return nil
}
