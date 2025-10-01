//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	xcur "golang.org/x/text/currency"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mistribe/subtracker/internal/adapters/persistence/repositories"
	"github.com/mistribe/subtracker/internal/domain/currency"
	"github.com/mistribe/subtracker/internal/domain/types"
)

func TestCurrencyRateRepository_CRUD(t *testing.T) {
	ctx := context.Background()
	repo := repositories.NewCurrencyRateRepository(GetDBContext())

	rateDate := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)
	rate := currency.NewRate(
		types.NewRateID(),
		xcur.MustParseISO("EUR"),
		xcur.MustParseISO("USD"),
		rateDate,
		1.2345,
		time.Now().UTC(),
		time.Now().UTC(),
	)

	// Create
	require.NoError(t, repo.Save(ctx, rate))

	// GetById
	stored, err := repo.GetById(ctx, rate.Id())
	require.NoError(t, err)
	require.NotNil(t, stored)
	assert.Equal(t, rate.Id(), stored.Id())
	assert.Equal(t, rate.ETag(), stored.ETag())

	// Exists
	exists, err := repo.Exists(ctx, rate.Id())
	require.NoError(t, err)
	assert.True(t, exists)

	// Latest update date (should be >= created at)
	latest, err := repo.GetLatestUpdateDate(ctx)
	require.NoError(t, err)
	assert.False(t, latest.IsZero())

	// Delete
	deleted, err := repo.Delete(ctx, rate.Id())
	require.NoError(t, err)
	assert.True(t, deleted)

	// Exists after delete
	exists, err = repo.Exists(ctx, rate.Id())
	require.NoError(t, err)
	assert.False(t, exists)
}

func TestCurrencyRateRepository_GetRatesByDate(t *testing.T) {
	ctx := context.Background()
	repo := repositories.NewCurrencyRateRepository(GetDBContext())

	rateDate := time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)
	baseTime := time.Now().UTC()

	r1 := currency.NewRate(types.NewRateID(), xcur.MustParseISO("EUR"), xcur.MustParseISO("USD"), rateDate, 1.10,
		baseTime, baseTime)
	r2 := currency.NewRate(types.NewRateID(), xcur.MustParseISO("EUR"), xcur.MustParseISO("GBP"), rateDate, 0.85,
		baseTime, baseTime)
	r3 := currency.NewRate(types.NewRateID(), xcur.MustParseISO("USD"), xcur.MustParseISO("JPY"), rateDate, 150.0,
		baseTime, baseTime)

	require.NoError(t, repo.Save(ctx, r1, r2, r3))

	rates, err := repo.GetRatesByDate(ctx, rateDate)
	require.NoError(t, err)
	require.Len(t, rates, 3)

	// Validate one of them
	found, ok := rates.FindExchangeRate(xcur.MustParseISO("EUR"), xcur.MustParseISO("USD"))
	assert.True(t, ok)
	assert.InDelta(t, 1.10, found, 0.00001)
}
