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
	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/domain/types"
)

func TestFamilyRepository_CRUDAndMembers(t *testing.T) {
	ctx := context.Background()
	repo := repositories.NewFamilyRepository(GetDBContext())

	ownerUserID := types.UserID("user-" + uuid.NewString())
	familyID := types.NewFamilyID()
	member1 := family.NewMember(
		types.NewFamilyMemberID(),
		familyID,
		"Alice",
		family.OwnerMemberType,
		nil,
		time.Now().UTC(),
		time.Now().UTC(),
	)
	member1User := ownerUserID
	member1.SetUserId(&member1User)

	fam := family.NewFamily(
		familyID,
		ownerUserID,
		"Test Family "+uuid.NewString()[0:8],
		[]family.Member{member1},
		time.Now().UTC(),
		time.Now().UTC(),
	)

	// Create
	require.NoError(t, repo.Save(ctx, fam))

	// GetById
	stored, err := repo.GetById(ctx, fam.Id())
	require.NoError(t, err)
	require.NotNil(t, stored)
	assert.Equal(t, fam.Id(), stored.Id())
	assert.Equal(t, 1, stored.Members().Len())

	// IsUserMemberOfFamily
	isMember, err := repo.IsUserMemberOfFamily(ctx, fam.Id(), ownerUserID)
	require.NoError(t, err)
	assert.True(t, isMember)

	// Add new member and update
	member2 := family.NewMember(
		types.NewFamilyMemberID(),
		familyID,
		"Bob",
		family.AdultMemberType,
		nil,
		time.Now().UTC(),
		time.Now().UTC(),
	)
	require.NoError(t, fam.AddMember(member2))
	require.NoError(t, repo.Save(ctx, fam))
	stored, err = repo.GetById(ctx, fam.Id())
	require.NoError(t, err)
	assert.Equal(t, 2, stored.Members().Len())

	// MemberExists
	exists, err := repo.MemberExists(ctx, fam.Id(), member1.Id(), member2.Id())
	require.NoError(t, err)
	assert.True(t, exists)

	// Exists
	familyExists, err := repo.Exists(ctx, fam.Id())
	require.NoError(t, err)
	assert.True(t, familyExists)

	// Delete
	deleted, err := repo.Delete(ctx, fam.Id())
	require.NoError(t, err)
	assert.True(t, deleted)

	familyExists, err = repo.Exists(ctx, fam.Id())
	require.NoError(t, err)
	assert.False(t, familyExists)
}
