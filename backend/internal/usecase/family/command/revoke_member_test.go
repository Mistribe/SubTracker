package command_test

import (
	"context"
	"errors"
	"testing"
	"time"

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

func TestRevokeMemberCommandHandler(t *testing.T) {
	ctx := context.Background()

	t.Run("GetById returns error", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		h := command.NewRevokeMemberCommandHandler(repo, authz)

		famID := types.NewFamilyID()
		cmd := command.RevokeMemberCommand{FamilyId: famID, FamilyMemberId: types.NewFamilyMemberID()}

		expected := errors.New("db error")
		repo.EXPECT().GetById(mock.Anything, famID).Return(nil, expected)

		res := h.Handle(ctx, cmd)

		err := result.Match(res, func(_ bool) error { return nil }, func(err error) error { return err })
		require.Error(t, err)
		assert.True(t, errors.Is(err, expected))
		repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	t.Run("Family not found returns ErrFamilyNotFound", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		h := command.NewRevokeMemberCommandHandler(repo, authz)

		famID := types.NewFamilyID()
		cmd := command.RevokeMemberCommand{FamilyId: famID, FamilyMemberId: types.NewFamilyMemberID()}

		repo.EXPECT().GetById(mock.Anything, famID).Return(nil, nil)

		res := h.Handle(ctx, cmd)

		err := result.Match(res, func(_ bool) error { return nil }, func(err error) error { return err })
		require.Error(t, err)
		assert.True(t, errors.Is(err, family.ErrFamilyNotFound))
		repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	t.Run("Authorization denies returns error", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		permReq := ports.NewMockPermissionRequest(t)
		h := command.NewRevokeMemberCommandHandler(repo, authz)

		famID := types.NewFamilyID()
		fam := helperNewFamilyWithID(famID, types.UserID("owner-1"), "Doe family", nil)
		cmd := command.RevokeMemberCommand{FamilyId: famID, FamilyMemberId: types.NewFamilyMemberID()}

		repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
		authz.EXPECT().Can(mock.Anything, authorization.PermissionWrite).Return(permReq)
		expected := errors.New("forbidden")
		permReq.EXPECT().For(fam).Return(expected)

		res := h.Handle(ctx, cmd)

		err := result.Match(res, func(_ bool) error { return nil }, func(err error) error { return err })
		require.Error(t, err)
		assert.True(t, errors.Is(err, expected))
		repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	t.Run("Member not found returns ErrFamilyMemberNotFound", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		permReq := ports.NewMockPermissionRequest(t)
		h := command.NewRevokeMemberCommandHandler(repo, authz)

		famID := types.NewFamilyID()
		m := helperNewMember(famID, "John", family.AdultMemberType)
		fam := helperNewFamilyWithID(famID, types.UserID("owner-1"), "Doe family", []family.Member{m})
		cmd := command.RevokeMemberCommand{FamilyId: famID, FamilyMemberId: types.NewFamilyMemberID()}

		repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
		authz.EXPECT().Can(mock.Anything, authorization.PermissionWrite).Return(permReq)
		permReq.EXPECT().For(fam).Return(nil)

		res := h.Handle(ctx, cmd)

		err := result.Match(res, func(_ bool) error { return nil }, func(err error) error { return err })
		require.Error(t, err)
		assert.True(t, errors.Is(err, family.ErrFamilyMemberNotFound))
		repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	t.Run("Member has no userId -> success true and no Save", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		permReq := ports.NewMockPermissionRequest(t)
		h := command.NewRevokeMemberCommandHandler(repo, authz)

		famID := types.NewFamilyID()
		m := helperNewMember(famID, "John", family.AdultMemberType)
		fam := helperNewFamilyWithID(famID, types.UserID("owner-1"), "Doe family", []family.Member{m})
		cmd := command.RevokeMemberCommand{FamilyId: famID, FamilyMemberId: m.Id()}

		repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
		authz.EXPECT().Can(mock.Anything, authorization.PermissionWrite).Return(permReq)
		permReq.EXPECT().For(fam).Return(nil)

		res := h.Handle(ctx, cmd)

		require.True(t, res.IsSuccess())
		ok := result.Match(res, func(v bool) bool { return v }, func(err error) bool { return false })
		assert.True(t, ok)
		repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	t.Run("Validation error after clearing userId -> fault and no Save", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		permReq := ports.NewMockPermissionRequest(t)
		h := command.NewRevokeMemberCommandHandler(repo, authz)

		famID := types.NewFamilyID()
		m := helperNewMember(famID, "John", family.AdultMemberType)
		uid := types.UserID("user-1")
		m.SetUserId(&uid)
		fam := family.NewFamily(famID, types.UserID("owner-1"), "x", []family.Member{m}, time.Now(), time.Now())
		cmd := command.RevokeMemberCommand{FamilyId: famID, FamilyMemberId: m.Id()}

		repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
		authz.EXPECT().Can(mock.Anything, authorization.PermissionWrite).Return(permReq)
		permReq.EXPECT().For(fam).Return(nil)

		res := h.Handle(ctx, cmd)

		require.True(t, res.IsFaulted())
		repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	t.Run("Repository Save error propagates", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		permReq := ports.NewMockPermissionRequest(t)
		h := command.NewRevokeMemberCommandHandler(repo, authz)

		famID := types.NewFamilyID()
		m := helperNewMember(famID, "John", family.AdultMemberType)
		uid := types.UserID("user-1")
		m.SetUserId(&uid)
		fam := helperNewFamilyWithID(famID, types.UserID("owner-1"), "Doe family", []family.Member{m})
		cmd := command.RevokeMemberCommand{FamilyId: famID, FamilyMemberId: m.Id()}

		repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
		authz.EXPECT().Can(mock.Anything, authorization.PermissionWrite).Return(permReq)
		permReq.EXPECT().For(fam).Return(nil)
		repo.EXPECT().Save(mock.Anything, mock.Anything).Return(errors.New("save failed"))

		res := h.Handle(ctx, cmd)

		require.True(t, res.IsFaulted())
	})

	t.Run("Successful revoke clears userId and saves", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		permReq := ports.NewMockPermissionRequest(t)
		h := command.NewRevokeMemberCommandHandler(repo, authz)

		famID := types.NewFamilyID()
		m := helperNewMember(famID, "John", family.AdultMemberType)
		uid := types.UserID("user-1")
		m.SetUserId(&uid)
		fam := helperNewFamilyWithID(famID, types.UserID("owner-1"), "Doe family", []family.Member{m})
		cmd := command.RevokeMemberCommand{FamilyId: famID, FamilyMemberId: m.Id()}

		repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
		authz.EXPECT().Can(mock.Anything, authorization.PermissionWrite).Return(permReq)
		permReq.EXPECT().For(fam).Return(nil)
		repo.EXPECT().Save(mock.Anything, mock.Anything).Return(nil)

		res := h.Handle(ctx, cmd)

		require.True(t, res.IsSuccess())
		ok := result.Match(res, func(v bool) bool { return v }, func(err error) bool { return false })
		assert.True(t, ok)
		updated := fam.GetMember(m.Id())
		require.NotNil(t, updated)
		assert.Nil(t, updated.UserId())
	})
}
