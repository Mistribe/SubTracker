package authorization

import (
	"context"
	"errors"
	"testing"
	"time"

	accountDomain "github.com/mistribe/subtracker/internal/domain/account"
	domainAuth "github.com/mistribe/subtracker/internal/domain/authorization"
	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// ownedEntity implements ports.EntityWithOwnership for authorization tests.
type ownedEntity struct{ owner types.Owner }

func (e ownedEntity) Owner() types.Owner { return e.owner }

// newTestFamily builds a minimal concrete family instance.
func newTestFamily(famId types.FamilyID, owner types.UserID) family.Family {
	return family.NewFamily(famId, owner, "fam", []family.Member{}, time.Now(), time.Now())
}

// createService builds the service with provided collaborators.
func createService(ar ports.AccountRepository, fr ports.FamilyRepository, c ports.Cache, auth ports.Authentication) service {
	return service{userRepository: ar, familyRepository: fr, cache: c, authentication: auth}
}

// mockConnectedAccount sets required expectations on connected account.
func mockConnectedAccount(t *testing.T, userId types.UserID, role types.Role) accountDomain.ConnectedAccount {
	ca := accountDomain.NewMockConnectedAccount(t)
	ca.EXPECT().UserID().Return(userId)
	ca.EXPECT().Role().Return(role)
	return ca
}

// newCacheWithLevel sets up a Cache mock returning the given leveled cache a fixed number of times.
func newCacheWithLevel(t *testing.T, leveled *ports.MockCacheLeveled, fromCalls int) *ports.MockCache {
	c := ports.NewMockCache(t)
	for i := 0; i < fromCalls; i++ {
		c.EXPECT().From(mock.Anything, ports.CacheLevelRequest).Return(leveled)
	}
	return c
}

func TestService_getFamilyId(t *testing.T) {
	ctx := context.Background()
	userId := types.UserID("user-1")

	t.Run("cache hit returns value without repository call", func(t *testing.T) {
		cacheKey := "family-user-" + userId.String()
		famId := types.NewFamilyID()
		leveled := ports.NewMockCacheLeveled(t)
		leveled.EXPECT().Get(cacheKey).Return(famId)
		cache := newCacheWithLevel(t, leveled, 1)
		fr := ports.NewMockFamilyRepository(t) // should not be called

		s := createService(nil, fr, cache, nil)
		got, err := s.getFamilyId(ctx, userId)
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, famId, *got)
	})

	t.Run("cache miss repository returns family and caches it", func(t *testing.T) {
		cacheKey := "family-user-" + userId.String()
		famId := types.NewFamilyID()
		leveled := ports.NewMockCacheLeveled(t)
		leveled.EXPECT().Get(cacheKey).Return(nil)   // miss
		leveled.EXPECT().Set(cacheKey, famId)        // write
		leveled.EXPECT().Get(cacheKey).Return(famId) // second call after caching
		cache := newCacheWithLevel(t, leveled, 3)    // three From() invocations
		fr := ports.NewMockFamilyRepository(t)
		fr.EXPECT().GetAccountFamily(mock.Anything, userId).Return(newTestFamily(famId, userId), nil)

		s := createService(nil, fr, cache, nil)
		got, err := s.getFamilyId(ctx, userId)
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, famId, *got)

		// second call hits cache
		got2, err2 := s.getFamilyId(ctx, userId)
		require.NoError(t, err2)
		require.NotNil(t, got2)
		assert.Equal(t, famId, *got2)
	})

	t.Run("cache miss repository error propagates", func(t *testing.T) {
		cacheKey := "family-user-" + userId.String()
		leveled := ports.NewMockCacheLeveled(t)
		leveled.EXPECT().Get(cacheKey).Return(nil)
		cache := newCacheWithLevel(t, leveled, 1)
		fr := ports.NewMockFamilyRepository(t)
		expectedErr := errors.New("db failure")
		fr.EXPECT().GetAccountFamily(mock.Anything, userId).Return(nil, expectedErr)

		s := createService(nil, fr, cache, nil)
		got, err := s.getFamilyId(ctx, userId)
		require.Error(t, err)
		assert.Nil(t, got)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("cache miss repository returns nil family", func(t *testing.T) {
		cacheKey := "family-user-" + userId.String()
		leveled := ports.NewMockCacheLeveled(t)
		leveled.EXPECT().Get(cacheKey).Return(nil)
		cache := newCacheWithLevel(t, leveled, 1)
		fr := ports.NewMockFamilyRepository(t)
		fr.EXPECT().GetAccountFamily(mock.Anything, userId).Return(nil, nil)

		s := createService(nil, fr, cache, nil)
		got, err := s.getFamilyId(ctx, userId)
		require.NoError(t, err)
		assert.Nil(t, got)
	})
}

func TestPermissionRequest_For(t *testing.T) {
	userId := types.UserID("user-2")
	famId := types.NewFamilyID()

	baseReq := func(role types.Role, family *types.FamilyID, err error) permissionRequest {
		return permissionRequest{userId: userId, userRole: role, userFamily: family, error: err}
	}

	newEntity := func(owner types.Owner) ownedEntity { return ownedEntity{owner: owner} }

	t.Run("error in request returned immediately", func(t *testing.T) {
		pr := baseReq(types.RoleUser, nil, errors.New("boom"))
		err := pr.For(newEntity(types.SystemOwner))
		assert.EqualError(t, err, "boom")
	})

	t.Run("admin role always authorized", func(t *testing.T) {
		pr := baseReq(types.RoleAdmin, nil, nil)
		err := pr.For(newEntity(types.SystemOwner))
		assert.NoError(t, err)
	})

	t.Run("system owner unauthorized for non-admin", func(t *testing.T) {
		pr := baseReq(types.RoleUser, nil, nil)
		err := pr.For(newEntity(types.SystemOwner))
		assert.ErrorIs(t, err, domainAuth.ErrUnauthorized)
	})

	t.Run("family owner mismatch", func(t *testing.T) {
		otherFam := types.NewFamilyID()
		pr := baseReq(types.RoleUser, &famId, nil)
		entity := newEntity(types.NewFamilyOwner(otherFam))
		err := pr.For(entity)
		assert.ErrorIs(t, err, domainAuth.ErrUnauthorized)
	})

	t.Run("family owner nil family in request => unauthorized", func(t *testing.T) {
		pr := baseReq(types.RoleUser, nil, nil)
		entity := newEntity(types.NewFamilyOwner(famId))
		err := pr.For(entity)
		assert.ErrorIs(t, err, domainAuth.ErrUnauthorized)
	})

	t.Run("family owner match", func(t *testing.T) {
		pr := baseReq(types.RoleUser, &famId, nil)
		entity := newEntity(types.NewFamilyOwner(famId))
		err := pr.For(entity)
		assert.NoError(t, err)
	})

	t.Run("personal owner mismatch", func(t *testing.T) {
		pr := baseReq(types.RoleUser, &famId, nil)
		entity := newEntity(types.NewPersonalOwner(types.UserID("other-user")))
		err := pr.For(entity)
		assert.ErrorIs(t, err, domainAuth.ErrUnauthorized)
	})

	t.Run("personal owner match", func(t *testing.T) {
		pr := baseReq(types.RoleUser, &famId, nil)
		entity := newEntity(types.NewPersonalOwner(userId))
		err := pr.For(entity)
		assert.NoError(t, err)
	})
}

func TestService_Can_BuildsPermissionRequest(t *testing.T) {
	ctx := context.Background()
	userId := types.UserID("user-3")
	famId := types.NewFamilyID()
	perm := domainAuth.PermissionRead

	cacheKey := "family-user-" + userId.String()
	leveled := ports.NewMockCacheLeveled(t)
	leveled.EXPECT().Get(cacheKey).Return(nil)
	leveled.EXPECT().Set(cacheKey, famId)
	cache := newCacheWithLevel(t, leveled, 2)

	fr := ports.NewMockFamilyRepository(t)
	fr.EXPECT().GetAccountFamily(mock.Anything, userId).Return(newTestFamily(famId, userId), nil)

	authMock := ports.NewMockAuthentication(t)
	connected := mockConnectedAccount(t, userId, types.RoleUser)
	authMock.EXPECT().MustGetConnectedAccount(mock.Anything).Return(connected)

	s := createService(nil, fr, cache, authMock)
	req := s.Can(ctx, perm)
	pr, ok := req.(*permissionRequest)
	require.True(t, ok)
	assert.Equal(t, userId, pr.userId)
	assert.Equal(t, types.RoleUser, pr.userRole)
	assert.Equal(t, perm, pr.permission)
	require.NotNil(t, pr.userFamily)
	assert.Equal(t, famId, *pr.userFamily)
	assert.NoError(t, pr.error)
}
