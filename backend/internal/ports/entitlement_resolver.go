package ports

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/billing"
)

// EntitlementResolver evaluates and consumes entitlements for an account.
type EntitlementResolver interface {
	// Resolve returns the effective entitlement for a feature at time 'at'.
	Resolve(ctx context.Context,
		account billing.Account,
		featureID billing.FeatureID) (billing.EffectiveEntitlement, error)

	// CheckBoolean returns whether a boolean feature is enabled.
	CheckBoolean(ctx context.Context,
		account billing.Account,
		featureID billing.FeatureID) (bool, error)

	// CheckQuota reports if 'needed' units can be consumed without exceeding the limit.
	// It also returns the current effective entitlement snapshot.
	CheckQuota(ctx context.Context,
		featureID billing.FeatureID,
		needed int64) (allowed bool, eff billing.EffectiveEntitlement, err error)

	CheckQuotaForAccount(ctx context.Context,
		account billing.Account,
		featureID billing.FeatureID,
		needed int64) (allowed bool, eff billing.EffectiveEntitlement, err error)
}
