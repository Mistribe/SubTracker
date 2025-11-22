//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	xcur "golang.org/x/text/currency"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mistribe/subtracker/internal/adapters/persistence/repositories"
	"github.com/mistribe/subtracker/internal/domain/currency"
	"github.com/mistribe/subtracker/internal/domain/family"
	providerDomain "github.com/mistribe/subtracker/internal/domain/provider"
	subdom "github.com/mistribe/subtracker/internal/domain/subscription"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
)

func TestSubscriptionRepository_CRUD(t *testing.T) {
	ctx := context.Background()
	provRepo := repositories.NewProviderRepository(GetDBContext())
	subRepo := repositories.NewSubscriptionRepository(GetDBContext())
	famRepo := repositories.NewFamilyRepository(GetDBContext())
	familyId := types.NewFamilyID()

	// create family
	ownerUserID := types.UserID(uuid.NewString())
	m := family.NewMember(
		types.NewFamilyMemberID(),
		familyId,
		"Owner",
		family.OwnerMemberType,
		nil,
		time.Now().UTC(),
		time.Now().UTC(),
	)
	fam := family.NewFamily(
		familyId,
		ownerUserID,
		"Fam-"+uuid.NewString()[0:8],
		[]family.Member{m},
		time.Now().UTC(),
		time.Now().UTC(),
	)
	require.NoError(t, famRepo.Save(ctx, fam))

	// Create a provider (system owner) required by subscription FK
	prov := providerDomain.NewProvider(
		types.NewProviderID(),
		"Prov-"+uuid.NewString()[0:8],
		nil,
		nil,
		nil,
		nil,
		[]types.LabelID{},
		types.SystemOwner,
		time.Now().UTC(),
		time.Now().UTC(),
	)
	require.NoError(t, provRepo.Save(ctx, prov))

	price := currency.NewAmount(9.99, xcur.MustParseISO("USD"))
	priceWrapper := subdom.NewPrice(price)
	friendly := "My Subscription"

	sub := subdom.NewSubscription(
		types.NewSubscriptionID(),
		&friendly,
		nil, // free trial
		prov.Id(),
		priceWrapper,
		types.NewFamilyOwner(familyId), // system owned subscription
		nil,                            // payer
		[]types.FamilyMemberID{},
		[]subdom.LabelRef{},
		time.Now().Add(-24*time.Hour).UTC(), // start yesterday
		nil,                                 // no end date
		subdom.MonthlyRecurrency,
		nil, // custom recurrency
		time.Now().UTC(),
		time.Now().UTC(),
	)

	require.NoError(t, subRepo.Save(ctx, sub))

	// GetById
	stored, err := subRepo.GetById(ctx, sub.Id())
	require.NoError(t, err)
	require.NotNil(t, stored)
	assert.Equal(t, sub.Id(), stored.Id())
	assert.True(t, stored.IsActive())

	// Exists
	exists, err := subRepo.Exists(ctx, sub.Id())
	require.NoError(t, err)
	assert.True(t, exists)

	// Pagination (GetAll)
	params := ports.NewSubscriptionQueryParameters(
		"",                        // searchText
		[]subdom.RecurrencyType{}, // recurrencies
		nil, nil,                  // from, to
		[]types.UserID{},     // users
		[]types.ProviderID{}, // providers
		true,                 // withInactive
		50, 0,                // limit, offset
	)
	list, total, err := subRepo.GetAll(ctx, params)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(1))
	found := false
	for _, s := range list {
		if s.Id() == sub.Id() {
			found = true
			break
		}
	}
	assert.True(t, found, "subscription should appear in GetAll list")

	// Delete subscription
	deleted, err := subRepo.Delete(ctx, sub.Id())
	require.NoError(t, err)
	assert.True(t, deleted)

	exists, err = subRepo.Exists(ctx, sub.Id())
	require.NoError(t, err)
	assert.False(t, exists)

	// Cleanup provider
	provDeleted, err := provRepo.Delete(ctx, prov.Id())
	require.NoError(t, err)
	assert.True(t, provDeleted)

	// Cleanup family
	famDeleted, err := famRepo.Delete(ctx, familyId)
	require.NoError(t, err)
	assert.True(t, famDeleted)
}
