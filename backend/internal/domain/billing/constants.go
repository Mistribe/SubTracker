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
	FeatureIdActiveSubscriptions
	FeatureIdCustomLabels
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
		FeatureIdActiveSubscriptions: {
			ID:          FeatureIdActiveSubscriptions,
			Type:        FeatureQuota,
			Description: "Number of active subscriptions",
		},
		FeatureIdCustomLabels: {
			ID:          FeatureIdCustomLabels,
			Type:        FeatureBoolean,
			Description: "Custom labels",
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
			FeatureIdSubscriptions:        newBoolEntitlement(PlanFree, FeatureIdSubscriptions, true),
			FeatureIdActiveSubscriptions:  newQuotaEntitlement(PlanFree, FeatureIdActiveSubscriptions, 10),
			FeatureIdCustomLabels:         newBoolEntitlement(PlanFree, FeatureIdCustomLabels, false),
			FeatureIdCustomProviders:      newBoolEntitlement(PlanFree, FeatureIdCustomProviders, true),
			FeatureIdCustomProvidersCount: newQuotaEntitlement(PlanFree, FeatureIdCustomProvidersCount, 5),
			FeatureIdFamily:               newBoolEntitlement(PlanFree, FeatureIdFamily, true),
			FeatureIdFamilyMembersCount:   newQuotaEntitlement(PlanFree, FeatureIdFamilyMembersCount, 3),
		},
		PlanPremium: {
			FeatureIdSubscriptions:        newBoolEntitlement(PlanPremium, FeatureIdSubscriptions, true),
			FeatureIdActiveSubscriptions:  newQuotaEntitlement(PlanPremium, FeatureIdActiveSubscriptions, 100),
			FeatureIdCustomLabels:         newBoolEntitlement(PlanPremium, FeatureIdCustomLabels, true),
			FeatureIdCustomProviders:      newBoolEntitlement(PlanPremium, FeatureIdCustomProviders, true),
			FeatureIdCustomProvidersCount: newQuotaEntitlement(PlanPremium, FeatureIdCustomProvidersCount, 100),
			FeatureIdFamily:               newBoolEntitlement(PlanPremium, FeatureIdFamily, true),
			FeatureIdFamilyMembersCount:   newQuotaEntitlement(PlanPremium, FeatureIdFamilyMembersCount, 25),
		},
	}
}
