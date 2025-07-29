package persistence

import (
	"github.com/oleexo/subtracker/internal/domain/user"
	"time"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/subscription"
	"github.com/oleexo/subtracker/pkg/slicesx"
)

type subscriptionPayerSqlModel struct {
	baseSqlModel

	FamilyId       uuid.UUID             `gorm:"type:uuid;not null"`
	Family         familySqlModel        `gorm:"foreignKey:FamilyId;references:Id"`
	FamilyMemberId *uuid.UUID            `gorm:"type:uuid"`
	FamilyMember   *familyMemberSqlModel `gorm:"foreignKey:FamilyMemberId;references:Id"`
	SubscriptionId uuid.UUID             `gorm:"type:uuid;not null"`
	Subscription   subscriptionSqlModel  `gorm:"foreignKey:SubscriptionId;references:Id"`
}

func (s subscriptionPayerSqlModel) TableName() string {
	return "subscription_payers"
}

func newSubscriptionPayerSqlModel(subscriptionId uuid.UUID, source subscription.Payer) subscriptionPayerSqlModel {
	model := subscriptionPayerSqlModel{
		baseSqlModel:   newBaseSqlModel(source),
		FamilyId:       source.FamilyId(),
		SubscriptionId: subscriptionId,
	}

	if source.Type() == subscription.FamilyMemberPayer {
		memberId := source.MemberId()
		model.FamilyMemberId = &memberId
	}
	return model
}

type subscriptionServiceUserModel struct {
	FamilyMemberId uuid.UUID            `gorm:"primaryKey;type:uuid;not null"`
	FamilyMember   familyMemberSqlModel `gorm:"foreignKey:FamilyMemberId;references:Id"`
	SubscriptionId uuid.UUID            `gorm:"primaryKey;type:uuid;not null"`
	Subscription   subscriptionSqlModel `gorm:"foreignKey:SubscriptionId;references:Id"`
}

func (s subscriptionServiceUserModel) TableName() string {
	return "subscription_service_users"
}

func newSubscriptionServiceUserModel(familyMemberId uuid.UUID, subscriptionId uuid.UUID) subscriptionServiceUserModel {
	return subscriptionServiceUserModel{
		FamilyMemberId: familyMemberId,
		SubscriptionId: subscriptionId,
	}
}

type subscriptionSqlModel struct {
	baseSqlModel
	baseOwnerSqlModel

	FriendlyName     *string                    `gorm:"type:varchar(100)"`
	FreeTrialDays    *uint                      `gorm:"type:integer;default:0"`
	ProviderId       uuid.UUID                  `gorm:"type:uuid;not null"`
	Provider         providerSqlModel           `gorm:"foreignKey:ProviderId;references:Id"`
	PlanId           uuid.UUID                  `gorm:"type:uuid;not null"`
	Plan             providerPlanSqlModel       `gorm:"foreignKey:PlanId;references:Id"`
	PriceId          uuid.UUID                  `gorm:"type:uuid;not null"`
	Price            providerPriceSqlModel      `gorm:"foreignKey:PriceId;references:Id"`
	PayerId          *uuid.UUID                 `gorm:"type:uuid"`
	Payer            *subscriptionPayerSqlModel `gorm:"foreignKey:SubscriptionId;references:Id"`
	StartDate        time.Time                  `gorm:"type:timestamp;not null"`
	EndDate          *time.Time                 `gorm:"type:timestamp"`
	Recurrency       string                     `gorm:"type:varchar(10);not null"`
	CustomRecurrency *uint                      `gorm:"type:integer"`

	ServiceUsers []subscriptionServiceUserModel
}

func (s subscriptionSqlModel) TableName() string {
	return "subscriptions"
}

func newSubscriptionSqlModel(source subscription.Subscription) subscriptionSqlModel {
	model := subscriptionSqlModel{
		baseSqlModel:     newBaseSqlModel(source),
		FriendlyName:     source.FriendlyName(),
		FreeTrialDays:    source.FreeTrialDays(),
		ProviderId:       source.ServiceProviderId(),
		PlanId:           source.PlanId(),
		PriceId:          source.PriceId(),
		StartDate:        source.StartDate(),
		EndDate:          source.EndDate(),
		Recurrency:       source.Recurrency().String(),
		CustomRecurrency: source.CustomRecurrency(),
	}

	if source.Payer() != nil {
		payerId := source.Payer().Id()
		model.PayerId = &payerId
	}

	model.OwnerType = source.Owner().Type().String()
	switch source.Owner().Type() {
	case user.FamilyOwner:
		familyId := source.Owner().FamilyId()
		model.OwnerFamilyId = &familyId
		break
	case user.PersonalOwner:
		userId := source.Owner().UserId()
		model.OwnerUserId = &userId
		break
	}

	return model
}

func newPayer(source subscriptionPayerSqlModel) subscription.Payer {
	payerType, err := subscription.ParsePayerType(source.Subscription.Recurrency)
	if err != nil {
		panic(err)
	}

	var payer subscription.Payer
	switch payerType {
	case subscription.FamilyPayer:
		payer = subscription.NewFamilyPayer(source.Id,
			source.FamilyId,
			source.CreatedAt,
			source.UpdatedAt)
	case subscription.FamilyMemberPayer:
		if source.FamilyMemberId == nil {
			panic("family member id is nil for family member payer")
		}
		payer = subscription.NewFamilyMemberPayer(source.Id,
			source.FamilyId,
			*source.FamilyMemberId,
			source.CreatedAt,
			source.UpdatedAt)
	default:
		panic("unknown payer type")
	}
	payer.Clean()
	return payer
}

func newSubscription(source subscriptionSqlModel) subscription.Subscription {
	var payer subscription.Payer
	if source.Payer != nil {
		payer = newPayer(*source.Payer)
	}
	var serviceUsers []uuid.UUID
	if source.ServiceUsers != nil && len(source.ServiceUsers) > 0 {
		serviceUsers = slicesx.Map(source.ServiceUsers, func(su subscriptionServiceUserModel) uuid.UUID {
			return su.FamilyMemberId
		})
	}
	recurrency, err := subscription.ParseRecurrencyType(source.Recurrency)
	if err != nil {
		panic(err)
	}

	ownerType, err := user.ParseOwnerType(source.OwnerType)
	if err != nil {
		panic(err)
	}
	owner := user.NewOwner(ownerType, source.OwnerFamilyId, source.OwnerUserId)

	sub := subscription.NewSubscription(
		source.Id,
		source.FriendlyName,
		source.FreeTrialDays,
		source.ProviderId,
		source.PlanId,
		source.PriceId,
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
