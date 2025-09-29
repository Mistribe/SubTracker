package query_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/mistribe/subtracker/internal/domain/account"
	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/family/query"
)

// local helpers (cannot reuse command package helpers due to different package)
func helperNewMemberWithParams(
	memberID types.FamilyMemberID,
	famID types.FamilyID,
	name string,
	mt family.MemberType,
	invitation *string,
) family.Member {
	return family.NewMember(memberID, famID, name, mt, invitation, time.Now(), time.Now())
}

func helperBuildFamily(
	famID types.FamilyID,
	owner types.UserID,
	name string,
	members []family.Member,
) family.Family {
	return family.NewFamily(famID, owner, name, members, time.Now(), time.Now())
}

func TestSeeInvitationQueryHandler_Handle(t *testing.T) {
	ctx := context.Background()

	t.Run("propagates repository error", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		auth := ports.NewMockAuthentication(t)
		handler := query.NewSeeInvitationQueryHandler(repo, auth)

		famID := types.NewFamilyID()
		memberID := types.NewFamilyMemberID()
		q := query.SeeInvitationQuery{FamilyId: famID, FamilyMemberId: memberID, InvitationCode: "CODE"}
		expectedErr := errors.New("db failure")
		repo.EXPECT().GetById(mock.Anything, famID).Return(nil, expectedErr)

		res := handler.Handle(ctx, q)
		require.True(t, res.IsFaulted())
		var got error
		res.IfFailure(func(e error) { got = e })
		assert.ErrorIs(t, got, expectedErr)
	})

	t.Run("returns ErrFamilyNotFound when repo returns nil family", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		auth := ports.NewMockAuthentication(t)
		handler := query.NewSeeInvitationQueryHandler(repo, auth)

		famID := types.NewFamilyID()
		memberID := types.NewFamilyMemberID()
		q := query.SeeInvitationQuery{FamilyId: famID, FamilyMemberId: memberID, InvitationCode: "CODE"}
		repo.EXPECT().GetById(mock.Anything, famID).Return(nil, nil)

		res := handler.Handle(ctx, q)
		require.True(t, res.IsFaulted())
		var got error
		res.IfFailure(func(e error) { got = e })
		assert.ErrorIs(t, got, family.ErrFamilyNotFound)
	})

	t.Run("returns ErrFamilyMemberNotFound when member not present", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		auth := ports.NewMockAuthentication(t)
		handler := query.NewSeeInvitationQueryHandler(repo, auth)

		famID := types.NewFamilyID()
		owner := types.UserID("owner-user")
		// Family has zero members
		fam := helperBuildFamily(famID, owner, "Fam", []family.Member{})
		memberID := types.NewFamilyMemberID()
		q := query.SeeInvitationQuery{FamilyId: famID, FamilyMemberId: memberID, InvitationCode: "CODE"}
		repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)

		res := handler.Handle(ctx, q)
		require.True(t, res.IsFaulted())
		var got error
		res.IfFailure(func(e error) { got = e })
		assert.ErrorIs(t, got, family.ErrFamilyMemberNotFound)
	})

	t.Run("returns ErrBadInvitationCode when member has nil invitation code", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		auth := ports.NewMockAuthentication(t)
		handler := query.NewSeeInvitationQueryHandler(repo, auth)

		famID := types.NewFamilyID()
		owner := types.UserID("owner-user")
		memberID := types.NewFamilyMemberID()
		mem := helperNewMemberWithParams(memberID, famID, "John", family.KidMemberType, nil)
		fam := helperBuildFamily(famID, owner, "Fam", []family.Member{mem})
		q := query.SeeInvitationQuery{FamilyId: famID, FamilyMemberId: memberID, InvitationCode: "CODE"}
		repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)

		res := handler.Handle(ctx, q)
		require.True(t, res.IsFaulted())
		var got error
		res.IfFailure(func(e error) { got = e })
		assert.ErrorIs(t, got, family.ErrBadInvitationCode)
	})

	t.Run("returns ErrBadInvitationCode when codes mismatch", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		auth := ports.NewMockAuthentication(t)
		handler := query.NewSeeInvitationQueryHandler(repo, auth)

		famID := types.NewFamilyID()
		owner := types.UserID("owner-user")
		memberID := types.NewFamilyMemberID()
		code := "ABC123"
		mem := helperNewMemberWithParams(memberID, famID, "John", family.KidMemberType, &code)
		fam := helperBuildFamily(famID, owner, "Fam", []family.Member{mem})
		q := query.SeeInvitationQuery{FamilyId: famID, FamilyMemberId: memberID, InvitationCode: "OTHER"}
		repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)

		res := handler.Handle(ctx, q)
		require.True(t, res.IsFaulted())
		var got error
		res.IfFailure(func(e error) { got = e })
		assert.ErrorIs(t, got, family.ErrBadInvitationCode)
	})

	t.Run("succeeds when invitation code matches", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		auth := ports.NewMockAuthentication(t)
		handler := query.NewSeeInvitationQueryHandler(repo, auth)

		famID := types.NewFamilyID()
		owner := types.UserID("owner-user")
		memberID := types.NewFamilyMemberID()
		code := "MATCH123"
		mem := helperNewMemberWithParams(memberID, famID, "John", family.KidMemberType, &code)
		fam := helperBuildFamily(famID, owner, "Fam", []family.Member{mem})
		q := query.SeeInvitationQuery{FamilyId: famID, FamilyMemberId: memberID, InvitationCode: code}
		repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)

		connected := account.NewMockConnectedAccount(t)
		connected.EXPECT().UserID().Return(types.UserID("current-user"))
		auth.EXPECT().MustGetConnectedAccount(mock.Anything).Return(connected)

		res := handler.Handle(ctx, q)
		require.True(t, res.IsSuccess())
		var response query.SeeInvitationQueryResponse
		res.IfSuccess(func(r query.SeeInvitationQueryResponse) { response = r })
		assert.Equal(t, fam.Id(), response.Family.Id())
		assert.Equal(t, memberID, response.InvitedInasmuchAs)
		assert.Equal(t, types.UserID("current-user"), response.UserId)
	})
}
