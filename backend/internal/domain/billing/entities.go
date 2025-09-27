package billing

import (
	"time"
)

// Strongly-typed identifiers for clarity.
type (
	FeatureType uint8
	PlanID      uint8
	FeatureID   uint8
)

// Feature describes a capability your product exposes.
type Feature struct {
	ID          FeatureID
	Type        FeatureType
	Description string

	// GatedBy optionally references another feature that must be enabled
	// (typically a boolean feature) for this feature to be usable.
	// When set, the effective Enabled result of this feature is AND-ed with the
	// effective Enabled of the gating feature.
	GatedBy *FeatureID
}

// PlanEntitlement defines what a plan grants for a feature.
// For boolean features, Allowed is used.
// For quota features, Limit is used; nil means "unlimited".
type PlanEntitlement struct {
	PlanID    PlanID
	FeatureID FeatureID
	Allowed   *bool  // only for FeatureBoolean
	Limit     *int64 // only for FeatureQuota; nil => unlimited
}

func newBoolEntitlement(planID PlanID, featureID FeatureID, allowed bool) PlanEntitlement {
	return PlanEntitlement{
		PlanID:    planID,
		FeatureID: featureID,
		Allowed:   &allowed,
	}
}

func newQuotaEntitlement(planID PlanID, featureID FeatureID, limit int64) PlanEntitlement {
	return PlanEntitlement{
		PlanID:    planID,
		FeatureID: featureID,
		Limit:     &limit,
	}
}

// UsageCounter tracks usage for a feature within a specific period.
type UsageCounter struct {
	FeatureID FeatureID
	Used      int64
	UpdatedAt time.Time
}

// EffectiveEntitlement represents the evaluated entitlements for an account.
type EffectiveEntitlement struct {
	FeatureID FeatureID
	Type      FeatureType

	// Enabled is true when:
	// - boolean feature -> allowed is true
	// - quota feature -> limit is unlimited or remaining > 0
	// - AND any gating (GatedBy) feature is enabled
	Enabled bool

	// For quota features:
	// Limit is nil for unlimited, non-nil otherwise.
	// Used and Remaining are only set for quota features; nil for boolean.
	Limit     *int64
	Used      *int64
	Remaining *int64
}
