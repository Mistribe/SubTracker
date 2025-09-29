package repositories

import (
	"context"
	"time"

	. "github.com/go-jet/jet/v2/postgres"

	"github.com/mistribe/subtracker/internal/adapters/persistence/db"
	. "github.com/mistribe/subtracker/internal/adapters/persistence/db/jet/app/public/table"
	"github.com/mistribe/subtracker/internal/domain/billing"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
)

type UsageRepository struct {
	dbContext *db.Context
}

func (u UsageRepository) GetAll(ctx context.Context, userID types.UserID) ([]billing.UsageCounter, error) {
	// Build reusable predicates
	now := time.Now()
	user := String(userID.String())

	// We'll reference Accounts via JOINs to safely access the user's family_id
	// Active subscriptions: owned by user or by user's family, and active at now
	activeOwnerPersonal := Subscriptions.OwnerType.EQ(String(types.PersonalOwnerType.String())).AND(Subscriptions.OwnerUserID.EQ(user))
	activeOwnerFamily := Subscriptions.OwnerType.EQ(String(types.FamilyOwnerType.String())).AND(Subscriptions.OwnerFamilyID.EQ(Accounts.FamilyID))
	activeOwnerCond := activeOwnerPersonal.OR(activeOwnerFamily)
	activeDateCond := Subscriptions.StartDate.LT_EQ(TimestampzT(now)).AND(Subscriptions.EndDate.IS_NULL().OR(Subscriptions.EndDate.GT(TimestampzT(now))))
	activeSubsCount := SELECT(COUNT(Subscriptions.ID)).
		FROM(Subscriptions.LEFT_JOIN(Accounts, Accounts.ID.EQ(user))).
		WHERE(activeOwnerCond.AND(activeDateCond))

	// Custom providers count: include personal and family owned
	provOwnerPersonal := Providers.OwnerType.EQ(String(types.PersonalOwnerType.String())).AND(Providers.OwnerUserID.EQ(user))
	provOwnerFamily := Providers.OwnerType.EQ(String(types.FamilyOwnerType.String())).AND(Providers.OwnerFamilyID.EQ(Accounts.FamilyID))
	provOwnerCond := provOwnerPersonal.OR(provOwnerFamily)
	customProvidersCount := SELECT(COUNT(Providers.ID)).
		FROM(Providers.LEFT_JOIN(Accounts, Accounts.ID.EQ(user))).
		WHERE(provOwnerCond)

	// Custom labels count: include personal and family owned
	labelOwnerPersonal := Labels.OwnerType.EQ(String(types.PersonalOwnerType.String())).AND(Labels.OwnerUserID.EQ(user))
	labelOwnerFamily := Labels.OwnerType.EQ(String(types.FamilyOwnerType.String())).AND(Labels.OwnerFamilyID.EQ(Accounts.FamilyID))
	labelOwnerCond := labelOwnerPersonal.OR(labelOwnerFamily)
	customLabelsCount := SELECT(COUNT(Labels.ID)).
		FROM(Labels.LEFT_JOIN(Accounts, Accounts.ID.EQ(user))).
		WHERE(labelOwnerCond)

	// Family members count: count all members within the user's family
	familyMembersCount := SELECT(COUNT(FamilyMembers.ID)).
		FROM(FamilyMembers.LEFT_JOIN(Accounts, Accounts.ID.EQ(user))).
		WHERE(FamilyMembers.FamilyID.EQ(Accounts.FamilyID))

	stmt := SELECT(
		activeSubsCount.AS("active_subscriptions_count"),
		customProvidersCount.AS("custom_providers_count"),
		customLabelsCount.AS("custom_labels_count"),
		familyMembersCount.AS("family_members_count"),
	)

	var row struct {
		ActiveSubscriptionsCount int64 `json:"active_subscriptions_count"`
		CustomProvidersCount     int64 `json:"custom_providers_count"`
		CustomLabelsCount        int64 `json:"custom_labels_count"`
		FamilyMembersCount       int64 `json:"family_members_count"`
	}

	if err := u.dbContext.Query(ctx, stmt, &row); err != nil {
		return nil, err
	}

	nowTs := time.Now()
	res := []billing.UsageCounter{
		{FeatureID: billing.FeatureIdActiveSubscriptionsCount, Used: row.ActiveSubscriptionsCount, UpdatedAt: nowTs},
		{FeatureID: billing.FeatureIdCustomProvidersCount, Used: row.CustomProvidersCount, UpdatedAt: nowTs},
		{FeatureID: billing.FeatureIdCustomLabelsCount, Used: row.CustomLabelsCount, UpdatedAt: nowTs},
		{FeatureID: billing.FeatureIdFamilyMembersCount, Used: row.FamilyMembersCount, UpdatedAt: nowTs},
	}
	return res, nil
}

func (u UsageRepository) Get(ctx context.Context, userID types.UserID, feature billing.Feature) (
	billing.UsageCounter,
	bool,
	error) {
	if !feature.IsQuota() {
		return billing.UsageCounter{}, false, billing.ErrCannotGetQuotaOnFeature
	}

	switch feature.ID {
	case billing.FeatureIdActiveSubscriptionsCount:
		return u.getActiveSubscriptionUsage(ctx, userID)
	case billing.FeatureIdCustomProvidersCount:
		return u.getCustomProviderUsage(ctx, userID)
	case billing.FeatureIdCustomLabelsCount:
		return u.getCustomLabelUsage(ctx, userID)
	case billing.FeatureIdFamilyMembersCount:
		return u.getFamilyMemberUsage(ctx, userID)
	default:
		return billing.UsageCounter{}, false, billing.ErrCannotGetQuotaOnFeature
	}
}

func (u UsageRepository) getActiveSubscriptionUsage(ctx context.Context, userID types.UserID) (
	billing.UsageCounter,
	bool,
	error) {
	now := time.Now()
	user := String(userID.String())

	ownerPersonal := Subscriptions.OwnerType.EQ(String(types.PersonalOwnerType.String())).AND(Subscriptions.OwnerUserID.EQ(user))
	ownerFamily := Subscriptions.OwnerType.EQ(String(types.FamilyOwnerType.String())).AND(Subscriptions.OwnerFamilyID.EQ(Accounts.FamilyID))
	ownerCond := ownerPersonal.OR(ownerFamily)
	dateCond := Subscriptions.StartDate.LT_EQ(TimestampzT(now)).AND(Subscriptions.EndDate.IS_NULL().OR(Subscriptions.EndDate.GT(TimestampzT(now))))

	stmt := SELECT(COUNT(Subscriptions.ID).AS("cnt")).
		FROM(Subscriptions.LEFT_JOIN(Accounts, Accounts.ID.EQ(user))).
		WHERE(ownerCond.AND(dateCond))

	var row struct{ Cnt int64 }
	if err := u.dbContext.Query(ctx, stmt, &row); err != nil {
		return billing.UsageCounter{}, false, err
	}
	return billing.UsageCounter{
		FeatureID: billing.FeatureIdActiveSubscriptionsCount, Used: row.Cnt, UpdatedAt: time.Now(),
	}, true, nil
}

func (u UsageRepository) getCustomProviderUsage(ctx context.Context, userID types.UserID) (
	billing.UsageCounter,
	bool,
	error) {
	user := String(userID.String())

	ownerPersonal := Providers.OwnerType.EQ(String(types.PersonalOwnerType.String())).AND(Providers.OwnerUserID.EQ(user))
	ownerFamily := Providers.OwnerType.EQ(String(types.FamilyOwnerType.String())).AND(Providers.OwnerFamilyID.EQ(Accounts.FamilyID))
	ownerCond := ownerPersonal.OR(ownerFamily)

	stmt := SELECT(COUNT(Providers.ID).AS("cnt")).
		FROM(Providers.LEFT_JOIN(Accounts, Accounts.ID.EQ(user))).
		WHERE(ownerCond)

	var row struct{ Cnt int64 }
	if err := u.dbContext.Query(ctx, stmt, &row); err != nil {
		return billing.UsageCounter{}, false, err
	}
	return billing.UsageCounter{
		FeatureID: billing.FeatureIdCustomProvidersCount, Used: row.Cnt, UpdatedAt: time.Now(),
	}, true, nil
}

func (u UsageRepository) getCustomLabelUsage(ctx context.Context, userID types.UserID) (
	billing.UsageCounter,
	bool,
	error) {
	user := String(userID.String())

	ownerPersonal := Labels.OwnerType.EQ(String(types.PersonalOwnerType.String())).AND(Labels.OwnerUserID.EQ(user))
	ownerFamily := Labels.OwnerType.EQ(String(types.FamilyOwnerType.String())).AND(Labels.OwnerFamilyID.EQ(Accounts.FamilyID))
	ownerCond := ownerPersonal.OR(ownerFamily)

	stmt := SELECT(COUNT(Labels.ID).AS("cnt")).
		FROM(Labels.LEFT_JOIN(Accounts, Accounts.ID.EQ(user))).
		WHERE(ownerCond)

	var row struct{ Cnt int64 }
	if err := u.dbContext.Query(ctx, stmt, &row); err != nil {
		return billing.UsageCounter{}, false, err
	}
	return billing.UsageCounter{
		FeatureID: billing.FeatureIdCustomLabelsCount, Used: row.Cnt, UpdatedAt: time.Now(),
	}, true, nil
}

func (u UsageRepository) getFamilyMemberUsage(ctx context.Context, userID types.UserID) (
	billing.UsageCounter,
	bool,
	error) {
	user := String(userID.String())

	stmt := SELECT(COUNT(FamilyMembers.ID).AS("cnt")).
		FROM(FamilyMembers.LEFT_JOIN(Accounts, Accounts.ID.EQ(user))).
		WHERE(FamilyMembers.FamilyID.EQ(Accounts.FamilyID))

	var row struct{ Cnt int64 }
	if err := u.dbContext.Query(ctx, stmt, &row); err != nil {
		return billing.UsageCounter{}, false, err
	}
	return billing.UsageCounter{
		FeatureID: billing.FeatureIdFamilyMembersCount, Used: row.Cnt, UpdatedAt: time.Now(),
	}, true, nil
}

func NewUsageRepository(dbContext *db.Context) ports.UsageRepository {
	return &UsageRepository{
		dbContext: dbContext,
	}
}
