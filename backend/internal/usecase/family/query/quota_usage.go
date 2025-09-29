package query

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/billing"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type GetQuotaUsage struct {
}

type GetQuotaUsageHandler struct {
	entitlement    ports.EntitlementResolver
	authentication ports.Authentication
}

func NewGetQuotaUsageHandler(
	entitlement ports.EntitlementResolver,
	authentication ports.Authentication) *GetQuotaUsageHandler {
	return &GetQuotaUsageHandler{
		entitlement:    entitlement,
		authentication: authentication,
	}
}

func (h GetQuotaUsageHandler) Handle(
	ctx context.Context,
	_ GetQuotaUsage) result.Result[billing.EffectiveEntitlement] {
	connectedAccount := h.authentication.MustGetConnectedAccount(ctx)
	eff, err := h.entitlement.Resolve(ctx, connectedAccount, billing.FeatureIdFamilyMembersCount)
	if err != nil {
		return result.Fail[billing.EffectiveEntitlement](err)
	}
	return result.Success(eff)
}
