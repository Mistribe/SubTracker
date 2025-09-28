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
	"github.com/mistribe/subtracker/internal/domain/billing"
	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/family/command"
	"github.com/mistribe/subtracker/pkg/langext/option"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

func TestCreateFamilyMemberCommandHandler(t *testing.T) {
	ctx := context.Background()

	t.Run("GetById returns error", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		ent := ports.NewMockEntitlementResolver(t)
		h := command.NewCreateFamilyMemberCommandHandler(repo, authz, ent)

		famID := types.NewFamilyID()
		cmd := command.CreateFamilyMemberCommand{FamilyID: famID, Name: "John", Type: family.AdultMemberType}

		expected := errors.New("db error")
		repo.EXPECT().GetById(mock.Anything, famID).Return(nil, expected)

		res := h.Handle(ctx, cmd)
		err := result.Match(res, func(_ family.Family) error { return nil }, func(e error) error { return e })
		require.Error(t, err)
		assert.True(t, errors.Is(err, expected))
		repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	t.Run("Family not found returns ErrFamilyNotFound", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		ent := ports.NewMockEntitlementResolver(t)
		h := command.NewCreateFamilyMemberCommandHandler(repo, authz, ent)

		famID := types.NewFamilyID()
		cmd := command.CreateFamilyMemberCommand{FamilyID: famID, Name: "Jane", Type: family.AdultMemberType}

		repo.EXPECT().GetById(mock.Anything, famID).Return(nil, nil)

		res := h.Handle(ctx, cmd)
		err := result.Match(res, func(_ family.Family) error { return nil }, func(e error) error { return e })
		require.Error(t, err)
		assert.True(t, errors.Is(err, family.ErrFamilyNotFound))
		repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	t.Run("Authorization Can.For returns error", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)
		ent := ports.NewMockEntitlementResolver(t)
		h := command.NewCreateFamilyMemberCommandHandler(repo, authz, ent)

		famID := types.NewFamilyID()
		m := helperNewMember(famID, "John", family.AdultMemberType)
		existing := helperNewFamilyWithID(famID, types.UserID("owner1"), "Doe family", []family.Member{m})
		cmd := command.CreateFamilyMemberCommand{FamilyID: existing.Id(), Name: "Jane", Type: family.AdultMemberType}

		repo.EXPECT().GetById(mock.Anything, existing.Id()).Return(existing, nil)
		authz.EXPECT().Can(mock.Anything, authorization.PermissionWrite).Return(perm)
		expected := errors.New("unauthorized")
		perm.EXPECT().For(existing).Return(expected)

		res := h.Handle(ctx, cmd)
		err := result.Match(res, func(_ family.Family) error { return nil }, func(e error) error { return e })
		require.Error(t, err)
		assert.True(t, errors.Is(err, expected))
		repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
		ent.AssertNotCalled(t, "CheckQuota", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("Entitlement resolver returns error", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)
		ent := ports.NewMockEntitlementResolver(t)
		h := command.NewCreateFamilyMemberCommandHandler(repo, authz, ent)

		famID := types.NewFamilyID()
		m := helperNewMember(famID, "John", family.AdultMemberType)
		existing := helperNewFamilyWithID(famID, types.UserID("owner1"), "Doe family", []family.Member{m})
		cmd := command.CreateFamilyMemberCommand{FamilyID: existing.Id(), Name: "Jane", Type: family.AdultMemberType}

		repo.EXPECT().GetById(mock.Anything, existing.Id()).Return(existing, nil)
		authz.EXPECT().Can(mock.Anything, authorization.PermissionWrite).Return(perm)
		perm.EXPECT().For(existing).Return(nil)
		expected := errors.New("quota error")
		ent.EXPECT().CheckQuota(mock.Anything, billing.FeatureIdFamilyMembersCount, int64(1)).Return(false,
			billing.EffectiveEntitlement{}, expected)

		res := h.Handle(ctx, cmd)
		err := result.Match(res, func(_ family.Family) error { return nil }, func(e error) error { return e })
		require.Error(t, err)
		assert.True(t, errors.Is(err, expected))
		repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	t.Run("Entitlement denies (limit reached) returns ErrFamilyMembersLimitReached", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)
		ent := ports.NewMockEntitlementResolver(t)
		h := command.NewCreateFamilyMemberCommandHandler(repo, authz, ent)

		famID := types.NewFamilyID()
		m := helperNewMember(famID, "John", family.AdultMemberType)
		existing := helperNewFamilyWithID(famID, types.UserID("owner1"), "Doe family", []family.Member{m})
		cmd := command.CreateFamilyMemberCommand{FamilyID: existing.Id(), Name: "Jane", Type: family.AdultMemberType}

		repo.EXPECT().GetById(mock.Anything, existing.Id()).Return(existing, nil)
		authz.EXPECT().Can(mock.Anything, authorization.PermissionWrite).Return(perm)
		perm.EXPECT().For(existing).Return(nil)
		ent.EXPECT().CheckQuota(mock.Anything, billing.FeatureIdFamilyMembersCount, int64(1)).Return(false,
			billing.EffectiveEntitlement{}, nil)

		res := h.Handle(ctx, cmd)
		err := result.Match(res, func(_ family.Family) error { return nil }, func(e error) error { return e })
		require.Error(t, err)
		assert.True(t, errors.Is(err, family.ErrFamilyMembersLimitReached))
		repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	t.Run("AddMember duplicate id returns success without save", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)
		ent := ports.NewMockEntitlementResolver(t)
		h := command.NewCreateFamilyMemberCommandHandler(repo, authz, ent)

		existingMemberID := types.NewFamilyMemberID()
		famID := types.NewFamilyID()
		existingMember := helperNewMemberWithID(famID, existingMemberID, "John", family.AdultMemberType)
		fam := helperNewFamilyWithID(famID, types.UserID("owner1"), "Doe family", []family.Member{existingMember})
		cmd := command.CreateFamilyMemberCommand{
			FamilyID: fam.Id(), FamilyMemberID: option.Some(existingMemberID), Name: "John",
			Type: family.AdultMemberType,
		}

		repo.EXPECT().GetById(mock.Anything, fam.Id()).Return(fam, nil)
		authz.EXPECT().Can(mock.Anything, authorization.PermissionWrite).Return(perm)
		perm.EXPECT().For(fam).Return(nil)
		ent.EXPECT().CheckQuota(mock.Anything, billing.FeatureIdFamilyMembersCount, int64(1)).Return(true,
			billing.EffectiveEntitlement{}, nil)

		before := fam.Members().Len()
		res := h.Handle(ctx, cmd)
		require.True(t, res.IsSuccess())
		assert.Equal(t, before, fam.Members().Len())
		repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	t.Run("Validation error after AddMember => fail and no save", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)
		ent := ports.NewMockEntitlementResolver(t)
		h := command.NewCreateFamilyMemberCommandHandler(repo, authz, ent)

		famID := types.NewFamilyID()
		m := helperNewMember(famID, "John", family.AdultMemberType)
		fam := helperNewFamilyWithID(famID, types.UserID("owner1"), "x", []family.Member{m})
		cmd := command.CreateFamilyMemberCommand{FamilyID: fam.Id(), Name: "Jane", Type: family.AdultMemberType}

		repo.EXPECT().GetById(mock.Anything, fam.Id()).Return(fam, nil)
		authz.EXPECT().Can(mock.Anything, authorization.PermissionWrite).Return(perm)
		perm.EXPECT().For(fam).Return(nil)
		ent.EXPECT().CheckQuota(mock.Anything, billing.FeatureIdFamilyMembersCount, int64(1)).Return(true,
			billing.EffectiveEntitlement{}, nil)

		res := h.Handle(ctx, cmd)
		require.True(t, res.IsFaulted())
		repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	t.Run("Repository Save returns error", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)
		ent := ports.NewMockEntitlementResolver(t)
		h := command.NewCreateFamilyMemberCommandHandler(repo, authz, ent)

		famID := types.NewFamilyID()
		m := helperNewMember(famID, "John", family.AdultMemberType)
		fam := helperNewFamilyWithID(famID, types.UserID("owner1"), "Doe family", []family.Member{m})
		cmd := command.CreateFamilyMemberCommand{FamilyID: fam.Id(), Name: "Jane", Type: family.AdultMemberType}

		repo.EXPECT().GetById(mock.Anything, fam.Id()).Return(fam, nil)
		authz.EXPECT().Can(mock.Anything, authorization.PermissionWrite).Return(perm)
		perm.EXPECT().For(fam).Return(nil)
		ent.EXPECT().CheckQuota(mock.Anything, billing.FeatureIdFamilyMembersCount, int64(1)).Return(true,
			billing.EffectiveEntitlement{}, nil)
		repo.EXPECT().Save(mock.Anything, mock.Anything).Return(errors.New("save failed"))

		res := h.Handle(ctx, cmd)
		require.True(t, res.IsFaulted())
	})

	t.Run("Successful add -> success and save called", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)
		ent := ports.NewMockEntitlementResolver(t)
		h := command.NewCreateFamilyMemberCommandHandler(repo, authz, ent)

		famID := types.NewFamilyID()
		m := helperNewMember(famID, "John", family.AdultMemberType)
		fam := helperNewFamilyWithID(famID, types.UserID("owner1"), "Doe family", []family.Member{m})
		cmd := command.CreateFamilyMemberCommand{
			FamilyID: fam.Id(), Name: "Jane", Type: family.AdultMemberType,
			CreatedAt: option.Some(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)),
		}

		repo.EXPECT().GetById(mock.Anything, fam.Id()).Return(fam, nil)
		authz.EXPECT().Can(mock.Anything, authorization.PermissionWrite).Return(perm)
		perm.EXPECT().For(fam).Return(nil)
		ent.EXPECT().CheckQuota(mock.Anything, billing.FeatureIdFamilyMembersCount, int64(1)).Return(true,
			billing.EffectiveEntitlement{}, nil)
		repo.EXPECT().Save(mock.Anything, mock.Anything).Return(nil)

		res := h.Handle(ctx, cmd)
		require.True(t, res.IsSuccess())
	})
}
