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
	FeatureIdImportSubscriptions
	FeatureIdExportSubscriptions
	FeatureIdCustomLabels
	FeatureIdCustomLabelsCount
	FeatureIdImportCustomLabels
	FeatureIdExportCustomLabels
	FeatureIdCustomProviders
	FeatureIdCustomProvidersCount // quota, gated by custom_providers
	FeatureIdImportCustomProviders
	FeatureIdExportCustomProviders
	FeatureIdFamily
	FeatureIdFamilyMembersCount
)

const (
	FeatureIdUnknownString                  = "unknown"
	FeatureIdSubscriptionsString            = "subscriptions"
	FeatureIdActiveSubscriptionsCountString = "active_subscriptions_count"
	FeatureIdImportSubscriptionsString      = "import_subscriptions"
	FeatureIdExportSubscriptionsString      = "export_subscriptions"
	FeatureIdCustomLabelsString             = "custom_labels"
	FeatureIdCustomLabelsCountString        = "custom_labels_count"
	FeatureIdImportCustomLabelsString       = "import_custom_labels"
	FeatureIdExportCustomLabelsString       = "export_custom_labels"
	FeatureIdCustomProvidersString          = "custom_providers"
	FeatureIdCustomProvidersCountString     = "custom_providers_count"
	FeatureIdImportCustomProvidersString    = "import_custom_providers"
	FeatureIdExportCustomProvidersString    = "export_custom_providers"
	FeatureIdFamilyString                   = "family"
	FeatureIdFamilyMembersCountString       = "family_members_count"
)

func FeatureIDToString(feature types.FeatureID) string {
	switch feature {
	case FeatureIdSubscriptions:
		return FeatureIdSubscriptionsString
	case FeatureIdActiveSubscriptionsCount:
		return FeatureIdActiveSubscriptionsCountString
	case FeatureIdImportSubscriptions:
		return FeatureIdImportSubscriptionsString
	case FeatureIdExportSubscriptions:
		return FeatureIdExportSubscriptionsString
	case FeatureIdCustomLabels:
		return FeatureIdCustomLabelsString
	case FeatureIdCustomLabelsCount:
		return FeatureIdCustomLabelsCountString
	case FeatureIdImportCustomLabels:
		return FeatureIdImportCustomLabelsString
	case FeatureIdExportCustomLabels:
		return FeatureIdExportCustomLabelsString
	case FeatureIdCustomProviders:
		return FeatureIdCustomProvidersString
	case FeatureIdCustomProvidersCount:
		return FeatureIdCustomProvidersCountString
	case FeatureIdImportCustomProviders:
		return FeatureIdImportCustomProvidersString
	case FeatureIdExportCustomProviders:
		return FeatureIdExportCustomProvidersString
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
			GatedBy:     x.P(FeatureIdSubscriptions),
		},
		FeatureIdImportSubscriptions: {
			ID:          FeatureIdImportSubscriptions,
			Type:        FeatureBoolean,
			Description: "Import subscriptions",
			GatedBy:     x.P(FeatureIdSubscriptions),
		},
		FeatureIdExportSubscriptions: {
			ID:          FeatureIdExportSubscriptions,
			Type:        FeatureBoolean,
			Description: "Export subscriptions",
			GatedBy:     x.P(FeatureIdSubscriptions),
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
			GatedBy:     x.P(FeatureIdCustomLabels),
		},
		FeatureIdImportCustomLabels: {
			ID:          FeatureIdImportCustomLabels,
			Type:        FeatureBoolean,
			Description: "Import custom labels",
			GatedBy:     x.P(FeatureIdCustomLabels),
		},
		FeatureIdExportCustomLabels: {
			ID:          FeatureIdExportCustomLabels,
			Type:        FeatureBoolean,
			Description: "Export custom labels",
			GatedBy:     x.P(FeatureIdCustomLabels),
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
		FeatureIdImportCustomProviders: {
			ID:          FeatureIdImportCustomProviders,
			Type:        FeatureBoolean,
			Description: "Import custom providers",
			GatedBy:     x.P(FeatureIdCustomProviders),
		},
		FeatureIdExportCustomProviders: {
			ID:          FeatureIdExportCustomProviders,
			Type:        FeatureBoolean,
			Description: "Export custom providers",
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
			FeatureIdImportSubscriptions: newBoolEntitlement(types.PlanFree, FeatureIdImportSubscriptions, true),
			FeatureIdExportSubscriptions: newBoolEntitlement(types.PlanFree, FeatureIdExportSubscriptions, true),

			FeatureIdCustomLabels:       newBoolEntitlement(types.PlanFree, FeatureIdCustomLabels, true),
			FeatureIdCustomLabelsCount:  newQuotaEntitlement(types.PlanFree, FeatureIdCustomLabelsCount, 5),
			FeatureIdImportCustomLabels: newBoolEntitlement(types.PlanFree, FeatureIdImportCustomLabels, false),
			FeatureIdExportCustomLabels: newBoolEntitlement(types.PlanFree, FeatureIdExportCustomLabels, false),

			FeatureIdCustomProviders:       newBoolEntitlement(types.PlanFree, FeatureIdCustomProviders, true),
			FeatureIdCustomProvidersCount:  newQuotaEntitlement(types.PlanFree, FeatureIdCustomProvidersCount, 5),
			FeatureIdImportCustomProviders: newBoolEntitlement(types.PlanFree, FeatureIdImportCustomProviders, false),
			FeatureIdExportCustomProviders: newBoolEntitlement(types.PlanFree, FeatureIdExportCustomProviders, false),

			FeatureIdFamily:             newBoolEntitlement(types.PlanFree, FeatureIdFamily, true),
			FeatureIdFamilyMembersCount: newQuotaEntitlement(types.PlanFree, FeatureIdFamilyMembersCount, 3),
		},
		types.PlanPremium: {
			FeatureIdSubscriptions: newBoolEntitlement(types.PlanPremium, FeatureIdSubscriptions, true),
			FeatureIdActiveSubscriptionsCount: newQuotaEntitlement(types.PlanPremium, FeatureIdActiveSubscriptionsCount,
				100),
			FeatureIdImportSubscriptions: newBoolEntitlement(types.PlanFree, FeatureIdImportSubscriptions, true),
			FeatureIdExportSubscriptions: newBoolEntitlement(types.PlanFree, FeatureIdExportSubscriptions, true),

			FeatureIdCustomLabels:       newBoolEntitlement(types.PlanPremium, FeatureIdCustomLabels, true),
			FeatureIdCustomLabelsCount:  newQuotaEntitlement(types.PlanPremium, FeatureIdCustomLabelsCount, 100),
			FeatureIdImportCustomLabels: newBoolEntitlement(types.PlanFree, FeatureIdImportCustomLabels, true),
			FeatureIdExportCustomLabels: newBoolEntitlement(types.PlanFree, FeatureIdExportCustomLabels, true),

			FeatureIdCustomProviders:       newBoolEntitlement(types.PlanPremium, FeatureIdCustomProviders, true),
			FeatureIdCustomProvidersCount:  newQuotaEntitlement(types.PlanPremium, FeatureIdCustomProvidersCount, 100),
			FeatureIdImportCustomProviders: newBoolEntitlement(types.PlanFree, FeatureIdImportCustomProviders, true),
			FeatureIdExportCustomProviders: newBoolEntitlement(types.PlanFree, FeatureIdExportCustomProviders, true),

			FeatureIdFamily:             newBoolEntitlement(types.PlanPremium, FeatureIdFamily, true),
			FeatureIdFamilyMembersCount: newQuotaEntitlement(types.PlanPremium, FeatureIdFamilyMembersCount, 25),
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
