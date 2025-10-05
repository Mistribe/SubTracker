package billing

import (
	"time"

	"github.com/mistribe/subtracker/internal/domain/types"
)

// Strongly-typed identifiers for clarity.

// Feature describes a capability your product exposes.
type Feature struct {
	ID          types.FeatureID
	Type        types.FeatureType
	Description string

	// GatedBy optionally references another feature that must be enabled
	// (typically a boolean feature) for this feature to be usable.
	// When set, the effective Enabled result of this feature is AND-ed with the
	// effective Enabled of the gating feature.
	GatedBy *types.FeatureID
}

func (f Feature) IsQuota() bool {
	return f.Type == FeatureQuota
}

func (f Feature) IsBoolean() bool {
	return f.Type == FeatureBoolean
}

// PlanEntitlement defines what a plan grants for a feature.
// For boolean features, Allowed is used.
// For quota features, Limit is used; nil means "unlimited".
type PlanEntitlement struct {
	PlanID    types.PlanID
	FeatureID types.FeatureID
	Allowed   *bool  // only for FeatureBoolean
	Limit     *int64 // only for FeatureQuota; nil => unlimited
}

func newBoolEntitlement(planID types.PlanID, featureID types.FeatureID, allowed bool) PlanEntitlement {
	return PlanEntitlement{
		PlanID:    planID,
		FeatureID: featureID,
		Allowed:   &allowed,
	}
}

func newQuotaEntitlement(planID types.PlanID, featureID types.FeatureID, limit int64) PlanEntitlement {
	return PlanEntitlement{
		PlanID:    planID,
		FeatureID: featureID,
		Limit:     &limit,
	}
}

// UsageCounter tracks usage for a feature within a specific period.
type UsageCounter struct {
	FeatureID types.FeatureID
	Used      int64
	UpdatedAt time.Time
}

// EffectiveEntitlement represents the evaluated entitlements for an account.
type EffectiveEntitlement struct {
	FeatureID types.FeatureID
	Type      types.FeatureType

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
