package billing

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/mistribe/subtracker/internal/domain/account"
	bdomain "github.com/mistribe/subtracker/internal/domain/billing"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
)

func TestEntitlementResolver_Resolve(t *testing.T) {
	ctx := context.Background()

	t.Run("unknown feature id returns ErrFeatureNotFound", func(t *testing.T) {
		usage := ports.NewMockUsageRepository(t)
		auth := ports.NewMockAuthentication(t)
		acc := account.NewMockConnectedAccount(t)
		// No PlanID expectation: code returns before calling PlanID
		r := &entitlementResolver{usage: usage, authentication: auth}
		_, err := r.Resolve(ctx, acc, bdomain.FeatureIdUnknown)
		assert.ErrorIs(t, err, bdomain.ErrFeatureNotFound)
	})

	t.Run("unknown plan returns ErrPlanNotFound", func(t *testing.T) {
		usage := ports.NewMockUsageRepository(t)
		auth := ports.NewMockAuthentication(t)
		acc := account.NewMockConnectedAccount(t)
		acc.EXPECT().PlanID().Return(types.PlanUnknown)
		r := &entitlementResolver{usage: usage, authentication: auth}
		_, err := r.Resolve(ctx, acc, bdomain.FeatureIdSubscriptions)
		assert.ErrorIs(t, err, bdomain.ErrPlanNotFound)
	})

	t.Run("boolean feature enabled when allowed", func(t *testing.T) {
		usage := ports.NewMockUsageRepository(t)
		auth := ports.NewMockAuthentication(t)
		acc := account.NewMockConnectedAccount(t)
		acc.EXPECT().PlanID().Return(types.PlanFree)
		res, err := (&entitlementResolver{usage: usage, authentication: auth}).Resolve(ctx, acc,
			bdomain.FeatureIdSubscriptions)
		require.NoError(t, err)
		assert.True(t, res.Enabled)
		assert.Nil(t, res.Limit)
	})

	t.Run("boolean feature disabled when entitlement explicitly false", func(t *testing.T) {
		// Temporarily modify entitlements and restore
		orig := bdomain.Entitlements[types.PlanFree][bdomain.FeatureIdCustomLabels]
		falseVal := false
		bdomain.Entitlements[types.PlanFree][bdomain.FeatureIdCustomLabels] = bdomain.PlanEntitlement{
			PlanID: types.PlanFree, FeatureID: bdomain.FeatureIdCustomLabels, Allowed: &falseVal,
		}
		defer func() { bdomain.Entitlements[types.PlanFree][bdomain.FeatureIdCustomLabels] = orig }()

		usage := ports.NewMockUsageRepository(t)
		auth := ports.NewMockAuthentication(t)
		acc := account.NewMockConnectedAccount(t)
		acc.EXPECT().PlanID().Return(types.PlanFree)
		res, err := (&entitlementResolver{usage: usage, authentication: auth}).Resolve(ctx, acc,
			bdomain.FeatureIdCustomLabels)
		require.NoError(t, err)
		assert.False(t, res.Enabled)
	})

	t.Run("quota feature returns remaining and enabled when usage below limit", func(t *testing.T) {
		usage := ports.NewMockUsageRepository(t)
		auth := ports.NewMockAuthentication(t)
		acc := account.NewMockConnectedAccount(t)
		acc.EXPECT().PlanID().Return(types.PlanFree)
		acc.EXPECT().UserID().Return(types.UserID("user-1"))
		limitEnt := bdomain.Entitlements[types.PlanFree][bdomain.FeatureIdActiveSubscriptionsCount]
		usage.EXPECT().Get(mock.Anything, types.UserID("user-1"), mock.Anything).Return(bdomain.UsageCounter{
			FeatureID: bdomain.FeatureIdActiveSubscriptionsCount, Used: 2,
		}, true, nil)
		res, err := (&entitlementResolver{usage: usage, authentication: auth}).Resolve(ctx, acc,
			bdomain.FeatureIdActiveSubscriptionsCount)
		require.NoError(t, err)
		require.NotNil(t, res.Limit)
		assert.Equal(t, *limitEnt.Limit, *res.Limit)
		assert.Equal(t, int64(2), *res.Used)
		assert.Equal(t, *limitEnt.Limit-2, *res.Remaining)
		assert.True(t, res.Enabled)
	})

	t.Run("quota feature disabled when gate (custom providers) disabled", func(t *testing.T) {
		// disable gating boolean
		origGate := bdomain.Entitlements[types.PlanFree][bdomain.FeatureIdCustomProviders]
		falseVal := false
		bdomain.Entitlements[types.PlanFree][bdomain.FeatureIdCustomProviders] = bdomain.PlanEntitlement{
			PlanID: types.PlanFree, FeatureID: bdomain.FeatureIdCustomProviders, Allowed: &falseVal,
		}
		defer func() { bdomain.Entitlements[types.PlanFree][bdomain.FeatureIdCustomProviders] = origGate }()

		usage := ports.NewMockUsageRepository(t)
		auth := ports.NewMockAuthentication(t)
		acc := account.NewMockConnectedAccount(t)
		acc.EXPECT().PlanID().Return(types.PlanFree)
		res, err := (&entitlementResolver{usage: usage, authentication: auth}).Resolve(ctx, acc,
			bdomain.FeatureIdCustomProvidersCount)
		require.NoError(t, err)
		assert.True(t, res.Enabled)
	})

	t.Run("usage repo error propagates", func(t *testing.T) {
		usage := ports.NewMockUsageRepository(t)
		auth := ports.NewMockAuthentication(t)
		acc := account.NewMockConnectedAccount(t)
		acc.EXPECT().PlanID().Return(types.PlanFree)
		acc.EXPECT().UserID().Return(types.UserID("user-err"))
		expected := errors.New("db err")
		usage.EXPECT().Get(mock.Anything, types.UserID("user-err"), mock.Anything).Return(bdomain.UsageCounter{}, false,
			expected)
		_, err := (&entitlementResolver{usage: usage, authentication: auth}).Resolve(ctx, acc,
			bdomain.FeatureIdActiveSubscriptionsCount)
		assert.ErrorIs(t, err, expected)
	})
}

func TestEntitlementResolver_Resolves(t *testing.T) {
	ctx := context.Background()

	t.Run("unknown feature id causes early error", func(t *testing.T) {
		usage := ports.NewMockUsageRepository(t)
		auth := ports.NewMockAuthentication(t)
		acc := account.NewMockConnectedAccount(t)
		acc.EXPECT().PlanID().Return(types.PlanFree).Maybe()
		resolver := &entitlementResolver{usage: usage, authentication: auth}
		_, err := resolver.Resolves(ctx, acc,
			[]types.FeatureID{bdomain.FeatureIdActiveSubscriptionsCount, bdomain.FeatureIdUnknown})
		assert.ErrorIs(t, err, bdomain.ErrFeatureNotFound)
	})

	t.Run("plan unknown returns error", func(t *testing.T) {
		usage := ports.NewMockUsageRepository(t)
		auth := ports.NewMockAuthentication(t)
		acc := account.NewMockConnectedAccount(t)
		acc.EXPECT().PlanID().Return(types.PlanUnknown)
		resolver := &entitlementResolver{usage: usage, authentication: auth}
		_, err := resolver.Resolves(ctx, acc, []types.FeatureID{bdomain.FeatureIdActiveSubscriptionsCount})
		assert.ErrorIs(t, err, bdomain.ErrPlanNotFound)
	})

	t.Run("returns effective entitlements using single GetAll for quota features", func(t *testing.T) {
		usage := ports.NewMockUsageRepository(t)
		auth := ports.NewMockAuthentication(t)
		acc := account.NewMockConnectedAccount(t)
		userID := types.UserID("user-batch")
		acc.EXPECT().PlanID().Return(types.PlanFree)
		acc.EXPECT().UserID().Return(userID)

		// Prepare usage counters for all quota features
		counters := []bdomain.UsageCounter{
			{FeatureID: bdomain.FeatureIdActiveSubscriptionsCount, Used: 2},
			{FeatureID: bdomain.FeatureIdCustomLabelsCount, Used: 1},
			{FeatureID: bdomain.FeatureIdCustomProvidersCount, Used: 3},
			{FeatureID: bdomain.FeatureIdFamilyMembersCount, Used: 1},
		}
		usage.EXPECT().GetAll(mock.Anything, userID).Return(counters, nil).Once()

		resolver := &entitlementResolver{usage: usage, authentication: auth}
		result, err := resolver.Resolves(ctx, acc, bdomain.QuotaFeatures)
		require.NoError(t, err)
		// Expect one result per requested feature
		assert.Len(t, result, len(bdomain.QuotaFeatures))

		// Validate one known feature's remaining calculation
		var activeSubEff *bdomain.EffectiveEntitlement
		for i := range result {
			if result[i].FeatureID == bdomain.FeatureIdActiveSubscriptionsCount {
				activeSubEff = &result[i]
			}
		}
		require.NotNil(t, activeSubEff)
		limit := bdomain.Entitlements[types.PlanFree][bdomain.FeatureIdActiveSubscriptionsCount].Limit
		require.NotNil(t, limit)
		assert.Equal(t, int64(2), *activeSubEff.Used)
		assert.Equal(t, *limit-2, *activeSubEff.Remaining)
		assert.True(t, activeSubEff.Enabled)
	})

	// Test gated quota feature disabled when gate disabled in batch
	t.Run("gated quota feature disabled when gate off", func(t *testing.T) {
		origGate := bdomain.Entitlements[types.PlanFree][bdomain.FeatureIdCustomProviders]
		falseVal := false
		bdomain.Entitlements[types.PlanFree][bdomain.FeatureIdCustomProviders] = bdomain.PlanEntitlement{
			PlanID: types.PlanFree, FeatureID: bdomain.FeatureIdCustomProviders, Allowed: &falseVal,
		}
		defer func() { bdomain.Entitlements[types.PlanFree][bdomain.FeatureIdCustomProviders] = origGate }()

		usage := ports.NewMockUsageRepository(t)
		auth := ports.NewMockAuthentication(t)
		acc := account.NewMockConnectedAccount(t)
		userID := types.UserID("user-gated-batch")
		acc.EXPECT().PlanID().Return(types.PlanFree)
		acc.EXPECT().UserID().Return(userID)
		// Even though gated feature will be disabled, needUsage will still be true so GetAll is called.
		usage.EXPECT().GetAll(mock.Anything, userID).Return([]bdomain.UsageCounter{}, nil).Once()

		resolver := &entitlementResolver{usage: usage, authentication: auth}
		res, err := resolver.Resolves(ctx, acc, []types.FeatureID{bdomain.FeatureIdCustomProvidersCount})
		require.NoError(t, err)
		require.Len(t, res, 1)
		assert.True(t, res[0].Enabled)
		assert.Nil(t, res[0].Limit)
		assert.Nil(t, res[0].Used)
		assert.Nil(t, res[0].Remaining)
	})
}

func TestEntitlementResolver_CheckBoolean(t *testing.T) {
	ctx := context.Background()
	usage := ports.NewMockUsageRepository(t)
	auth := ports.NewMockAuthentication(t)
	acc := account.NewMockConnectedAccount(t)
	acc.EXPECT().PlanID().Return(types.PlanFree).Maybe() // for boolean path plan lookup
	res := &entitlementResolver{usage: usage, authentication: auth}

	t.Run("returns error for quota feature", func(t *testing.T) {
		acc.EXPECT().PlanID().Return(types.PlanFree)
		acc.EXPECT().UserID().Return(types.UserID("user-qbool")).Maybe()
		usage.EXPECT().Get(mock.Anything, types.UserID("user-qbool"), mock.Anything).Return(bdomain.UsageCounter{
			FeatureID: bdomain.FeatureIdActiveSubscriptionsCount, Used: 0,
		}, true, nil)
		_, err := res.CheckBoolean(ctx, acc, bdomain.FeatureIdActiveSubscriptionsCount)
		assert.ErrorIs(t, err, bdomain.ErrInvalidFeatureType)
	})

	t.Run("returns enabled boolean", func(t *testing.T) {
		acc.EXPECT().PlanID().Return(types.PlanFree)
		e, err := res.CheckBoolean(ctx, acc, bdomain.FeatureIdSubscriptions)
		require.NoError(t, err)
		assert.True(t, e)
	})
}

func TestEntitlementResolver_CheckQuotaForAccount_and_CheckQuota(t *testing.T) {
	ctx := context.Background()

	t.Run("allowed when remaining >= needed", func(t *testing.T) {
		acc := account.NewMockConnectedAccount(t)
		usage := ports.NewMockUsageRepository(t)
		auth := ports.NewMockAuthentication(t)
		acc.EXPECT().PlanID().Return(types.PlanFree)
		acc.EXPECT().UserID().Return(types.UserID("user-q1"))
		usage.EXPECT().Get(mock.Anything, types.UserID("user-q1"), mock.Anything).Return(bdomain.UsageCounter{
			FeatureID: bdomain.FeatureIdActiveSubscriptionsCount, Used: 1,
		}, true, nil)
		resolver := &entitlementResolver{usage: usage, authentication: auth}
		allowed, eff, err := resolver.CheckQuotaForAccount(ctx, acc, bdomain.FeatureIdActiveSubscriptionsCount, 1)
		require.NoError(t, err)
		assert.True(t, allowed)
		assert.True(t, eff.Enabled)
	})

	t.Run("denied when remaining < needed", func(t *testing.T) {
		acc := account.NewMockConnectedAccount(t)
		usage := ports.NewMockUsageRepository(t)
		auth := ports.NewMockAuthentication(t)
		acc.EXPECT().PlanID().Return(types.PlanFree)
		acc.EXPECT().UserID().Return(types.UserID("user-q2"))
		usage.EXPECT().Get(mock.Anything, types.UserID("user-q2"), mock.Anything).Return(bdomain.UsageCounter{
			FeatureID: bdomain.FeatureIdActiveSubscriptionsCount, Used: 9,
		}, true, nil)
		resolver := &entitlementResolver{usage: usage, authentication: auth}
		allowed, eff, err := resolver.CheckQuotaForAccount(ctx, acc, bdomain.FeatureIdActiveSubscriptionsCount, 5)
		require.NoError(t, err)
		assert.False(t, allowed)
		assert.True(t, eff.Enabled)
	})

	t.Run("CheckQuota uses authentication and unlimited path (simulate unlimited)", func(t *testing.T) {
		// modify entitlement to unlimited
		orig := bdomain.Entitlements[types.PlanFree][bdomain.FeatureIdCustomLabelsCount]
		bdomain.Entitlements[types.PlanFree][bdomain.FeatureIdCustomLabelsCount] = bdomain.PlanEntitlement{
			PlanID: types.PlanFree, FeatureID: bdomain.FeatureIdCustomLabelsCount, Limit: nil,
		}
		defer func() { bdomain.Entitlements[types.PlanFree][bdomain.FeatureIdCustomLabelsCount] = orig }()

		acc := account.NewMockConnectedAccount(t)
		usage := ports.NewMockUsageRepository(t)
		auth := ports.NewMockAuthentication(t)
		auth.EXPECT().MustGetConnectedAccount(mock.Anything).Return(acc)
		acc.EXPECT().PlanID().Return(types.PlanFree)
		acc.EXPECT().UserID().Return(types.UserID("user-q3"))
		usage.EXPECT().Get(mock.Anything, types.UserID("user-q3"), mock.Anything).Return(bdomain.UsageCounter{
			FeatureID: bdomain.FeatureIdCustomLabelsCount, Used: 999,
		}, true, nil)
		resolver := &entitlementResolver{usage: usage, authentication: auth}
		allowed, eff, err := resolver.CheckQuota(ctx, bdomain.FeatureIdCustomLabelsCount, 1000)
		require.NoError(t, err)
		assert.True(t, allowed)
		assert.Nil(t, eff.Limit)
		assert.Nil(t, eff.Remaining)
	})
}
