package subscription

import (
	"time"

	"github.com/mistribe/subtracker/internal/domain/currency"
	"github.com/mistribe/subtracker/internal/domain/entity"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/pkg/slicesx"
	"github.com/mistribe/subtracker/pkg/x/herd"
	"github.com/mistribe/subtracker/pkg/x/validation"
)

const (
	MaxLabelCount                       int = 10
	MaxFamilyMemberPerSubscriptionCount int = 10
)

type Subscription interface {
	entity.Entity[types.SubscriptionID]
	entity.ETagEntity

	// FriendlyName is an override name to the provider name
	FriendlyName() *string
	//FreeTrial is the configuration of the free trial for the subscription
	FreeTrial() FreeTrial
	// ProviderId is the provider of the subscription
	ProviderId() types.ProviderID
	Price() Price
	// Owner is who own the subscription
	Owner() types.Owner
	// Payer is who is paying the subscription
	Payer() Payer
	FamilyUsers() *slicesx.Tracked[types.FamilyMemberID]
	// StartDate is when the subscription as start
	StartDate() time.Time
	// EndDate is when the subscription as stop
	EndDate() *time.Time
	// Recurrency defined the type of recurrency for the subscription
	Recurrency() RecurrencyType
	// CustomRecurrency When the recurrency is custom, this is the number of month between each recurrence
	CustomRecurrency() *int32
	Labels() *slicesx.Tracked[LabelRef]

	// GetRecurrencyAmount returns the amount of the subscription for the given recurrency type
	GetRecurrencyAmount(to RecurrencyType) currency.Amount
	// GetNextRenewalDate returns the next renewal date of the subscription
	GetNextRenewalDate() *time.Time
	// GetPrice returns the price of the subscription from custom or price LabelID
	GetPrice() currency.Amount
	// GetTotalSpent returns the total spent of the subscription
	GetTotalSpent() currency.Amount
	// IsActive returns true if the subscription is active
	IsActive() bool
	// IsStarted returns true if the subscription is started
	IsStarted() bool
	// IsActiveAt determines whether the subscription is active at the given time t.
	IsActiveAt(t time.Time) bool
	// GetTotalDuration returns the total duration of the subscription
	GetTotalDuration() time.Duration

	SetFriendlyName(name *string)
	SetFreeTrial(trial FreeTrial)
	SetPrice(amount currency.Amount)
	SetOwner(owner types.Owner)
	SetPayer(payer Payer)
	SetFamilyUsers(familyMembers []types.FamilyMemberID)
	SetStartDate(startDate time.Time)
	SetEndDate(endDate *time.Time)
	SetRecurrency(recurrency RecurrencyType)
	SetCustomRecurrency(customRecurrency *int32)

	Equal(other Subscription) bool
	GetValidationErrors() validation.Errors
}

type subscription struct {
	*entity.Base[types.SubscriptionID]

	friendlyName     *string
	freeTrial        FreeTrial
	providerId       types.ProviderID
	price            Price
	owner            types.Owner
	payer            Payer
	familyUsers      *slicesx.Tracked[types.FamilyMemberID]
	labels           *slicesx.Tracked[LabelRef]
	startDate        time.Time
	endDate          *time.Time
	recurrency       RecurrencyType
	customRecurrency *int32
}

func NewSubscription(
	id types.SubscriptionID,
	friendlyName *string,
	freeTrial FreeTrial,
	providerId types.ProviderID,
	customPrice Price,
	owner types.Owner,
	payer Payer,
	familyUsers []types.FamilyMemberID,
	labels []LabelRef,
	startDate time.Time,
	endDate *time.Time,
	recurrency RecurrencyType,
	customRecurrency *int32,
	createdAt time.Time,
	updatedAt time.Time,
) Subscription {
	return &subscription{
		Base:             entity.NewBase[types.SubscriptionID](id, createdAt, updatedAt, true, false),
		friendlyName:     friendlyName,
		freeTrial:        freeTrial,
		providerId:       providerId,
		price:            customPrice,
		owner:            owner,
		payer:            payer,
		familyUsers:      slicesx.NewTracked(familyUsers, types.FamilyMemberIdComparer, types.FamilyMemberIdComparer),
		labels:           slicesx.NewTracked(labels, LabelRefUniqueComparer, LabelRefComparer),
		startDate:        startDate,
		endDate:          endDate,
		recurrency:       recurrency,
		customRecurrency: customRecurrency,
	}
}

func (s *subscription) Labels() *slicesx.Tracked[LabelRef] {
	return s.labels
}

func (s *subscription) IsStarted() bool {
	if s.startDate.After(time.Now()) {
		return false
	}
	return true
}

func (s *subscription) GetTotalDuration() time.Duration {
	if s.startDate.After(time.Now()) {
		return time.Duration(0)
	}
	if s.endDate != nil {
		return s.endDate.Sub(s.startDate)
	}

	return time.Now().Sub(s.startDate)
}

func (s *subscription) IsActiveAt(t time.Time) bool {
	if s.startDate.After(t) {
		return false
	}
	if s.endDate != nil {
		return s.endDate.After(t)
	}
	return true
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
	if s.price == nil {
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
	if s.price != nil {
		return currency.NewAmount(s.price.Amount().Value(), s.price.Amount().Currency())
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
	if s.price != nil {
		price = s.price.Amount().Value()
	} else {
		return currency.NewInvalidAmount()
	}

	return currency.NewAmount(price/float64(numberOfMonths), s.price.Amount().Currency())
}

func (s *subscription) GetRecurrencyAmount(to RecurrencyType) currency.Amount {
	if s.price != nil {
		monthlyPrice := s.getMonthlyPrice()

		return currency.NewAmount(
			monthlyPrice.Value()*float64(to.TotalMonths()),
			s.price.Amount().Currency(),
		)
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

func (s *subscription) Price() Price {
	return s.price
}

func (s *subscription) SetPrice(amount currency.Amount) {
	s.price.SetAmount(amount)
	s.SetAsDirty()
}

func (s *subscription) FreeTrial() FreeTrial {
	return s.freeTrial
}

func (s *subscription) ProviderId() types.ProviderID {
	return s.providerId
}

func (s *subscription) Owner() types.Owner {
	return s.owner
}

func (s *subscription) Payer() Payer {
	return s.payer
}

func (s *subscription) FamilyUsers() *slicesx.Tracked[types.FamilyMemberID] {
	return s.familyUsers
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

func (s *subscription) SetOwner(owner types.Owner) {
	s.owner = owner
	s.SetAsDirty()
}

func (s *subscription) SetPayer(payer Payer) {
	s.payer = payer
	s.SetAsDirty()
}

func (s *subscription) SetFamilyUsers(familyMembers []types.FamilyMemberID) {
	s.familyUsers = slicesx.NewTracked(familyMembers, types.FamilyMemberIdComparer, types.FamilyMemberIdComparer)
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

	if s.price != nil {
		fields = append(fields, s.price.ETag())
	}

	if s.payer != nil {
		fields = append(fields, s.payer.ETag())
	}

	if s.familyUsers != nil && s.familyUsers.Len() > 0 {
		fields = append(fields, herd.Select(s.familyUsers.Values(), func(su types.FamilyMemberID) string {
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

func (s *subscription) GetValidationErrors() validation.Errors {
	var errors validation.Errors

	if s.friendlyName != nil && len(*s.friendlyName) > 100 {
		errors = append(errors,
			validation.NewError("friendlyName", "FriendlyName cannot be longer than 100 characters"))
	}

	if s.freeTrial != nil {
		if err := s.freeTrial.GetValidationErrors(); err != nil {
			errors = append(errors, err...)
		}
	}

	if s.price != nil {
		if err := s.price.GetValidationErrors(); err != nil {
			errors = append(errors, err...)
		}
	}

	if s.familyUsers != nil && s.familyUsers.Len() > MaxFamilyMemberPerSubscriptionCount {
		errors = append(errors, validation.NewError("familyUsers", "Number of service users cannot exceed 10"))
	}

	if s.recurrency == CustomRecurrency && s.customRecurrency == nil {
		errors = append(errors,
			validation.NewError("customRecurrency",
				"CustomRecurrency is required when Recurrency is CustomRecurrency"))
	}

	if s.customRecurrency != nil && *s.customRecurrency <= 0 {
		errors = append(errors, validation.NewError("customRecurrency", "CustomRecurrency must be greater than 0"))
	}

	if s.endDate != nil && !s.endDate.IsZero() && s.endDate.Before(s.startDate) {
		errors = append(errors, validation.NewError("endDate", "EndDate must be after StartDate"))
	}

	if errors.HasErrors() {
		return errors
	}

	return nil
}
