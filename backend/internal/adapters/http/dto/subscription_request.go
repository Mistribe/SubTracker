package dto

import (
	"time"
)

// SubscriptionSummaryRequest represents a request for fetching subscription summary details such as costs and renewals.
type SubscriptionSummaryRequest struct {
	TopProviders     uint8 `json:"top_providers,omitempty" form:"top_providers" example:"5" description:"Number of top providers to return"`
	UpcomingRenewals uint8 `json:"upcoming_renewals,omitempty" form:"upcoming_renewals" example:"3" description:"Number of upcoming renewals to return"`
	TopLabels        uint8 `json:"top_labels,omitempty" form:"top_labels"  example:"5" description:"Number of top labels to return"`
	TotalMonthly     bool  `json:"total_monthly,omitempty" form:"total_monthly"  example:"true" description:"Include monthly total costs"`
	TotalYearly      bool  `json:"total_yearly,omitempty" form:"total_yearly"  example:"true" description:"Include yearly total costs"`
}

type CreateSubscriptionRequest struct {
	Id               *string                         `json:"id,omitempty"`
	FriendlyName     *string                         `json:"friendly_name,omitempty"`
	FreeTrial        *SubscriptionFreeTrialModel     `json:"free_trial,omitempty"`
	ProviderId       *string                         `json:"provider_id" binding:"required"`
	ProviderKey      *string                         `json:"provider_key,omitempty"`
	Price            AmountModel                     `json:"price"`
	FamilyUsers      []string                        `json:"family_users,omitempty"`
	Labels           []string                        `json:"labels,omitempty"`
	StartDate        time.Time                       `json:"start_date" binding:"required" format:"date-time"`
	EndDate          *time.Time                      `json:"end_date,omitempty" format:"date-time"`
	Recurrency       string                          `json:"recurrency" binding:"required"`
	CustomRecurrency *int32                          `json:"custom_recurrency,omitempty"`
	Payer            *EditableSubscriptionPayerModel `json:"payer,omitempty"`
	Owner            string                          `json:"owner" binding:"required" example:"personal" enums:"personal,family,system"`
	CreatedAt        *time.Time                      `json:"created_at,omitempty"`
}

type UpdateSubscriptionRequest struct {
	FriendlyName     *string                         `json:"friendly_name,omitempty"`
	FreeTrial        *SubscriptionFreeTrialModel     `json:"free_trial,omitempty"`
	ProviderId       *string                         `json:"provider_id" binding:"required"`
	ProviderKey      *string                         `json:"provider_key,omitempty"`
	Price            AmountModel                     `json:"price"`
	ServiceUsers     []string                        `json:"service_users,omitempty"`
	Labels           []string                        `json:"labels,omitempty"`
	StartDate        time.Time                       `json:"start_date" binding:"required" format:"date-time"`
	EndDate          *time.Time                      `json:"end_date,omitempty" format:"date-time"`
	Recurrency       string                          `json:"recurrency" binding:"required"`
	CustomRecurrency *int32                          `json:"custom_recurrency,omitempty"`
	Payer            *EditableSubscriptionPayerModel `json:"payer,omitempty"`
	Owner            string                          `json:"owner" binding:"required" example:"personal" enums:"personal,family,system"`
	UpdatedAt        *time.Time                      `json:"updated_at,omitempty" format:"date-time"`
}
