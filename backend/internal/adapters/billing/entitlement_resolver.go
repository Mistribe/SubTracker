package billing

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/account"
	"github.com/mistribe/subtracker/internal/domain/billing"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/x"
)

type entitlementResolver struct {
	usage          ports.UsageRepository
	authentication ports.Authentication
}

func NewEntitlementResolver(
	usage ports.UsageRepository,
	authentication ports.Authentication) ports.EntitlementResolver {
	return &entitlementResolver{
		usage:          usage,
		authentication: authentication,
	}
}

func (r *entitlementResolver) Resolve(
	ctx context.Context,
	account account.ConnectedAccount,
	featureID types.FeatureID) (billing.EffectiveEntitlement, error) {
	if featureID == billing.FeatureIdUnknown {
		return billing.EffectiveEntitlement{}, billing.ErrFeatureNotFound
	}

	planID := account.PlanID()
	if planID == types.PlanUnknown {
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
				Enabled:   true,
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

		usage, foundUsage, err := r.usage.Get(ctx, account.UserID(), feature)
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

func (r *entitlementResolver) Resolves(
	ctx context.Context,
	account account.ConnectedAccount,
	featureIDs []types.FeatureID) ([]billing.EffectiveEntitlement, error) {
	// If any feature is unknown upfront, fail early (mirrors Resolve behavior).
	for _, fid := range featureIDs {
		if fid == billing.FeatureIdUnknown {
			return nil, billing.ErrFeatureNotFound
		}
	}

	planID := account.PlanID()
	if planID == types.PlanUnknown {
		return nil, billing.ErrPlanNotFound
	}

	// Determine whether we need usage (any quota feature requested).
	needUsage := false
	for _, fid := range featureIDs {
		if ent, ok := billing.Entitlements[planID][fid]; ok { // mimic Resolve path using entitlement->feature lookup
			if feature, ok2 := billing.Features[ent.FeatureID]; ok2 && feature.IsQuota() {
				needUsage = true
				break
			}
		}
	}

	usageMap := make(map[types.FeatureID]billing.UsageCounter)
	if needUsage {
		counters, err := r.usage.GetAll(ctx, account.UserID())
		if err != nil {
			return nil, err
		}
		for _, c := range counters {
			usageMap[c.FeatureID] = c
		}
	}

	results := make([]billing.EffectiveEntitlement, 0, len(featureIDs))
	for _, fid := range featureIDs {
		entitlement, hasEntitlement := billing.Entitlements[planID][fid]
		feature, hasFeature := billing.Features[entitlement.FeatureID]
		if !hasFeature {
			return nil, billing.ErrFeatureNotFound
		}

		gateAllowed := true
		if feature.GatedBy != nil {
			allowed, err := r.gateAllows(account, *feature.GatedBy)
			if err != nil {
				return nil, err
			}
			gateAllowed = allowed
		}

		switch feature.Type {
		case billing.FeatureBoolean:
			enabled := false
			if hasEntitlement && entitlement.Allowed != nil {
				enabled = *entitlement.Allowed
			}
			enabled = enabled && gateAllowed
			results = append(results, billing.EffectiveEntitlement{
				FeatureID: feature.ID,
				Type:      feature.Type,
				Enabled:   enabled,
				Limit:     nil,
				Used:      nil,
				Remaining: nil,
			})
		case billing.FeatureQuota:
			if !gateAllowed {
				results = append(results, billing.EffectiveEntitlement{
					FeatureID: feature.ID,
					Type:      feature.Type,
					Enabled:   true,
					Limit:     nil,
					Used:      nil,
					Remaining: nil,
				})
				continue
			}
			var limit *int64
			if hasEntitlement {
				limit = entitlement.Limit
			} else {
				zero := int64(0)
				limit = &zero
			}
			used := int64(0)
			if uc, ok := usageMap[fid]; ok {
				used = uc.Used
			}
			var remaining *int64
			var enabled bool
			if limit == nil { // unlimited
				remaining = nil
				enabled = true
			} else {
				rv := *limit - used
				if rv < 0 {
					rv = 0
				}
				remaining = x.P(rv)
				enabled = rv > 0
			}
			results = append(results, billing.EffectiveEntitlement{
				FeatureID: feature.ID,
				Type:      feature.Type,
				Enabled:   enabled,
				Limit:     limit,
				Used:      x.P(used),
				Remaining: remaining,
			})
		default:
			return nil, billing.ErrInvalidFeatureType
		}
	}
	return results, nil
}

func (r *entitlementResolver) gateAllows(account account.ConnectedAccount, gate types.FeatureID) (
	bool,
	error) {
	planID := account.PlanID()
	if planID == types.PlanUnknown {
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

func (r *entitlementResolver) CheckBoolean(
	ctx context.Context, account account.ConnectedAccount,
	featureID types.FeatureID) (bool, error) {
	eff, err := r.Resolve(ctx, account, featureID)
	if err != nil {
		return false, err
	}
	if eff.Type != billing.FeatureBoolean {
		return false, billing.ErrInvalidFeatureType
	}
	return eff.Enabled, nil
}

func (r *entitlementResolver) CheckQuotaForAccount(
	ctx context.Context, account account.ConnectedAccount,
	featureID types.FeatureID, needed int64) (bool, billing.EffectiveEntitlement, error) {
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

func (r *entitlementResolver) CheckQuota(
	ctx context.Context, featureID types.FeatureID,
	needed int64) (
	bool,
	billing.EffectiveEntitlement, error) {
	connectedAccount := r.authentication.MustGetConnectedAccount(ctx)
	return r.CheckQuotaForAccount(ctx, connectedAccount, featureID, needed)
}
