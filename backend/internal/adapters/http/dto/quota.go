package dto

import (
	"github.com/mistribe/subtracker/internal/domain/billing"
)

type QuotaUsageModel struct {
	Feature   string `json:"feature" enums:"unknown,subscriptions,active_subscriptions_count,custom_labels,custom_labels_count,custom_providers,custom_providers_count,family,family_members_count"`
	Type      string `json:"type" enums:"boolean,quota,unknown"`
	Enabled   bool   `json:"enabled" example:"true"`
	Limit     *int64 `json:"limit,omitempty"`
	Used      *int64 `json:"used,omitempty"`
	Remaining *int64 `json:"remaining,omitempty"`
}

func NewQuotaUsageModel(eff billing.EffectiveEntitlement) QuotaUsageModel {
	return QuotaUsageModel{
		Feature:   billing.FeatureIDToString(eff.FeatureID),
		Type:      billing.FeatureTypeToString(eff.Type),
		Enabled:   eff.Enabled,
		Limit:     eff.Limit,
		Used:      eff.Used,
		Remaining: eff.Remaining,
	}
}
