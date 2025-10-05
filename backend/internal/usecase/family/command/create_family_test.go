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

	"github.com/mistribe/subtracker/internal/domain/account"
	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/family/command"
	"github.com/mistribe/subtracker/pkg/langext/option"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

// simple stub for Authentication to reduce brittleness of mocks
type stubAuth struct{ user types.UserID }

func (s stubAuth) MustGetConnectedAccount(ctx context.Context) account.ConnectedAccount {
	return stubAccount{s.user}
}

type stubAccount struct{ uid types.UserID }

func (a stubAccount) UserID() types.UserID { return a.uid }
func (a stubAccount) PlanID() types.PlanID { return types.PlanFree }
func (a stubAccount) Role() types.Role     { return types.RoleUser }

func TestCreateFamilyCommandHandler_Handle(t *testing.T) {
	ctx := context.Background()
	const currentUserStr = "user-123"

	// returns ErrFamilyAlreadyExists when provided FamilyID already exists
	// (no auth call path)
	t.Run("returns ErrFamilyAlreadyExists when provided FamilyID already exists", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		auth := ports.NewMockAuthentication(t)
		h := command.NewCreateFamilyCommandHandler(repo, auth)

		providedId := types.FamilyID(uuid.MustParse("00000000-0000-7000-8000-000000000001"))
		cmd := command.CreateFamilyCommand{FamilyId: option.Some(providedId), Name: "Any", CreatorName: "Owner"}

		repo.EXPECT().GetById(mock.Anything, providedId).Return(family.NewFamily(providedId, types.UserID("other"),
			"ExistingFam", []family.Member{}, time.Now(), time.Now()), nil)

		res := h.Handle(ctx, cmd)

		err := result.Match(res, func(_ family.Family) error { return nil }, func(e error) error { return e })
		require.Error(t, err)
		assert.True(t, errors.Is(err, family.ErrFamilyAlreadyExists))
		repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	// propagates repository error from GetById when provided FamilyID
	t.Run("propagates repository error from GetById when provided FamilyID", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		auth := ports.NewMockAuthentication(t)
		h := command.NewCreateFamilyCommandHandler(repo, auth)

		providedId := types.FamilyID(uuid.MustParse("00000000-0000-7000-8000-000000000002"))
		expectedErr := errors.New("db error")
		cmd := command.CreateFamilyCommand{FamilyId: option.Some(providedId), Name: "Any", CreatorName: "Owner"}

		repo.EXPECT().GetById(mock.Anything, providedId).Return(nil, expectedErr)

		res := h.Handle(ctx, cmd)
		err := result.Match(res, func(_ family.Family) error { return nil }, func(e error) error { return e })
		require.Error(t, err)
		assert.True(t, errors.Is(err, expectedErr))
		repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	// succeeds when FamilyID is None (new id generated)
	t.Run("succeeds when FamilyID is None (new id generated)", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		auth := stubAuth{user: types.UserID(currentUserStr)}
		h := command.NewCreateFamilyCommandHandler(repo, auth)

		repo.EXPECT().Save(mock.Anything, mock.Anything).Return(nil)

		res := h.Handle(ctx, command.CreateFamilyCommand{Name: "My Family", CreatorName: "Owner"})
		require.True(t, res.IsSuccess())
	})

	// fails when name is empty due to validation errors
	t.Run("fails when name is empty due to validation errors", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		auth := stubAuth{user: types.UserID(currentUserStr)}
		h := command.NewCreateFamilyCommandHandler(repo, auth)

		res := h.Handle(ctx, command.CreateFamilyCommand{Name: "", CreatorName: "Owner"})
		require.True(t, res.IsFaulted())
		repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	// propagates repository error from Save
	t.Run("propagates repository error from Save", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		auth := stubAuth{user: types.UserID(currentUserStr)}
		h := command.NewCreateFamilyCommandHandler(repo, auth)

		repo.EXPECT().Save(mock.Anything, mock.Anything).Return(errors.New("save failed"))

		res := h.Handle(ctx, command.CreateFamilyCommand{Name: "My Family", CreatorName: "Owner"})
		require.True(t, res.IsFaulted())
	})

	// calls repository Save exactly once on success
	t.Run("calls repository Save exactly once on success", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		auth := stubAuth{user: types.UserID(currentUserStr)}
		h := command.NewCreateFamilyCommandHandler(repo, auth)

		repo.EXPECT().Save(mock.Anything, mock.Anything).Return(nil).Once()

		_ = h.Handle(ctx, command.CreateFamilyCommand{Name: "My Family", CreatorName: "Owner"})
		repo.AssertNumberOfCalls(t, "Save", 1)
	})

	// sets owner user id to current user on created family
	t.Run("sets owner user id to current user on created family", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		auth := stubAuth{user: types.UserID(currentUserStr)}
		h := command.NewCreateFamilyCommandHandler(repo, auth)

		repo.EXPECT().Save(mock.Anything, mock.Anything).Return(nil)

		res := h.Handle(ctx, command.CreateFamilyCommand{Name: "My Family", CreatorName: "Owner"})
		fam := result.Match(res, func(v family.Family) family.Family { return v },
			func(err error) family.Family { t.Fatalf("unexpected error: %v", err); return nil })
		require.NotNil(t, fam)
		assert.Equal(t, currentUserStr, fam.Owner().UserId().String())
	})

	// adds exactly one member to the created family
	t.Run("adds exactly one member to the created family", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		auth := stubAuth{user: types.UserID(currentUserStr)}
		h := command.NewCreateFamilyCommandHandler(repo, auth)
		repo.EXPECT().Save(mock.Anything, mock.Anything).Return(nil)

		res := h.Handle(ctx, command.CreateFamilyCommand{Name: "My Family", CreatorName: "Owner"})
		fam := result.Match(res, func(v family.Family) family.Family { return v },
			func(err error) family.Family { t.Fatalf("unexpected error: %v", err); return nil })
		require.NotNil(t, fam)
		assert.Equal(t, 1, fam.Members().Len())
	})

	// sets created member type to owner
	t.Run("sets created member type to owner", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		auth := stubAuth{user: types.UserID(currentUserStr)}
		h := command.NewCreateFamilyCommandHandler(repo, auth)
		repo.EXPECT().Save(mock.Anything, mock.Anything).Return(nil)

		res := h.Handle(ctx, command.CreateFamilyCommand{Name: "My Family", CreatorName: "Owner"})
		fam := result.Match(res, func(v family.Family) family.Family { return v },
			func(err error) family.Family { t.Fatalf("unexpected error: %v", err); return nil })
		var memberType family.MemberType
		for m := range fam.Members().It() {
			memberType = m.Type()
			break
		}
		assert.Equal(t, family.OwnerMemberType, memberType)
	})

	// sets creator name as member name
	t.Run("sets creator name as member name", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		auth := stubAuth{user: types.UserID(currentUserStr)}
		h := command.NewCreateFamilyCommandHandler(repo, auth)
		repo.EXPECT().Save(mock.Anything, mock.Anything).Return(nil)

		const creatorName = "John Doe"
		res := h.Handle(ctx, command.CreateFamilyCommand{Name: "My Family", CreatorName: creatorName})
		fam := result.Match(res, func(v family.Family) family.Family { return v },
			func(err error) family.Family { t.Fatalf("unexpected error: %v", err); return nil })
		var name string
		for m := range fam.Members().It() {
			name = m.Name()
			break
		}
		assert.Equal(t, creatorName, name)
	})

	// sets member user id to current user
	t.Run("sets member user id to current user", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		auth := stubAuth{user: types.UserID(currentUserStr)}
		h := command.NewCreateFamilyCommandHandler(repo, auth)
		repo.EXPECT().Save(mock.Anything, mock.Anything).Return(nil)

		res := h.Handle(ctx, command.CreateFamilyCommand{Name: "My Family", CreatorName: "Owner"})
		fam := result.Match(res, func(v family.Family) family.Family { return v },
			func(err error) family.Family { t.Fatalf("unexpected error: %v", err); return nil })
		var uid *types.UserID
		for m := range fam.Members().It() {
			uid = m.UserId()
			break
		}
		require.NotNil(t, uid)
		assert.Equal(t, currentUserStr, uid.String())
	})

	// uses provided CreatedAt when present
	t.Run("uses provided CreatedAt when present", func(t *testing.T) {
		repo := ports.NewMockFamilyRepository(t)
		auth := ports.NewMockAuthentication(t)
		fake := account.NewMockConnectedAccount(t)
		fake.EXPECT().UserID().Return(types.UserID(currentUserStr))
		auth.EXPECT().MustGetConnectedAccount(mock.Anything).Return(fake)
		h := command.NewCreateFamilyCommandHandler(repo, auth)
		customTime := time.Date(2024, 7, 21, 10, 9, 8, 0, time.UTC)

		repo.EXPECT().Save(mock.Anything, mock.Anything).Return(nil)

		res := h.Handle(ctx,
			command.CreateFamilyCommand{Name: "My Family", CreatorName: "Owner", CreatedAt: option.Some(customTime)})
		fam := result.Match(res, func(v family.Family) family.Family { return v },
			func(err error) family.Family { t.Fatalf("unexpected error: %v", err); return nil })
		require.NotNil(t, fam)
		assert.True(t, fam.CreatedAt().Equal(customTime), "expected CreatedAt %v got %v", customTime, fam.CreatedAt())
	})
}
