package query_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/family/query"
)

// helpers (duplicated minimal versions since command test helpers live in another package)
func helperNewMember(famID types.FamilyID, name string, mt family.MemberType) family.Member {
	return family.NewMember(types.NewFamilyMemberID(), famID, name, mt, nil, time.Now(), time.Now())
}

func helperNewFamily(owner types.UserID, name string, members []family.Member) family.Family {
	return family.NewFamily(types.NewFamilyID(), owner, name, members, time.Now(), time.Now())
}

func TestFindUserFamilyQueryHandler_Handle(t *testing.T) {
	ctx := context.Background()

	t.Run("returns repository error", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		handler := query.NewFindOneQueryHandler(repo, authz)

		expectedErr := errors.New("db failure")
		uid := types.UserID("user-1")
		repo.EXPECT().GetAccountFamily(mock.Anything, uid).Return(nil, expectedErr)

		res := handler.Handle(ctx, query.FindUserFamilyQuery{UserID: uid})
		require.True(t, res.IsFaulted())
		var gotErr error
		res.IfFailure(func(e error) { gotErr = e })
		assert.ErrorIs(t, gotErr, expectedErr)
	})

	t.Run("returns ErrFamilyNotFound when repository returns nil without error", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		handler := query.NewFindOneQueryHandler(repo, authz)

		uid := types.UserID("user-2")
		repo.EXPECT().GetAccountFamily(mock.Anything, uid).Return(nil, nil)

		res := handler.Handle(ctx, query.FindUserFamilyQuery{UserID: uid})
		require.True(t, res.IsFaulted())
		var gotErr error
		res.IfFailure(func(e error) { gotErr = e })
		assert.ErrorIs(t, gotErr, family.ErrFamilyNotFound)
	})

	t.Run("succeeds returning family", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		handler := query.NewFindOneQueryHandler(repo, authz)

		uid := types.UserID("user-3")
		fam := helperNewFamily(uid, "TestFam", []family.Member{helperNewMember(types.NewFamilyID(), "Owner", family.OwnerMemberType)})
		repo.EXPECT().GetAccountFamily(mock.Anything, uid).Return(fam, nil)

		res := handler.Handle(ctx, query.FindUserFamilyQuery{UserID: uid})
		require.True(t, res.IsSuccess())
		var response query.FindUserFamilyQueryResponse
		res.IfSuccess(func(r query.FindUserFamilyQueryResponse) { response = r })
		require.NotNil(t, response.Family)
		assert.Equal(t, fam.Id(), response.Family.Id())
	})
}
