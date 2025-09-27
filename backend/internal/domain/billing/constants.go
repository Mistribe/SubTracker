package billing

import (
	"errors"

	"github.com/mistribe/subtracker/pkg/x"
)

const (
	PlanUnknown PlanID = iota
	PlanFree
	PlanPremium
)

var (
	ErrInvalidPlan = errors.New("invalid plan")
)

const (
	PlanUnknownString = "unknown"
	PlanFreeString    = "free"
	PlanPremiumString = "premium"
)

func (p PlanID) String() string {
	switch p {
	case PlanFree:
		return PlanFreeString
	case PlanPremium:
		return PlanPremiumString
	default:
		return PlanUnknownString
	}
}

func ParsePlan(input string) (PlanID, error) {
	switch input {
	case PlanUnknownString:
		return PlanUnknown, nil
	case PlanFreeString:
		return PlanFree, nil
	case PlanPremiumString:
		return PlanPremium, nil
	}

	return PlanUnknown, ErrInvalidPlan
}

func ParsePlanOrDefault(input string, defaultValue PlanID) PlanID {
	p, err := ParsePlan(input)
	if err != nil {
		return defaultValue
	}
	return p
}

const (
	FeatureUnknown FeatureType = iota
	FeatureBoolean
	FeatureQuota
)

const (
	FeatureIdUnknown FeatureID = iota
	FeatureIdSubscriptions
	FeatureIdActiveSubscriptionsCount
	FeatureIdCustomLabels
	FeatureIdCustomLabelsCount
	FeatureIdCustomProviders
	FeatureIdCustomProvidersCount // quota, gated by custom_providers
	FeatureIdFamily
	FeatureIdFamilyMembersCount
)

var (
	Features = map[FeatureID]Feature{
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
	Entitlements map[PlanID]map[FeatureID]PlanEntitlement
)

func init() {
	Entitlements = map[PlanID]map[FeatureID]PlanEntitlement{
		PlanFree: {
			FeatureIdSubscriptions:            newBoolEntitlement(PlanFree, FeatureIdSubscriptions, true),
			FeatureIdActiveSubscriptionsCount: newQuotaEntitlement(PlanFree, FeatureIdActiveSubscriptionsCount, 10),
			FeatureIdCustomLabels:             newBoolEntitlement(PlanFree, FeatureIdCustomLabels, true),
			FeatureIdCustomLabelsCount:        newQuotaEntitlement(PlanFree, FeatureIdCustomLabelsCount, 5),
			FeatureIdCustomProviders:          newBoolEntitlement(PlanFree, FeatureIdCustomProviders, true),
			FeatureIdCustomProvidersCount:     newQuotaEntitlement(PlanFree, FeatureIdCustomProvidersCount, 5),
			FeatureIdFamily:                   newBoolEntitlement(PlanFree, FeatureIdFamily, true),
			FeatureIdFamilyMembersCount:       newQuotaEntitlement(PlanFree, FeatureIdFamilyMembersCount, 3),
		},
		PlanPremium: {
			FeatureIdSubscriptions:            newBoolEntitlement(PlanPremium, FeatureIdSubscriptions, true),
			FeatureIdActiveSubscriptionsCount: newQuotaEntitlement(PlanPremium, FeatureIdActiveSubscriptionsCount, 100),
			FeatureIdCustomLabels:             newBoolEntitlement(PlanPremium, FeatureIdCustomLabels, true),
			FeatureIdCustomLabelsCount:        newQuotaEntitlement(PlanFree, FeatureIdCustomLabelsCount, 100),

			FeatureIdCustomProviders:      newBoolEntitlement(PlanPremium, FeatureIdCustomProviders, true),
			FeatureIdCustomProvidersCount: newQuotaEntitlement(PlanPremium, FeatureIdCustomProvidersCount, 100),
			FeatureIdFamily:               newBoolEntitlement(PlanPremium, FeatureIdFamily, true),
			FeatureIdFamilyMembersCount:   newQuotaEntitlement(PlanPremium, FeatureIdFamilyMembersCount, 25),
		},
	}
}
