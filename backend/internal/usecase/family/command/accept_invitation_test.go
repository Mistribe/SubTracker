package command_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/mistribe/subtracker/internal/domain/account"
	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/family/command"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

func TestAcceptInvitationCommandHandler(t *testing.T) {
	ctx := context.Background()

	t.Run("GetById returns error", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		auth := ports.NewMockAuthentication(t)
		h := command.NewAcceptInvitationCommandHandler(repo, auth)

		famID := types.NewFamilyID()
		cmd := command.AcceptInvitationCommand{FamilyId: famID}

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
		auth := ports.NewMockAuthentication(t)
		h := command.NewAcceptInvitationCommandHandler(repo, auth)

		famID := types.NewFamilyID()
		cmd := command.AcceptInvitationCommand{FamilyId: famID}

		repo.EXPECT().GetById(mock.Anything, famID).Return(nil, nil)

		res := h.Handle(ctx, cmd)

		receivedErr := result.Match(res, func(_ bool) error { return nil }, func(err error) error { return err })
		require.Error(t, receivedErr)
		assert.True(t, errors.Is(receivedErr, family.ErrFamilyNotFound))
		repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	t.Run("Member not found returns ErrFamilyMemberNotFound", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		auth := ports.NewMockAuthentication(t)
		h := command.NewAcceptInvitationCommandHandler(repo, auth)

		famID := types.NewFamilyID()
		other := helperNewMember(famID, "John", family.AdultMemberType)
		fam := helperNewFamilyWithID(famID, types.UserID("owner1"), "Doe family", []family.Member{other})
		cmd := command.AcceptInvitationCommand{FamilyId: famID, FamilyMemberId: types.NewFamilyMemberID()}

		repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)

		res := h.Handle(ctx, cmd)

		receivedErr := result.Match(res, func(_ bool) error { return nil }, func(err error) error { return err })
		require.Error(t, receivedErr)
		assert.True(t, errors.Is(receivedErr, family.ErrFamilyMemberNotFound))
		repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	t.Run("Invitation code is nil returns ErrBadInvitationCode", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		auth := ports.NewMockAuthentication(t)
		h := command.NewAcceptInvitationCommandHandler(repo, auth)

		famID := types.NewFamilyID()
		m := helperNewMember(famID, "John", family.AdultMemberType) // no invitation code set
		fam := helperNewFamilyWithID(famID, types.UserID("owner1"), "Doe family", []family.Member{m})
		cmd := command.AcceptInvitationCommand{FamilyId: famID, FamilyMemberId: m.Id(), InvitationCode: "expected"}

		repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)

		res := h.Handle(ctx, cmd)

		receivedErr := result.Match(res, func(_ bool) error { return nil }, func(err error) error { return err })
		require.Error(t, receivedErr)
		assert.True(t, errors.Is(receivedErr, family.ErrBadInvitationCode))
		repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	t.Run("Invitation code mismatch returns ErrBadInvitationCode", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		auth := ports.NewMockAuthentication(t)
		h := command.NewAcceptInvitationCommandHandler(repo, auth)

		famID := types.NewFamilyID()
		code := "expected"
		m := helperNewMember(famID, "John", family.AdultMemberType)
		m.SetInvitationCode(&code)
		fam := helperNewFamilyWithID(famID, types.UserID("owner1"), "Doe family", []family.Member{m})
		cmd := command.AcceptInvitationCommand{FamilyId: famID, FamilyMemberId: m.Id(), InvitationCode: "wrong"}

		repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)

		res := h.Handle(ctx, cmd)

		receivedErr := result.Match(res, func(_ bool) error { return nil }, func(err error) error { return err })
		require.Error(t, receivedErr)
		assert.True(t, errors.Is(receivedErr, family.ErrBadInvitationCode))
		repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	t.Run("Validation error after update -> fail and no save", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		auth := ports.NewMockAuthentication(t)
		ca := account.NewMockConnectedAccount(t)
		h := command.NewAcceptInvitationCommandHandler(repo, auth)

		famID := types.NewFamilyID()
		code := "expected"
		m := helperNewMember(famID, "John", family.AdultMemberType)
		m.SetInvitationCode(&code)
		// invalid family name (<3)
		fam := helperNewFamilyWithID(famID, types.UserID("owner1"), "x", []family.Member{m})
		cmd := command.AcceptInvitationCommand{FamilyId: famID, FamilyMemberId: m.Id(), InvitationCode: code}

		repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
		ca.On("UserID").Return(types.UserID("user-123"))
		auth.EXPECT().MustGetConnectedAccount(mock.Anything).Return(ca)

		res := h.Handle(ctx, cmd)
		require.True(t, res.IsFaulted())
		repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
		ca.AssertExpectations(t)
	})

	t.Run("Repository Save returns error", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		auth := ports.NewMockAuthentication(t)
		ca := account.NewMockConnectedAccount(t)
		h := command.NewAcceptInvitationCommandHandler(repo, auth)

		famID := types.NewFamilyID()
		code := "expected"
		m := helperNewMember(famID, "John", family.AdultMemberType)
		m.SetInvitationCode(&code)
		fam := helperNewFamilyWithID(famID, types.UserID("owner1"), "Doe family", []family.Member{m})
		cmd := command.AcceptInvitationCommand{FamilyId: famID, FamilyMemberId: m.Id(), InvitationCode: code}

		repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
		ca.On("UserID").Return(types.UserID("user-123"))
		auth.EXPECT().MustGetConnectedAccount(mock.Anything).Return(ca)
		repo.EXPECT().Save(mock.Anything, mock.Anything).Return(errors.New("save failed"))

		res := h.Handle(ctx, cmd)
		require.True(t, res.IsFaulted())
		ca.AssertExpectations(t)
	})

	t.Run("Successful acceptance -> success true, userId set, save called", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		auth := ports.NewMockAuthentication(t)
		ca := account.NewMockConnectedAccount(t)
		h := command.NewAcceptInvitationCommandHandler(repo, auth)

		famID := types.NewFamilyID()
		code := "expected"
		m := helperNewMember(famID, "John", family.AdultMemberType)
		m.SetInvitationCode(&code)
		fam := helperNewFamilyWithID(famID, types.UserID("owner1"), "Doe family", []family.Member{m})
		cmd := command.AcceptInvitationCommand{FamilyId: famID, FamilyMemberId: m.Id(), InvitationCode: code}

		repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
		ca.On("UserID").Return(types.UserID("user-123"))
		auth.EXPECT().MustGetConnectedAccount(mock.Anything).Return(ca)
		repo.EXPECT().Save(mock.Anything, mock.Anything).Return(nil)

		res := h.Handle(ctx, cmd)

		require.True(t, res.IsSuccess())
		updated := fam.GetMember(m.Id())
		require.NotNil(t, updated)
		require.NotNil(t, updated.UserId())
		assert.Equal(t, types.UserID("user-123"), *updated.UserId())
		ca.AssertExpectations(t)
	})
}
