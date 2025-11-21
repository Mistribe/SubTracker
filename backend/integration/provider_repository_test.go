//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mistribe/subtracker/internal/adapters/persistence/repositories"
	"github.com/mistribe/subtracker/internal/domain/provider"
	"github.com/mistribe/subtracker/internal/domain/types"
)

func TestProviderRepository_CRUD(t *testing.T) {
	ctx := context.Background()
	repo := repositories.NewProviderRepository(GetDBContext())

	prov := provider.NewProvider(
		types.NewProviderID(),
		"Provider "+uuid.NewString()[0:8],
		nil, // key
		nil, // description
		nil, // icon
		nil, // url
		[]types.LabelID{},
		types.SystemOwner,
		time.Now().UTC(),
		time.Now().UTC(),
	)

	require.NoError(t, repo.Save(ctx, prov))

	stored, err := repo.GetById(ctx, prov.Id())
	require.NoError(t, err)
	require.NotNil(t, stored)
	assert.Equal(t, prov.Id(), stored.Id())

	exists, err := repo.Exists(ctx, prov.Id())
	require.NoError(t, err)
	assert.True(t, exists)

	// Update provider name
	newName := "Updated-" + uuid.NewString()[0:6]
	stored.SetName(newName)
	require.NoError(t, repo.Save(ctx, stored))

	updated, err := repo.GetById(ctx, prov.Id())
	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, newName, updated.Name())

	systems, total, err := repo.GetSystemProviders(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(1))
	found := false
	for _, p := range systems {
		if p.Id() == prov.Id() {
			found = true
			break
		}
	}
	assert.True(t, found, "saved system provider should be in system list")

	deleted, err := repo.Delete(ctx, prov.Id())
	require.NoError(t, err)
	assert.True(t, deleted)

	exists, err = repo.Exists(ctx, prov.Id())
	require.NoError(t, err)
	assert.False(t, exists)
}
