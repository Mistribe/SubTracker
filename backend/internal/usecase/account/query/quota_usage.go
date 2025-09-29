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
	_ GetQuotaUsage) result.Result[[]billing.EffectiveEntitlement] {
	connectedAccount := h.authentication.MustGetConnectedAccount(ctx)
	effs, err := h.entitlement.Resolves(ctx, connectedAccount, billing.QuotaFeatures)
	if err != nil {
		return result.Fail[[]billing.EffectiveEntitlement](err)
	}
	return result.Success(effs)
}
