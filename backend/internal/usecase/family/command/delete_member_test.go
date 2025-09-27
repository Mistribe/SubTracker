package command_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/family/command"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

func makeMember(famID uuid.UUID, name string) family.Member {
	return family.NewMember(uuid.New(), famID, name, family.AdultMemberType, nil, time.Now(), time.Now())
}

func makeFamilyWithMembers(famID uuid.UUID, members []family.Member) family.Family {
	return family.NewFamily(famID, "owner-1", "Doe family", members, time.Now(), time.Now())
}

func TestDeleteFamilyMemberCommandHandler(t *testing.T) {
	ctx := context.Background()

	t.Run("GetById returns error", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		auth := stubAuthorization{}
		h := command.NewDeleteFamilyMemberCommandHandler(repo, auth)

		famID := uuid.New()
		cmd := command.DeleteFamilyMemberCommand{FamilyID: famID, FamilyMemberID: uuid.New()}

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
		auth := stubAuthorization{}
		h := command.NewDeleteFamilyMemberCommandHandler(repo, auth)

		famID := uuid.New()
		cmd := command.DeleteFamilyMemberCommand{FamilyID: famID, FamilyMemberID: uuid.New()}

		repo.EXPECT().GetById(mock.Anything, famID).Return(nil, nil)

		res := h.Handle(ctx, cmd)

		receivedErr := result.Match(res, func(_ bool) error { return nil }, func(err error) error { return err })
		require.Error(t, receivedErr)
		assert.True(t, errors.Is(receivedErr, family.ErrFamilyNotFound))
		repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	t.Run("Authorization denies returns error", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		permReq := ports.NewMockPermissionRequest(t)
		auth := stubAuthorization{req: permReq}
		h := command.NewDeleteFamilyMemberCommandHandler(repo, auth)

		famID := uuid.New()
		m := makeMember(famID, "John")
		fam := makeFamilyWithMembers(famID, []family.Member{m})
		cmd := command.DeleteFamilyMemberCommand{FamilyID: famID, FamilyMemberID: m.Id()}

		expectedErr := errors.New("forbidden")
		repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
		permReq.EXPECT().For(fam).Return(expectedErr)

		res := h.Handle(ctx, cmd)

		receivedErr := result.Match(res, func(_ bool) error { return nil }, func(err error) error { return err })
		require.Error(t, receivedErr)
		assert.True(t, errors.Is(receivedErr, expectedErr))
		repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	t.Run("Member not found returns ErrFamilyMemberNotFound", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		permReq := ports.NewMockPermissionRequest(t)
		auth := stubAuthorization{req: permReq}
		h := command.NewDeleteFamilyMemberCommandHandler(repo, auth)

		famID := uuid.New()
		// family with a different member
		other := makeMember(famID, "John")
		fam := makeFamilyWithMembers(famID, []family.Member{other})
		cmd := command.DeleteFamilyMemberCommand{FamilyID: famID, FamilyMemberID: uuid.New()}

		repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
		permReq.EXPECT().For(fam).Return(nil)

		res := h.Handle(ctx, cmd)

		receivedErr := result.Match(res, func(_ bool) error { return nil }, func(err error) error { return err })
		require.Error(t, receivedErr)
		assert.True(t, errors.Is(receivedErr, family.ErrFamilyMemberNotFound))
		repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	t.Run("Validation error after removal -> fail and no save", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		permReq := ports.NewMockPermissionRequest(t)
		auth := stubAuthorization{req: permReq}
		h := command.NewDeleteFamilyMemberCommandHandler(repo, auth)

		famID := uuid.New()
		// single-member family so removing it breaks validation (family must have at least one member)
		m := makeMember(famID, "Solo")
		fam := makeFamilyWithMembers(famID, []family.Member{m})
		cmd := command.DeleteFamilyMemberCommand{FamilyID: famID, FamilyMemberID: m.Id()}

		repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
		permReq.EXPECT().For(fam).Return(nil)

		res := h.Handle(ctx, cmd)

		// We only assert fault and that Save was not called
		require.True(t, res.IsFaulted())
		repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	t.Run("Repository Save returns error", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		permReq := ports.NewMockPermissionRequest(t)
		auth := stubAuthorization{req: permReq}
		h := command.NewDeleteFamilyMemberCommandHandler(repo, auth)

		famID := uuid.New()
		// two-member family so validation passes after removal
		m1 := makeMember(famID, "John")
		m2 := makeMember(famID, "Jane")
		fam := makeFamilyWithMembers(famID, []family.Member{m1, m2})
		cmd := command.DeleteFamilyMemberCommand{FamilyID: famID, FamilyMemberID: m1.Id()}

		repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
		permReq.EXPECT().For(fam).Return(nil)
		repo.EXPECT().Save(mock.Anything, mock.Anything).Return(errors.New("save failed"))

		res := h.Handle(ctx, cmd)

		require.True(t, res.IsFaulted())
	})

	t.Run("Successful deletion returns success true and Save called", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		permReq := ports.NewMockPermissionRequest(t)
		auth := stubAuthorization{req: permReq}
		h := command.NewDeleteFamilyMemberCommandHandler(repo, auth)

		famID := uuid.New()
		// two-member family so validation passes after removal
		m1 := makeMember(famID, "John")
		m2 := makeMember(famID, "Jane")
		fam := makeFamilyWithMembers(famID, []family.Member{m1, m2})
		cmd := command.DeleteFamilyMemberCommand{FamilyID: famID, FamilyMemberID: m1.Id()}

		repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
		permReq.EXPECT().For(fam).Return(nil)
		repo.EXPECT().Save(mock.Anything, mock.Anything).Return(nil)

		res := h.Handle(ctx, cmd)

		require.True(t, res.IsSuccess())
		ok := result.Match(res, func(v bool) bool { return v }, func(err error) bool { return false })
		assert.True(t, ok)
		// ensure the removed member is gone
		require.Nil(t, fam.GetMember(m1.Id()))
		require.NotNil(t, fam.GetMember(m2.Id()))
	})
}
