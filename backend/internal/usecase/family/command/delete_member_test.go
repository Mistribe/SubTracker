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

func TestDeleteFamilyMemberCommandHandler(t *testing.T) {
	ctx := context.Background()

	t.Run("GetById returns error", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		h := command.NewDeleteFamilyMemberCommandHandler(repo, authz)

		famID := types.NewFamilyID()
		cmd := command.DeleteFamilyMemberCommand{FamilyID: famID, FamilyMemberID: types.NewFamilyMemberID()}

		expectedErr := errors.New("db error")
		repo.EXPECT().GetById(mock.Anything, famID).Return(nil, expectedErr)

		res := h.Handle(ctx, cmd)

		receivedErr := result.Match(res, func(_ bool) error { return nil }, func(err error) error { return err })
		require.Error(t, receivedErr)
		assert.True(t, errors.Is(receivedErr, expectedErr))
		repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	t.Run("Family not found returns ErrFamilyNotFound", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		h := command.NewDeleteFamilyMemberCommandHandler(repo, authz)

		famID := types.NewFamilyID()
		cmd := command.DeleteFamilyMemberCommand{FamilyID: famID, FamilyMemberID: types.NewFamilyMemberID()}

		repo.EXPECT().GetById(mock.Anything, famID).Return(nil, nil)

		res := h.Handle(ctx, cmd)

		receivedErr := result.Match(res, func(_ bool) error { return nil }, func(err error) error { return err })
		require.Error(t, receivedErr)
		assert.True(t, errors.Is(receivedErr, family.ErrFamilyNotFound))
		repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	t.Run("Authorization denies returns error", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		permReq := ports.NewMockPermissionRequest(t)
		h := command.NewDeleteFamilyMemberCommandHandler(repo, authz)

		famID := types.NewFamilyID()
		m := helperNewMember(famID, "John", family.AdultMemberType)
		fam := helperNewFamilyWithID(famID, types.UserID("owner-1"), "Doe family", []family.Member{m})
		cmd := command.DeleteFamilyMemberCommand{FamilyID: famID, FamilyMemberID: m.Id()}

		repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
		authz.EXPECT().Can(mock.Anything, authorization.PermissionDelete).Return(permReq)
		expected := errors.New("forbidden")
		permReq.EXPECT().For(fam).Return(expected)

		res := h.Handle(ctx, cmd)

		receivedErr := result.Match(res, func(_ bool) error { return nil }, func(err error) error { return err })
		require.Error(t, receivedErr)
		assert.True(t, errors.Is(receivedErr, expected))
		repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	t.Run("Member not found returns ErrFamilyMemberNotFound", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		permReq := ports.NewMockPermissionRequest(t)
		h := command.NewDeleteFamilyMemberCommandHandler(repo, authz)

		famID := types.NewFamilyID()
		other := helperNewMember(famID, "John", family.AdultMemberType)
		fam := helperNewFamilyWithID(famID, types.UserID("owner-1"), "Doe family", []family.Member{other})
		cmd := command.DeleteFamilyMemberCommand{FamilyID: famID, FamilyMemberID: types.NewFamilyMemberID()}

		repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
		authz.EXPECT().Can(mock.Anything, authorization.PermissionDelete).Return(permReq)
		permReq.EXPECT().For(fam).Return(nil)

		res := h.Handle(ctx, cmd)

		receivedErr := result.Match(res, func(_ bool) error { return nil }, func(err error) error { return err })
		require.Error(t, receivedErr)
		assert.True(t, errors.Is(receivedErr, family.ErrFamilyMemberNotFound))
		repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	t.Run("Validation error after removal -> fail and no save", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		permReq := ports.NewMockPermissionRequest(t)
		h := command.NewDeleteFamilyMemberCommandHandler(repo, authz)

		famID := types.NewFamilyID()
		m := helperNewMember(famID, "Solo", family.AdultMemberType)
		fam := helperNewFamilyWithID(famID, types.UserID("owner-1"), "Doe family", []family.Member{m})
		cmd := command.DeleteFamilyMemberCommand{FamilyID: famID, FamilyMemberID: m.Id()}

		repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
		authz.EXPECT().Can(mock.Anything, authorization.PermissionDelete).Return(permReq)
		permReq.EXPECT().For(fam).Return(nil)

		res := h.Handle(ctx, cmd)

		require.True(t, res.IsFaulted())
		repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	t.Run("Repository Save returns error", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		permReq := ports.NewMockPermissionRequest(t)
		h := command.NewDeleteFamilyMemberCommandHandler(repo, authz)

		famID := types.NewFamilyID()
		m1 := helperNewMember(famID, "John", family.AdultMemberType)
		m2 := helperNewMember(famID, "Jane", family.AdultMemberType)
		fam := helperNewFamilyWithID(famID, types.UserID("owner-1"), "Doe family", []family.Member{m1, m2})
		cmd := command.DeleteFamilyMemberCommand{FamilyID: famID, FamilyMemberID: m1.Id()}

		repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
		authz.EXPECT().Can(mock.Anything, authorization.PermissionDelete).Return(permReq)
		permReq.EXPECT().For(fam).Return(nil)
		repo.EXPECT().Save(mock.Anything, mock.Anything).Return(errors.New("save failed"))

		res := h.Handle(ctx, cmd)

		require.True(t, res.IsFaulted())
	})

	t.Run("Successful deletion returns success true and Save called", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		permReq := ports.NewMockPermissionRequest(t)
		h := command.NewDeleteFamilyMemberCommandHandler(repo, authz)

		famID := types.NewFamilyID()
		m1 := helperNewMember(famID, "John", family.AdultMemberType)
		m2 := helperNewMember(famID, "Jane", family.AdultMemberType)
		fam := helperNewFamilyWithID(famID, types.UserID("owner-1"), "Doe family", []family.Member{m1, m2})
		cmd := command.DeleteFamilyMemberCommand{FamilyID: famID, FamilyMemberID: m1.Id()}

		repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
		authz.EXPECT().Can(mock.Anything, authorization.PermissionDelete).Return(permReq)
		permReq.EXPECT().For(fam).Return(nil)
		repo.EXPECT().Save(mock.Anything, mock.Anything).Return(nil)

		res := h.Handle(ctx, cmd)

		require.True(t, res.IsSuccess())
		ok := result.Match(res, func(v bool) bool { return v }, func(err error) bool { return false })
		assert.True(t, ok)
		require.Nil(t, fam.GetMember(m1.Id()))
		require.NotNil(t, fam.GetMember(m2.Id()))
	})
}
