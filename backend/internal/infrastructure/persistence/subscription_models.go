package persistence

import (
	"time"

	"github.com/google/uuid"
)

type subscriptionPayerSqlModel struct {
	baseSqlModel

	FamilyId       *uuid.UUID            `gorm:"type:uuid"`
	Family         *familySqlModel       `gorm:"foreignKey:FamilyId;references:Id"`
	FamilyMemberId *uuid.UUID            `gorm:"type:uuid"`
	FamilyMember   *familyMemberSqlModel `gorm:"foreignKey:FamilyMemberId;references:Id"`
	SubscriptionId uuid.UUID             `gorm:"type:uuid;not null"`
	Subscription   subscriptionSqlModel  `gorm:"foreignKey:SubscriptionId;references:Id"`
}

func (s subscriptionPayerSqlModel) TableName() string {
	return "subscription_payers"
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

type subscriptionSqlModel struct {
	baseSqlModel

	FriendlyName     *string                    `gorm:"type:varchar(100)"`
	FreeTrialDays    *uint                      `gorm:"type:integer;default:0"`
	ProviderId       uuid.UUID                  `gorm:"type:uuid;not null"`
	Provider         providerSqlModel           `gorm:"foreignKey:ProviderId;references:Id"`
	PlanId           uuid.UUID                  `gorm:"type:uuid;not null"`
	Plan             providerPlanSqlModel       `gorm:"foreignKey:PlanId;references:Id"`
	PriceId          *uuid.UUID                 `gorm:"type:uuid"`
	Price            providerPriceSqlModel      `gorm:"foreignKey:PriceId;references:Id"`
	OwnerId          *uuid.UUID                 `gorm:"type:uuid"`
	Owner            *ownerSqlModel             `gorm:"foreignKey:OwnerId;references:Id"`
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
