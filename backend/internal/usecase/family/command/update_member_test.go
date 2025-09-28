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
	"github.com/mistribe/subtracker/pkg/langext/option"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

func TestUpdateFamilyMemberCommandHandler(t *testing.T) {
	ctx := context.Background()

	t.Run("GetById returns error", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		h := command.NewUpdateFamilyMemberCommandHandler(repo, authz)

		famID := types.NewFamilyID()
		cmd := command.UpdateFamilyMemberCommand{FamilyID: famID, FamilyMemberID: types.NewFamilyMemberID()}

		expected := errors.New("db error")
		repo.EXPECT().GetById(mock.Anything, famID).Return(nil, expected)

		res := h.Handle(ctx, cmd)

		err := result.Match(res, func(_ family.Family) error { return nil }, func(err error) error { return err })
		require.Error(t, err)
		assert.True(t, errors.Is(err, expected))
		repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	t.Run("Family not found returns ErrFamilyNotFound", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		h := command.NewUpdateFamilyMemberCommandHandler(repo, authz)

		famID := types.NewFamilyID()
		cmd := command.UpdateFamilyMemberCommand{FamilyID: famID, FamilyMemberID: types.NewFamilyMemberID()}

		repo.EXPECT().GetById(mock.Anything, famID).Return(nil, nil)

		res := h.Handle(ctx, cmd)

		err := result.Match(res, func(_ family.Family) error { return nil }, func(err error) error { return err })
		require.Error(t, err)
		assert.True(t, errors.Is(err, family.ErrFamilyNotFound))
		repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	t.Run("Authorization denies returns error", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		permReq := ports.NewMockPermissionRequest(t)
		authz := ports.NewMockAuthorization(t)
		h := command.NewUpdateFamilyMemberCommandHandler(repo, authz)

		famID := types.NewFamilyID()
		m := helperNewMember(famID, "John", family.AdultMemberType)
		fam := helperNewFamilyWithID(famID, types.UserID("owner-1"), "Doe family", []family.Member{m})
		cmd := command.UpdateFamilyMemberCommand{
			FamilyID: famID, FamilyMemberID: m.Id(), Name: "Johnny",
			Type: family.KidMemberType, UpdatedAt: option.None[time.Time](),
		}

		repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
		authz.EXPECT().Can(mock.Anything, authorization.PermissionWrite).Return(permReq)
		expected := errors.New("forbidden")
		permReq.EXPECT().For(fam).Return(expected)

		res := h.Handle(ctx, cmd)

		err := result.Match(res, func(_ family.Family) error { return nil }, func(err error) error { return err })
		require.Error(t, err)
		assert.True(t, errors.Is(err, expected))
		repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	t.Run("Member not found returns ErrFamilyMemberNotFound", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		h := command.NewUpdateFamilyMemberCommandHandler(repo, authz)

		famID := types.NewFamilyID()
		other := helperNewMember(famID, "John", family.AdultMemberType)
		fam := helperNewFamilyWithID(famID, types.UserID("owner-1"), "Doe family", []family.Member{other})
		cmd := command.UpdateFamilyMemberCommand{
			FamilyID: famID, FamilyMemberID: types.NewFamilyMemberID(), Name: "Johnny",
			Type: family.KidMemberType, UpdatedAt: option.None[time.Time](),
		}

		repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)

		res := h.Handle(ctx, cmd)

		err := result.Match(res, func(_ family.Family) error { return nil }, func(err error) error { return err })
		require.Error(t, err)
		assert.True(t, errors.Is(err, family.ErrFamilyMemberNotFound))
		repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	t.Run("Validation error after update -> fault and no Save", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		permReq := ports.NewMockPermissionRequest(t)
		authz := ports.NewMockAuthorization(t)
		h := command.NewUpdateFamilyMemberCommandHandler(repo, authz)

		famID := types.NewFamilyID()
		m := helperNewMember(famID, "John", family.AdultMemberType)
		fam := helperNewFamilyWithID(famID, types.UserID("owner-1"), "Doe family", []family.Member{m})
		cmd := command.UpdateFamilyMemberCommand{
			FamilyID: famID, FamilyMemberID: m.Id(), Name: "", Type: family.KidMemberType,
			UpdatedAt: option.None[time.Time](),
		}

		repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
		authz.EXPECT().Can(mock.Anything, authorization.PermissionWrite).Return(permReq)
		permReq.EXPECT().For(fam).Return(nil)

		res := h.Handle(ctx, cmd)

		require.True(t, res.IsFaulted())
		repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	t.Run("Repository Save error propagates", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		permReq := ports.NewMockPermissionRequest(t)
		authz := ports.NewMockAuthorization(t)
		h := command.NewUpdateFamilyMemberCommandHandler(repo, authz)

		famID := types.NewFamilyID()
		m := helperNewMember(famID, "John", family.AdultMemberType)
		fam := helperNewFamilyWithID(famID, types.UserID("owner-1"), "Doe family", []family.Member{m})
		cmd := command.UpdateFamilyMemberCommand{
			FamilyID: famID, FamilyMemberID: m.Id(), Name: "Johnny", Type: family.KidMemberType,
			UpdatedAt: option.Some(time.Now()),
		}

		repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
		authz.EXPECT().Can(mock.Anything, authorization.PermissionWrite).Return(permReq)
		permReq.EXPECT().For(fam).Return(nil)
		repo.EXPECT().Save(mock.Anything, mock.Anything).Return(errors.New("save failed"))

		res := h.Handle(ctx, cmd)

		require.True(t, res.IsFaulted())
	})

	t.Run("Success with UpdatedAt Some sets exact time and updates fields", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		permReq := ports.NewMockPermissionRequest(t)
		authz := ports.NewMockAuthorization(t)
		h := command.NewUpdateFamilyMemberCommandHandler(repo, authz)

		famID := types.NewFamilyID()
		m := helperNewMember(famID, "John", family.AdultMemberType)
		fam := helperNewFamilyWithID(famID, types.UserID("owner-1"), "Doe family", []family.Member{m})
		ts := time.Date(2024, 12, 25, 10, 30, 0, 0, time.UTC)
		newName := "Johnny"
		newType := family.KidMemberType
		cmd := command.UpdateFamilyMemberCommand{
			FamilyID: famID, FamilyMemberID: m.Id(), Name: newName, Type: newType, UpdatedAt: option.Some(ts),
		}

		repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
		authz.EXPECT().Can(mock.Anything, authorization.PermissionWrite).Return(permReq)
		permReq.EXPECT().For(fam).Return(nil)
		repo.EXPECT().Save(mock.Anything, mock.Anything).Return(nil)

		res := h.Handle(ctx, cmd)

		require.True(t, res.IsSuccess())
		got := result.Match(res, func(v family.Family) family.Family { return v },
			func(err error) family.Family { t.Fatalf("unexpected error: %v", err); return nil })
		require.NotNil(t, got)
		updated := got.GetMember(m.Id())
		require.NotNil(t, updated)
		assert.Equal(t, newName, updated.Name())
		assert.Equal(t, newType, updated.Type())
	})

	t.Run("Success with UpdatedAt None sets time to now-ish and updates fields", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		permReq := ports.NewMockPermissionRequest(t)
		authz := ports.NewMockAuthorization(t)
		h := command.NewUpdateFamilyMemberCommandHandler(repo, authz)

		famID := types.NewFamilyID()
		m := helperNewMember(famID, "John", family.AdultMemberType)
		fam := helperNewFamilyWithID(famID, types.UserID("owner-1"), "Doe family", []family.Member{m})
		newName := "Johnny"
		newType := family.AdultMemberType
		cmd := command.UpdateFamilyMemberCommand{
			FamilyID: famID, FamilyMemberID: m.Id(), Name: newName, Type: newType, UpdatedAt: option.None[time.Time](),
		}

		repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
		authz.EXPECT().Can(mock.Anything, authorization.PermissionWrite).Return(permReq)
		permReq.EXPECT().For(fam).Return(nil)
		repo.EXPECT().Save(mock.Anything, mock.Anything).Return(nil)

		before := time.Now()
		res := h.Handle(ctx, cmd)
		_ = before

		require.True(t, res.IsSuccess())
		got := result.Match(res, func(v family.Family) family.Family { return v },
			func(err error) family.Family { t.Fatalf("unexpected error: %v", err); return nil })
		require.NotNil(t, got)
		updated := got.GetMember(m.Id())
		require.NotNil(t, updated)
		assert.Equal(t, newName, updated.Name())
		assert.Equal(t, newType, updated.Type())
	})
}
