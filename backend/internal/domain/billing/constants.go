package billing

import (
	"errors"

	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/pkg/x"
)

var (
	ErrInvalidPlan = errors.New("invalid plan")
)

const (
	FeatureUnknown types.FeatureType = iota
	FeatureBoolean
	FeatureQuota
)

const (
	FeatureUnknownString = "unknown"
	FeatureBooleanString = "boolean"
	FeatureQuotaString   = "quota"
)

func FeatureTypeToString(feature types.FeatureType) string {
	switch feature {
	case FeatureBoolean:
		return FeatureBooleanString
	case FeatureQuota:
		return FeatureQuotaString
	default:
		return FeatureUnknownString
	}
}

const (
	FeatureIdUnknown types.FeatureID = iota
	FeatureIdSubscriptions
	FeatureIdActiveSubscriptionsCount
	FeatureIdCustomLabels
	FeatureIdCustomLabelsCount
	FeatureIdCustomProviders
	FeatureIdCustomProvidersCount // quota, gated by custom_providers
	FeatureIdFamily
	FeatureIdFamilyMembersCount
)

const (
	FeatureIdUnknownString                  = "unknown"
	FeatureIdSubscriptionsString            = "subscriptions"
	FeatureIdActiveSubscriptionsCountString = "active_subscriptions_count"
	FeatureIdCustomLabelsString             = "custom_labels"
	FeatureIdCustomLabelsCountString        = "custom_labels_count"
	FeatureIdCustomProvidersString          = "custom_providers"
	FeatureIdCustomProvidersCountString     = "custom_providers_count"
	FeatureIdFamilyString                   = "family"
	FeatureIdFamilyMembersCountString       = "family_members_count"
)

func FeatureIDToString(feature types.FeatureID) string {
	switch feature {
	case FeatureIdSubscriptions:
		return FeatureIdSubscriptionsString
	case FeatureIdActiveSubscriptionsCount:
		return FeatureIdActiveSubscriptionsCountString
	case FeatureIdCustomLabels:
		return FeatureIdCustomLabelsString
	case FeatureIdCustomLabelsCount:
		return FeatureIdCustomLabelsCountString
	case FeatureIdCustomProviders:
		return FeatureIdCustomProvidersString
	case FeatureIdCustomProvidersCount:
		return FeatureIdCustomProvidersCountString
	case FeatureIdFamily:
		return FeatureIdFamilyString
	case FeatureIdFamilyMembersCount:
		return FeatureIdFamilyMembersCountString
	default:
		return FeatureIdUnknownString
	}
}

var (
	Features = map[types.FeatureID]Feature{
		FeatureIdSubscriptions: {
			ID:          FeatureIdSubscriptions,
			Type:        FeatureBoolean,
			Description: "Subscriptions",
		},
		FeatureIdActiveSubscriptionsCount: {
			ID:          FeatureIdActiveSubscriptionsCount,
			Type:        FeatureQuota,
			Description: "Number of active subscriptions",
		},
		FeatureIdCustomLabels: {
			ID:          FeatureIdCustomLabels,
			Type:        FeatureBoolean,
			Description: "Custom labels",
		},
		FeatureIdCustomLabelsCount: {
			ID:          FeatureIdCustomLabelsCount,
			Type:        FeatureQuota,
			Description: "Number of custom labels",
		},
		FeatureIdCustomProviders: {
			ID:          FeatureIdCustomProviders,
			Type:        FeatureBoolean,
			Description: "Custom providers",
		},
		FeatureIdCustomProvidersCount: {
			ID:          FeatureIdCustomProvidersCount,
			Type:        FeatureQuota,
			Description: "Maximum number of custom providers",
			GatedBy:     x.P(FeatureIdCustomProviders),
		},
		FeatureIdFamily: {
			ID:          FeatureIdFamily,
			Type:        FeatureBoolean,
			Description: "Family",
		},
		FeatureIdFamilyMembersCount: {
			ID:          FeatureIdFamilyMembersCount,
			Type:        FeatureQuota,
			Description: "Family members",
		},
	}
	Entitlements    map[types.PlanID]map[types.FeatureID]PlanEntitlement
	AllFeatures     []types.FeatureID
	QuotaFeatures   []types.FeatureID
	BooleanFeatures []types.FeatureID
)

func init() {
	Entitlements = map[types.PlanID]map[types.FeatureID]PlanEntitlement{
		types.PlanFree: {
			FeatureIdSubscriptions: newBoolEntitlement(types.PlanFree, FeatureIdSubscriptions, true),
			FeatureIdActiveSubscriptionsCount: newQuotaEntitlement(types.PlanFree, FeatureIdActiveSubscriptionsCount,
				10),
			FeatureIdCustomLabels:         newBoolEntitlement(types.PlanFree, FeatureIdCustomLabels, true),
			FeatureIdCustomLabelsCount:    newQuotaEntitlement(types.PlanFree, FeatureIdCustomLabelsCount, 5),
			FeatureIdCustomProviders:      newBoolEntitlement(types.PlanFree, FeatureIdCustomProviders, true),
			FeatureIdCustomProvidersCount: newQuotaEntitlement(types.PlanFree, FeatureIdCustomProvidersCount, 5),
			FeatureIdFamily:               newBoolEntitlement(types.PlanFree, FeatureIdFamily, true),
			FeatureIdFamilyMembersCount:   newQuotaEntitlement(types.PlanFree, FeatureIdFamilyMembersCount, 3),
		},
		types.PlanPremium: {
			FeatureIdSubscriptions: newBoolEntitlement(types.PlanPremium, FeatureIdSubscriptions, true),
			FeatureIdActiveSubscriptionsCount: newQuotaEntitlement(types.PlanPremium, FeatureIdActiveSubscriptionsCount,
				100),
			FeatureIdCustomLabels:      newBoolEntitlement(types.PlanPremium, FeatureIdCustomLabels, true),
			FeatureIdCustomLabelsCount: newQuotaEntitlement(types.PlanPremium, FeatureIdCustomLabelsCount, 100),

			FeatureIdCustomProviders:      newBoolEntitlement(types.PlanPremium, FeatureIdCustomProviders, true),
			FeatureIdCustomProvidersCount: newQuotaEntitlement(types.PlanPremium, FeatureIdCustomProvidersCount, 100),
			FeatureIdFamily:               newBoolEntitlement(types.PlanPremium, FeatureIdFamily, true),
			FeatureIdFamilyMembersCount:   newQuotaEntitlement(types.PlanPremium, FeatureIdFamilyMembersCount, 25),
		},
	}

	for _, feature := range Features {
		if feature.IsQuota() {
			QuotaFeatures = append(QuotaFeatures, feature.ID)
		}
		if feature.IsBoolean() {
			BooleanFeatures = append(BooleanFeatures, feature.ID)
		}
		AllFeatures = append(AllFeatures, feature.ID)
	}
}
