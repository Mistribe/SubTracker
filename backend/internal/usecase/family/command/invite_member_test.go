package command_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/mistribe/subtracker/internal/domain/authorization"
	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/family/command"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

func TestInviteMemberCommandHandler(t *testing.T) {
	ctx := context.Background()

	t.Run("GetById returns error", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		h := command.NewInviteMemberCommandHandler(repo, authz)

		famID := types.NewFamilyID()
		cmd := command.InviteMemberCommand{FamilyId: famID, FamilyMemberId: types.NewFamilyMemberID()}

		repo.EXPECT().GetById(mock.Anything, famID).Return(nil, errors.New("db error"))

		res := h.Handle(ctx, cmd)
		require.True(t, res.IsFaulted())
		// authorization not reached
		authz.AssertNotCalled(t, "Can", mock.Anything, mock.Anything)
	})

	t.Run("Family not found returns ErrFamilyNotFound", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		permReq := ports.NewMockPermissionRequest(t)
		h := command.NewInviteMemberCommandHandler(repo, authz)

		famID := types.NewFamilyID()
		cmd := command.InviteMemberCommand{FamilyId: famID, FamilyMemberId: types.NewFamilyMemberID()}

		repo.EXPECT().GetById(mock.Anything, famID).Return(nil, nil)

		res := h.Handle(ctx, cmd)
		err := result.Match(res, func(_ command.InviteMemberResponse) error { return nil },
			func(err error) error { return err })
		require.Error(t, err)
		assert.True(t, errors.Is(err, family.ErrFamilyNotFound))
		authz.AssertNotCalled(t, "Can", mock.Anything, mock.Anything)
		permReq.AssertNotCalled(t, "For", mock.Anything)
	})

	t.Run("Authorization denies returns error", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		permReq := ports.NewMockPermissionRequest(t)
		h := command.NewInviteMemberCommandHandler(repo, authz)

		famID := types.NewFamilyID()
		fam := helperNewFamilyWithID(famID, types.UserID("owner-1"), "Doe family", nil)
		cmd := command.InviteMemberCommand{FamilyId: famID, FamilyMemberId: types.NewFamilyMemberID()}

		repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
		authz.EXPECT().Can(mock.Anything, authorization.PermissionWrite).Return(permReq)
		expected := errors.New("forbidden")
		permReq.EXPECT().For(fam).Return(expected)

		res := h.Handle(ctx, cmd)
		err := result.Match(res, func(_ command.InviteMemberResponse) error { return nil },
			func(err error) error { return err })
		require.Error(t, err)
		assert.True(t, errors.Is(err, expected))
		repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	t.Run("Existing member with user already set returns ErrCannotInviteUser", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		permReq := ports.NewMockPermissionRequest(t)
		h := command.NewInviteMemberCommandHandler(repo, authz)

		famID := types.NewFamilyID()
		m := helperNewMember(famID, "John", family.AdultMemberType)
		uid := types.UserID("user-1")
		m.SetUserId(&uid)
		fam := helperNewFamilyWithID(famID, types.UserID("owner-1"), "Doe family", []family.Member{m})
		cmd := command.InviteMemberCommand{FamilyId: famID, FamilyMemberId: m.Id()}

		repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
		authz.EXPECT().Can(mock.Anything, authorization.PermissionWrite).Return(permReq)
		permReq.EXPECT().For(fam).Return(nil)

		res := h.Handle(ctx, cmd)
		err := result.Match(res, func(_ command.InviteMemberResponse) error { return nil },
			func(err error) error { return err })
		require.Error(t, err)
		assert.True(t, errors.Is(err, family.ErrCannotInviteUser))
		repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	t.Run("Existing member invite success updates code and saves", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		permReq := ports.NewMockPermissionRequest(t)
		h := command.NewInviteMemberCommandHandler(repo, authz)

		famID := types.NewFamilyID()
		m := helperNewMember(famID, "John", family.AdultMemberType)
		fam := helperNewFamilyWithID(famID, types.UserID("owner-1"), "Doe family", []family.Member{m})
		cmd := command.InviteMemberCommand{FamilyId: famID, FamilyMemberId: m.Id()}

		repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
		authz.EXPECT().Can(mock.Anything, authorization.PermissionWrite).Return(permReq)
		permReq.EXPECT().For(fam).Return(nil)
		repo.EXPECT().Save(mock.Anything, mock.Anything).Return(nil)

		res := h.Handle(ctx, cmd)

		require.True(t, res.IsSuccess())
		out := result.Match(res, func(v command.InviteMemberResponse) command.InviteMemberResponse { return v },
			func(err error) command.InviteMemberResponse {
				t.Fatalf("unexpected error: %v",
					err)
				return command.InviteMemberResponse{}
			})
		require.NotEmpty(t, out.Code)
		updated := fam.GetMember(m.Id())
		require.NotNil(t, updated)
		require.NotNil(t, updated.InvitationCode())
		assert.Equal(t, out.Code, *updated.InvitationCode())
	})

	t.Run("New member invite success with default name and type", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		permReq := ports.NewMockPermissionRequest(t)
		h := command.NewInviteMemberCommandHandler(repo, authz)

		famID := types.NewFamilyID()
		fam := helperNewFamilyWithID(famID, types.UserID("owner-1"), "Doe family", nil)
		cmd := command.InviteMemberCommand{
			FamilyId: famID, FamilyMemberId: types.NewFamilyMemberID(), Name: nil, Type: nil,
		}

		repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
		authz.EXPECT().Can(mock.Anything, authorization.PermissionWrite).Return(permReq)
		permReq.EXPECT().For(fam).Return(nil)
		repo.EXPECT().Save(mock.Anything, mock.Anything).Return(nil)

		res := h.Handle(ctx, cmd)

		require.True(t, res.IsSuccess())
		out := result.Match(res, func(v command.InviteMemberResponse) command.InviteMemberResponse { return v },
			func(err error) command.InviteMemberResponse {
				t.Fatalf("unexpected error: %v",
					err)
				return command.InviteMemberResponse{}
			})
		created := fam.GetMember(out.FamilyMemberId)
		require.NotNil(t, created)
		assert.Equal(t, "Invited User", created.Name())
		assert.Equal(t, family.AdultMemberType, created.Type())
		require.NotNil(t, created.InvitationCode())
		assert.Equal(t, out.Code, *created.InvitationCode())
	})

	t.Run("Existing member: validation error prevents Save", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		permReq := ports.NewMockPermissionRequest(t)
		h := command.NewInviteMemberCommandHandler(repo, authz)

		famID := types.NewFamilyID()
		m := helperNewMember(famID, "John", family.AdultMemberType)
		fam := helperNewFamilyWithID(famID, types.UserID("owner-1"), "x", []family.Member{m})
		cmd := command.InviteMemberCommand{FamilyId: famID, FamilyMemberId: m.Id()}

		repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
		authz.EXPECT().Can(mock.Anything, authorization.PermissionWrite).Return(permReq)
		permReq.EXPECT().For(fam).Return(nil)

		res := h.Handle(ctx, cmd)

		require.True(t, res.IsFaulted())
		repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	t.Run("Existing member: repository Save error propagates", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		permReq := ports.NewMockPermissionRequest(t)
		h := command.NewInviteMemberCommandHandler(repo, authz)

		famID := types.NewFamilyID()
		m := helperNewMember(famID, "John", family.AdultMemberType)
		fam := helperNewFamilyWithID(famID, types.UserID("owner-1"), "Doe family", []family.Member{m})
		cmd := command.InviteMemberCommand{FamilyId: famID, FamilyMemberId: m.Id()}

		repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
		authz.EXPECT().Can(mock.Anything, authorization.PermissionWrite).Return(permReq)
		permReq.EXPECT().For(fam).Return(nil)
		repo.EXPECT().Save(mock.Anything, mock.Anything).Return(errors.New("save failed"))

		res := h.Handle(ctx, cmd)
		require.True(t, res.IsFaulted())
	})
}
