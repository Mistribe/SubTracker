package billing

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/billing"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/x"
)

type entitlementResolver struct {
	usage billing.UsageRepository
}

// NewEntitlementResolver constructs a Resolver with the given dependencies.
// periodCalc may be nil to use the default monthly calculator.
func NewEntitlementResolver(usage billing.UsageRepository) ports.EntitlementResolver {
	return &entitlementResolver{
		usage: usage,
	}
}

func (r *entitlementResolver) Resolve(ctx context.Context,
	account billing.Account,
	featureID billing.FeatureID) (billing.EffectiveEntitlement, error) {
	if featureID == billing.FeatureIdUnknown {
		return billing.EffectiveEntitlement{}, billing.ErrFeatureNotFound
	}

	planID := account.PlanID()
	if planID == billing.PlanUnknown {
		return billing.EffectiveEntitlement{}, billing.ErrPlanNotFound
	}

	entitlement, hasEntitlement := billing.Entitlements[planID][featureID]

	feature, hasFeature := billing.Features[entitlement.FeatureID]
	if !hasFeature {
		return billing.EffectiveEntitlement{}, billing.ErrFeatureNotFound
	}
	gateAllowed := true
	var err error
	if feature.GatedBy != nil {
		gateAllowed, err = r.gateAllows(account, *feature.GatedBy)
		if err != nil {
			return billing.EffectiveEntitlement{}, err
		}
	}

	switch feature.Type {
	case billing.FeatureBoolean:
		enabled := false
		if hasEntitlement && entitlement.Allowed != nil {
			enabled = *entitlement.Allowed
		}
		enabled = enabled && gateAllowed
		return billing.EffectiveEntitlement{
			FeatureID: feature.ID,
			Type:      feature.Type,
			Enabled:   enabled,
			Limit:     nil,
			Used:      nil,
			Remaining: nil,
		}, nil

	case billing.FeatureQuota:
		// If gate is disabled, short-circuit: the quota is not usable at all.
		if !gateAllowed {
			return billing.EffectiveEntitlement{
				FeatureID: feature.ID,
				Type:      feature.Type,
				Enabled:   false,
				Limit:     nil,
				Used:      nil,
				Remaining: nil,
			}, nil
		}

		var limit *int64
		// start from plan entitlement
		if hasEntitlement {
			limit = entitlement.Limit // may be nil (unlimited)
		} else {
			// No entitlement given -> treat as zero limit
			zero := int64(0)
			limit = &zero
		}

		usage, foundUsage, err := r.usage.Get(ctx, account, feature.ID)
		if err != nil {
			return billing.EffectiveEntitlement{}, err
		}
		used := int64(0)
		if foundUsage {
			used = usage.Used
		}
		var remaining *int64
		var enabled bool
		if limit == nil {
			remaining = nil // unlimited
			enabled = true  // unlimited is always enabled for quota features
		} else {
			rv := *limit - used
			if rv < 0 {
				rv = 0
			}
			remaining = x.P(rv)
			enabled = rv > 0
		}

		return billing.EffectiveEntitlement{
			FeatureID: feature.ID,
			Type:      feature.Type,
			Enabled:   enabled,
			Limit:     limit,
			Used:      x.P(used),
			Remaining: remaining,
		}, nil

	default:
		return billing.EffectiveEntitlement{}, billing.ErrInvalidFeatureType
	}
}

func (r *entitlementResolver) gateAllows(account billing.Account, gate billing.FeatureID) (bool,
	error) {
	planID := account.PlanID()
	if planID == billing.PlanUnknown {
		return false, billing.ErrPlanNotFound
	}
	entitlement, hasEntitlement := billing.Entitlements[planID][gate]
	if !hasEntitlement {
		return false, billing.ErrPlanNotFound
	}
	enabled := false
	if entitlement.Allowed != nil {
		enabled = *entitlement.Allowed
	}
	return enabled, nil
}

func (r *entitlementResolver) CheckBoolean(ctx context.Context, account billing.Account,
	featureID billing.FeatureID) (bool, error) {
	eff, err := r.Resolve(ctx, account, featureID)
	if err != nil {
		return false, err
	}
	if eff.Type != billing.FeatureBoolean {
		return false, billing.ErrInvalidFeatureType
	}
	return eff.Enabled, nil
}

func (r *entitlementResolver) CheckQuota(ctx context.Context, account billing.Account, featureID billing.FeatureID,
	needed int64) (bool,
	billing.EffectiveEntitlement, error) {
	if needed <= 0 {
		needed = 1
	}
	eff, err := r.Resolve(ctx, account, featureID)
	if err != nil {
		return false, billing.EffectiveEntitlement{}, err
	}
	if eff.Type != billing.FeatureQuota {
		return false, billing.EffectiveEntitlement{}, billing.ErrInvalidFeatureType
	}
	if !eff.Enabled {
		return false, eff, nil
	}
	if eff.Limit == nil {
		// unlimited
		return true, eff, nil
	}
	remaining := int64(0)
	if eff.Remaining != nil {
		remaining = *eff.Remaining
	}
	allowed := remaining >= needed
	return allowed, eff, nil
}
