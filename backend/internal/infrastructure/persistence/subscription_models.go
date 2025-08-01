package persistence

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"golang.org/x/text/currency"

	"github.com/oleexo/subtracker/internal/domain/auth"

	"github.com/oleexo/subtracker/internal/domain/subscription"
	"github.com/oleexo/subtracker/pkg/slicesx"
)

type SubscriptionServiceUserModel struct {
	FamilyMemberId uuid.UUID            `gorm:"primaryKey;type:uuid;not null"`
	FamilyMember   FamilyMemberSqlModel `gorm:"foreignKey:FamilyMemberId;references:Id"`
	SubscriptionId uuid.UUID            `gorm:"primaryKey;type:uuid;not null"`
	Subscription   SubscriptionSqlModel `gorm:"foreignKey:SubscriptionId;references:Id"`
}

func (s SubscriptionServiceUserModel) TableName() string {
	return "subscription_service_users"
}

func newSubscriptionServiceUserModel(familyMemberId uuid.UUID, subscriptionId uuid.UUID) SubscriptionServiceUserModel {
	return SubscriptionServiceUserModel{
		FamilyMemberId: familyMemberId,
		SubscriptionId: subscriptionId,
	}
}

type SubscriptionSqlModel struct {
	BaseSqlModel      `gorm:"embedded"`
	BaseOwnerSqlModel `gorm:"embedded"`

	FriendlyName        sql.NullString        `gorm:"type:varchar(100)"`
	FreeTrialStartDate  *time.Time            `gorm:"type:timestamp"`
	FreeTrialEndDate    *time.Time            `gorm:"type:timestamp"`
	ProviderId          uuid.UUID             `gorm:"type:uuid;not null"`
	Provider            ProviderSqlModel      `gorm:"foreignKey:ProviderId;references:Id"`
	PlanId              *uuid.UUID            `gorm:"type:uuid"`
	Plan                providerPlanSqlModel  `gorm:"foreignKey:PlanId;references:Id"`
	PriceId             *uuid.UUID            `gorm:"type:uuid"`
	Price               providerPriceSqlModel `gorm:"foreignKey:PriceId;references:Id"`
	FamilyId            *uuid.UUID            `gorm:"type:uuid"`
	Family              FamilySqlModel        `gorm:"foreignKey:FamilyId;references:Id"`
	PayerType           sql.NullString        `gorm:"type:varchar(10)"`
	PayerMemberId       *uuid.UUID            `gorm:"type:uuid"`
	PayerMember         FamilyMemberSqlModel  `gorm:"foreignKey:PayerMemberId;references:Id"`
	StartDate           time.Time             `gorm:"type:timestamp;not null"`
	EndDate             *time.Time            `gorm:"type:timestamp"`
	Recurrency          string                `gorm:"type:varchar(10);not null"`
	CustomRecurrency    *uint                 `gorm:"type:integer"`
	CustomPriceCurrency *string               `gorm:"type:varchar(3)"`
	CustomPriceAmount   *float64

	ServiceUsers []SubscriptionServiceUserModel `gorm:"foreignKey:SubscriptionId;references:Id"`
}

func (s SubscriptionSqlModel) TableName() string {
	return "subscriptions"
}

func newSubscriptionSqlModel(source subscription.Subscription) SubscriptionSqlModel {
	model := SubscriptionSqlModel{
		BaseSqlModel:     newBaseSqlModel(source, source.ETag()),
		FriendlyName:     stringToSqlNull(source.FriendlyName()),
		ProviderId:       source.ProviderId(),
		PlanId:           source.PlanId(),
		PriceId:          source.PriceId(),
		StartDate:        source.StartDate(),
		EndDate:          source.EndDate(),
		Recurrency:       source.Recurrency().String(),
		CustomRecurrency: source.CustomRecurrency(),
	}

	if source.FreeTrial() != nil {
		startDate := source.FreeTrial().StartDate()
		model.FreeTrialStartDate = &startDate
		endDate := source.FreeTrial().EndDate()
		model.FreeTrialEndDate = &endDate
	}

	if source.CustomPrice() != nil {
		cry := source.CustomPrice().Currency().String()
		amount := source.CustomPrice().Amount()
		model.CustomPriceCurrency = &cry
		model.CustomPriceAmount = &amount
	}

	if source.Payer() != nil {
		payerType := source.Payer().Type().String()
		model.PayerType = stringToSqlNull(&payerType)
		if source.Payer().Type() == subscription.FamilyMemberPayer {
			memberId := source.Payer().MemberId()
			model.PayerMemberId = &memberId
		}
	}

	model.OwnerType = source.Owner().Type().String()
	switch source.Owner().Type() {
	case auth.FamilyOwnerType:
		familyId := source.Owner().FamilyId()
		model.OwnerFamilyId = &familyId
	case auth.PersonalOwnerType:
		userId := source.Owner().UserId()
		model.OwnerUserId = stringToSqlNull(&userId)
	}

	return model
}

func newSubscription(source SubscriptionSqlModel) subscription.Subscription {
	var payer subscription.Payer
	if source.PayerType.Valid {
		payerType, err := subscription.ParsePayerType(source.PayerType.String)
		if err != nil {
			panic(err)
		}
		if source.FamilyId == nil {
			panic("missing family id")
		}
		payer = subscription.NewPayer(payerType, *source.FamilyId, source.PayerMemberId)
	}
	var serviceUsers []uuid.UUID
	if len(source.ServiceUsers) > 0 {
		serviceUsers = slicesx.Select(source.ServiceUsers, func(su SubscriptionServiceUserModel) uuid.UUID {
			return su.FamilyMemberId
		})
	}
	recurrency, err := subscription.ParseRecurrencyType(source.Recurrency)
	if err != nil {
		panic(err)
	}

	ownerType, err := auth.ParseOwnerType(source.OwnerType)
	if err != nil {
		panic(err)
	}
	owner := auth.NewOwner(ownerType, source.OwnerFamilyId, sqlNullToString(source.OwnerUserId))
	var freeTrial subscription.FreeTrial
	if source.FreeTrialEndDate != nil &&
		source.FreeTrialStartDate != nil {
		freeTrial = subscription.NewFreeTrial(
			*source.FreeTrialStartDate,
			*source.FreeTrialEndDate,
		)
	}
	var customPrice subscription.CustomPrice
	if source.CustomPriceAmount != nil &&
		source.CustomPriceCurrency != nil {
		cry := currency.MustParseISO(*source.CustomPriceCurrency)
		customPrice = subscription.NewCustomPrice(
			*source.CustomPriceAmount,
			cry,
		)
	}

	sub := subscription.NewSubscription(
		source.Id,
		sqlNullToString(source.FriendlyName),
		freeTrial,
		source.ProviderId,
		source.PlanId,
		source.PriceId,
		customPrice,
		owner,
		payer,
		serviceUsers,
		source.StartDate,
		source.EndDate,
		recurrency,
		source.CustomRecurrency,
		source.CreatedAt,
		source.UpdatedAt,
	)
	sub.Clean()
	return sub
}
