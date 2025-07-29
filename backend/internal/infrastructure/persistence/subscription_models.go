package persistence

import (
	"database/sql"
	"time"

	"github.com/oleexo/subtracker/internal/domain/user"

	"github.com/google/uuid"

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

	FriendlyName     sql.NullString        `gorm:"type:varchar(100)"`
	FreeTrialDays    *uint                 `gorm:"type:integer;default:0"`
	ProviderId       uuid.UUID             `gorm:"type:uuid;not null"`
	Provider         ProviderSqlModel      `gorm:"foreignKey:ProviderId;references:Id"`
	PlanId           uuid.UUID             `gorm:"type:uuid;not null"`
	Plan             providerPlanSqlModel  `gorm:"foreignKey:PlanId;references:Id"`
	PriceId          uuid.UUID             `gorm:"type:uuid;not null"`
	Price            providerPriceSqlModel `gorm:"foreignKey:PriceId;references:Id"`
	FamilyId         *uuid.UUID            `gorm:"type:uuid"`
	Family           FamilySqlModel        `gorm:"foreignKey:FamilyId;references:Id"`
	PayerType        sql.NullString        `gorm:"type:varchar(10)"`
	PayerMemberId    *uuid.UUID            `gorm:"type:uuid"`
	PayerMember      FamilyMemberSqlModel  `gorm:"foreignKey:PayerMemberId;references:Id"`
	StartDate        time.Time             `gorm:"type:timestamp;not null"`
	EndDate          *time.Time            `gorm:"type:timestamp"`
	Recurrency       string                `gorm:"type:varchar(10);not null"`
	CustomRecurrency *uint                 `gorm:"type:integer"`

	ServiceUsers []SubscriptionServiceUserModel `gorm:"foreignKey:SubscriptionId;references:Id"`
}

func (s SubscriptionSqlModel) TableName() string {
	return "subscriptions"
}

func newSubscriptionSqlModel(source subscription.Subscription) SubscriptionSqlModel {
	model := SubscriptionSqlModel{
		BaseSqlModel:     newBaseSqlModel(source),
		FriendlyName:     stringToSqlNull(source.FriendlyName()),
		FreeTrialDays:    source.FreeTrialDays(),
		ProviderId:       source.ServiceProviderId(),
		PlanId:           source.PlanId(),
		PriceId:          source.PriceId(),
		StartDate:        source.StartDate(),
		EndDate:          source.EndDate(),
		FamilyId:         source.FamilyId(),
		Recurrency:       source.Recurrency().String(),
		CustomRecurrency: source.CustomRecurrency(),
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
	case user.FamilyOwner:
		familyId := source.Owner().FamilyId()
		model.OwnerFamilyId = &familyId
		break
	case user.PersonalOwner:
		userId := source.Owner().UserId()
		model.OwnerUserId = stringToSqlNull(&userId)
		break
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
	if source.ServiceUsers != nil && len(source.ServiceUsers) > 0 {
		serviceUsers = slicesx.Map(source.ServiceUsers, func(su SubscriptionServiceUserModel) uuid.UUID {
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
	owner := user.NewOwner(ownerType, source.OwnerFamilyId, sqlNullToString(source.OwnerUserId))

	sub := subscription.NewSubscription(
		source.Id,
		sqlNullToString(source.FriendlyName),
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
