package dto

import (
	"time"
)

// SubscriptionSummaryTopProviderResponse represents the response structure for a top provider in subscription summary.
type SubscriptionSummaryTopProviderResponse struct {
	ProviderId string      `json:"provider_id" binding:"required"`
	Total      AmountModel `json:"total"`
	Duration   string      `json:"duration"`
}

// SubscriptionSummaryTopLabelResponse represents the response structure for a top label in the subscription summary.
type SubscriptionSummaryTopLabelResponse struct {
	LabelId string      `json:"label_id" binding:"required"`
	Total   AmountModel `json:"total"`
}

// SubscriptionSummaryUpcomingRenewalResponse represents the structure for upcoming subscription renewal details.
type SubscriptionSummaryUpcomingRenewalResponse struct {
	ProviderId string       `json:"provider_id" binding:"required"`
	At         time.Time    `json:"at" binding:"required" format:"date-time"`
	Total      AmountModel  `json:"total"`
	Source     *AmountModel `json:"source"`
}

// SubscriptionSummaryResponse represents a summary of subscription details including costs, providers, labels, and renewals.
type SubscriptionSummaryResponse struct {
	Active            uint16                                       `json:"active" example:"10" description:"Number of active subscriptions"`
	ActivePersonal    uint16                                       `json:"active_personal" example:"5" description:"Number of active personal subscriptions"`
	ActiveFamily      uint16                                       `json:"active_family" example:"5" description:"Number of active family subscriptions"`
	TotalMonthly      AmountModel                                  `json:"total_monthly" description:"Total monthly subscription costs"`
	TotalYearly       AmountModel                                  `json:"total_yearly"  description:"Total yearly subscription costs"`
	TotalLastMonth    AmountModel                                  `json:"total_last_month" description:"Total monthly subscription costs for last month"`
	TotalLastYear     AmountModel                                  `json:"total_last_year" description:"Total yearly subscription costs for last year"`
	PersonalMonthly   AmountModel                                  `json:"personal_monthly" description:"Personal monthly subscription costs"`
	PersonalYearly    AmountModel                                  `json:"personal_yearly" description:"Personal yearly subscription costs"`
	PersonalLastMonth AmountModel                                  `json:"personal_last_month" description:"Personal monthly subscription costs for last month"`
	PersonalLastYear  AmountModel                                  `json:"personal_last_year" description:"Personal yearly subscription costs for last year"`
	FamilyMonthly     AmountModel                                  `json:"family_monthly" description:"Family monthly subscription costs"`
	FamilyYearly      AmountModel                                  `json:"family_yearly" description:"Family yearly subscription costs"`
	FamilyLastMonth   AmountModel                                  `json:"family_last_month" description:"Family monthly subscription costs for last month"`
	FamilyLastYear    AmountModel                                  `json:"family_last_year" description:"Family yearly subscription costs for last year"`
	TopProviders      []SubscriptionSummaryTopProviderResponse     `json:"top_providers" description:"List of top providers by cost"`
	UpcomingRenewals  []SubscriptionSummaryUpcomingRenewalResponse `json:"upcoming_renewals" description:"List of upcoming subscription renewals"`
	TopLabels         []SubscriptionSummaryTopLabelResponse        `json:"top_labels" description:"List of top labels by cost"`
}
