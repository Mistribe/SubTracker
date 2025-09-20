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
	"github.com/mistribe/subtracker/internal/domain/user"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/shared"
	"github.com/mistribe/subtracker/internal/usecase/family/command"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type permissionRequestStub struct { // implements ports.PermissionRequest
	err    error
	called bool
}

func (p *permissionRequestStub) For(entity ports.EntityWithOwnership) error {
	p.called = true
	return p.err
}

type authStub struct { // implements ports.Authorization
	perm         ports.PermissionRequest
	ensureErr    error
	ensureCalled bool
}

func (a *authStub) EnsureWithinLimit(ctx context.Context, feature user.Feature, delta int) error {
	a.ensureCalled = true
	return a.ensureErr
}

func (a *authStub) Can(ctx context.Context, permission user.Permission) ports.PermissionRequest {
	return a.perm
}

func (a *authStub) GetCurrentLimit(ctx context.Context, feature user.Feature) (shared.Limit, error) {
	return nil, nil
}

func (a *authStub) GetCurrentLimits(ctx context.Context, category user.Category) (shared.Limits, error) {
	return nil, nil
}

func newMember(familyID uuid.UUID, name string) family.Member {
	return family.NewMember(
		uuid.New(),
		familyID,
		name,
		family.AdultMemberType,
		nil,
		time.Now(),
		time.Now(),
	)
}

func newFamily(ownerID string, name string, members []family.Member) family.Family {
	return family.NewFamily(
		uuid.New(),
		ownerID,
		name,
		members,
		time.Now(),
		time.Now(),
	)
}

func TestCreateFamilyMemberCommandHandler(t *testing.T) {
	ctx := context.Background()

	t.Run("GetById returns error", func(t *testing.T) {
		// Arrange
		repo := ports.NewMockFamilyRepository(t)
		auth := &authStub{}
		h := command.NewCreateFamilyMemberCommandHandler(repo, auth)

		familyID := uuid.New()
		cmd := command.CreateFamilyMemberCommand{FamilyId: familyID, Member: nil}

		expectedErr := errors.New("db error")
		repo.EXPECT().GetById(mock.Anything, familyID).Return(nil, expectedErr)

		// Act
		res := h.Handle(ctx, cmd)

		// Assert
		receivedErr := result.Match(res, func(_ family.Family) error { return nil }, func(err error) error { return err })
		require.Error(t, receivedErr)
		assert.True(t, errors.Is(receivedErr, expectedErr))
		repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	t.Run("Family not found returns ErrFamilyNotFound", func(t *testing.T) {
		// Arrange
		repo := ports.NewMockFamilyRepository(t)
		auth := &authStub{}
		h := command.NewCreateFamilyMemberCommandHandler(repo, auth)

		familyID := uuid.New()
		cmd := command.CreateFamilyMemberCommand{FamilyId: familyID, Member: nil}

		repo.EXPECT().GetById(mock.Anything, familyID).Return(nil, nil)

		// Act
		res := h.Handle(ctx, cmd)

		// Assert
		receivedErr := result.Match(res, func(_ family.Family) error { return nil }, func(err error) error { return err })
		require.Error(t, receivedErr)
		assert.True(t, errors.Is(receivedErr, family.ErrFamilyNotFound))
		repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	t.Run("Authorization Can.For returns error", func(t *testing.T) {
		// Arrange
		repo := ports.NewMockFamilyRepository(t)
		auth := &authStub{}
		h := command.NewCreateFamilyMemberCommandHandler(repo, auth)

		fam := newFamily("owner1", "Doe family", []family.Member{newMember(uuid.New(), "John")})
		cmd := command.CreateFamilyMemberCommand{FamilyId: fam.Id(), Member: newMember(fam.Id(), "Jane")}

		repo.EXPECT().GetById(mock.Anything, fam.Id()).Return(fam, nil)
		expectedErr := errors.New("unauthorized")
		auth.perm = &permissionRequestStub{err: expectedErr}

		// Act
		res := h.Handle(ctx, cmd)

		// Assert
		receivedErr := result.Match(res, func(_ family.Family) error { return nil }, func(err error) error { return err })
		require.Error(t, receivedErr)
		assert.True(t, errors.Is(receivedErr, expectedErr))
		assert.False(t, auth.ensureCalled, "EnsureWithinLimit should not be called when authorization fails")
		repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	t.Run("EnsureWithinLimit returns error", func(t *testing.T) {
		// Arrange
		repo := ports.NewMockFamilyRepository(t)
		auth := &authStub{}
		h := command.NewCreateFamilyMemberCommandHandler(repo, auth)

		fam := newFamily("owner1", "Doe family", []family.Member{newMember(uuid.New(), "John")})
		cmd := command.CreateFamilyMemberCommand{FamilyId: fam.Id(), Member: newMember(fam.Id(), "Jane")}

		repo.EXPECT().GetById(mock.Anything, fam.Id()).Return(fam, nil)
		expectedErr := errors.New("limit exceeded")
		auth.perm = &permissionRequestStub{err: nil}
		auth.ensureErr = expectedErr

		// Act
		res := h.Handle(ctx, cmd)

		// Assert
		receivedErr := result.Match(res, func(_ family.Family) error { return nil }, func(err error) error { return err })
		require.Error(t, receivedErr)
		assert.True(t, errors.Is(receivedErr, expectedErr))
		assert.True(t, auth.ensureCalled, "EnsureWithinLimit should be called")
		repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	t.Run("AddMember returns error => success without save", func(t *testing.T) {
		// Arrange
		repo := ports.NewMockFamilyRepository(t)
		auth := &authStub{perm: &permissionRequestStub{err: nil}}
		h := command.NewCreateFamilyMemberCommandHandler(repo, auth)

		famID := uuid.New()
		existing := family.NewMember(uuid.New(), famID, "John", family.AdultMemberType, nil, time.Now(), time.Now())
		// Use same ID to force duplicate
		dup := family.NewMember(existing.Id(), famID, "John", family.AdultMemberType, nil, time.Now(), time.Now())
		fam := newFamily("owner1", "Doe family", []family.Member{existing})
		cmd := command.CreateFamilyMemberCommand{FamilyId: fam.Id(), Member: dup}

		repo.EXPECT().GetById(mock.Anything, fam.Id()).Return(fam, nil)
		auth.ensureErr = nil

		beforeLen := fam.Members().Len()

		// Act
		res := h.Handle(ctx, cmd)

		// Assert
		require.True(t, res.IsSuccess())
		assert.Equal(t, beforeLen, fam.Members().Len(), "family members should be unchanged on duplicate add")
		repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	t.Run("Validation error after AddMember => fail and no save", func(t *testing.T) {
		// Arrange: invalid name triggers validation error regardless of members
		repo := ports.NewMockFamilyRepository(t)
		auth := &authStub{perm: &permissionRequestStub{err: nil}}
		h := command.NewCreateFamilyMemberCommandHandler(repo, auth)

		fam := newFamily("owner1", "x", []family.Member{newMember(uuid.New(), "John")}) // invalid name (<3)
		cmd := command.CreateFamilyMemberCommand{FamilyId: fam.Id(), Member: newMember(fam.Id(), "Jane")}

		repo.EXPECT().GetById(mock.Anything, fam.Id()).Return(fam, nil)
		auth.ensureErr = nil

		// Act
		res := h.Handle(ctx, cmd)

		// Assert
		receivedErr := result.Match(res, func(_ family.Family) error { return nil }, func(err error) error { return err })
		require.Error(t, receivedErr)
		repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	t.Run("Repository Save returns error", func(t *testing.T) {
		// Arrange
		repo := ports.NewMockFamilyRepository(t)
		auth := &authStub{perm: &permissionRequestStub{err: nil}, ensureErr: nil}
		h := command.NewCreateFamilyMemberCommandHandler(repo, auth)

		fam := newFamily("owner1", "Doe family", []family.Member{newMember(uuid.New(), "John")})
		newM := newMember(fam.Id(), "Jane")
		cmd := command.CreateFamilyMemberCommand{FamilyId: fam.Id(), Member: newM}

		repo.EXPECT().GetById(mock.Anything, fam.Id()).Return(fam, nil)
		repo.EXPECT().Save(mock.Anything, mock.Anything).Return(errors.New("save failed"))

		// Act
		res := h.Handle(ctx, cmd)

		// Assert
		require.True(t, res.IsFaulted())
	})

	t.Run("Successful add -> success and save called", func(t *testing.T) {
		// Arrange
		repo := ports.NewMockFamilyRepository(t)
		auth := &authStub{perm: &permissionRequestStub{err: nil}, ensureErr: nil}
		h := command.NewCreateFamilyMemberCommandHandler(repo, auth)

		fam := newFamily("owner1", "Doe family", []family.Member{newMember(uuid.New(), "John")})
		newM := newMember(fam.Id(), "Jane")
		cmd := command.CreateFamilyMemberCommand{FamilyId: fam.Id(), Member: newM}

		repo.EXPECT().GetById(mock.Anything, fam.Id()).Return(fam, nil)
		repo.EXPECT().Save(mock.Anything, mock.Anything).Return(nil)

		// Act
		res := h.Handle(ctx, cmd)

		// Assert
		require.True(t, res.IsSuccess())
	})
}

// Helpers to satisfy variadic mock matching for Save which expects a slice of family.Family as second arg in mock
func mockAnyFamilies() []family.Family { return nil }

func mockFamiliesContaining(f family.Family) []family.Family { return []family.Family{f} }
