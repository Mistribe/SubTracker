package command_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/domain/user"
	"github.com/mistribe/subtracker/internal/usecase/family/command"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type fakeFamilyRepo struct {
	getByIdFn func(ctx context.Context, id uuid.UUID) (family.Family, error)
	saveFn    func(ctx context.Context, entities ...family.Family) error

	lastGetById uuid.UUID
	lastSaved   []family.Family
	saveCalls   int
}

func (f *fakeFamilyRepo) GetById(ctx context.Context, id uuid.UUID) (family.Family, error) {
	f.lastGetById = id
	if f.getByIdFn != nil {
		return f.getByIdFn(ctx, id)
	}
	return nil, nil
}

func (f *fakeFamilyRepo) Save(ctx context.Context, entities ...family.Family) error {
	f.saveCalls++
	f.lastSaved = entities
	if f.saveFn != nil {
		return f.saveFn(ctx, entities...)
	}
	return nil
}

// Unused in these tests but required by interface
func (f *fakeFamilyRepo) Delete(ctx context.Context, entityId uuid.UUID) (bool, error) {
	return false, nil
}
func (f *fakeFamilyRepo) Exists(ctx context.Context, ids ...uuid.UUID) (bool, error) {
	return false, nil
}
func (f *fakeFamilyRepo) GetUserFamily(ctx context.Context, userId string) (family.Family, error) {
	return nil, nil
}
func (f *fakeFamilyRepo) MemberExists(ctx context.Context, familyId uuid.UUID, members ...uuid.UUID) (bool, error) {
	return false, nil
}
func (f *fakeFamilyRepo) IsUserMemberOfFamily(ctx context.Context, familyId uuid.UUID, userId string) (bool, error) {
	return false, nil
}

type fakeAuth struct{ userId string }

func (a fakeAuth) MustGetUserId(ctx context.Context) string      { return a.userId }
func (a fakeAuth) MustGetUserRole(ctx context.Context) user.Role { return user.RoleUser }

func TestCreateFamilyCommandHandler_Handle(t *testing.T) {
	ctx := context.Background()
	const currentUser = "user-123"

	t.Run("returns ErrFamilyAlreadyExists when provided FamilyId already exists", func(t *testing.T) {
		repo := &fakeFamilyRepo{}
		providedId := uuid.MustParse("00000000-0000-7000-8000-000000000001")
		repo.getByIdFn = func(ctx context.Context, id uuid.UUID) (family.Family, error) {
			return family.NewFamily(id, "other", "ExistingFam", nil, time.Now(), time.Now()), nil
		}
		h := command.NewCreateFamilyCommandHandler(repo, fakeAuth{userId: currentUser})
		res := h.Handle(ctx, command.CreateFamilyCommand{FamilyId: &providedId, Name: "Any", CreatorName: "Owner"})

		err := result.Match(res, func(_ family.Family) error { return nil }, func(e error) error { return e })
		if !errors.Is(err, family.ErrFamilyAlreadyExists) {
			t.Fatalf("expected ErrFamilyAlreadyExists, got %v", err)
		}
	})

	t.Run("propagates repository error from GetById when provided FamilyId", func(t *testing.T) {
		repo := &fakeFamilyRepo{}
		providedId := uuid.MustParse("00000000-0000-7000-8000-000000000002")
		expectedErr := errors.New("db error")
		repo.getByIdFn = func(ctx context.Context, id uuid.UUID) (family.Family, error) { return nil, expectedErr }
		h := command.NewCreateFamilyCommandHandler(repo, fakeAuth{userId: currentUser})
		res := h.Handle(ctx, command.CreateFamilyCommand{FamilyId: &providedId, Name: "Any", CreatorName: "Owner"})

		err := result.Match(res, func(_ family.Family) error { return nil }, func(e error) error { return e })
		if !errors.Is(err, expectedErr) {
			t.Fatalf("expected error %v, got %v", expectedErr, err)
		}
	})

	t.Run("succeeds when FamilyId is nil (new id generated)", func(t *testing.T) {
		repo := &fakeFamilyRepo{}
		h := command.NewCreateFamilyCommandHandler(repo, fakeAuth{userId: currentUser})
		res := h.Handle(ctx, command.CreateFamilyCommand{Name: "My Family", CreatorName: "Owner"})
		if !res.IsSuccess() {
			t.Fatalf("expected success, got failure")
		}
	})

	t.Run("fails when name is empty due to validation errors", func(t *testing.T) {
		repo := &fakeFamilyRepo{}
		h := command.NewCreateFamilyCommandHandler(repo, fakeAuth{userId: currentUser})
		res := h.Handle(ctx, command.CreateFamilyCommand{Name: "", CreatorName: "Owner"})
		if !res.IsFaulted() {
			t.Fatalf("expected failure due to validation errors, got success")
		}
	})

	t.Run("propagates repository error from Save", func(t *testing.T) {
		repo := &fakeFamilyRepo{}
		repo.saveFn = func(ctx context.Context, entities ...family.Family) error { return errors.New("save failed") }
		h := command.NewCreateFamilyCommandHandler(repo, fakeAuth{userId: currentUser})
		res := h.Handle(ctx, command.CreateFamilyCommand{Name: "My Family", CreatorName: "Owner"})
		if !res.IsFaulted() {
			t.Fatalf("expected failure when save fails, got success")
		}
	})

	t.Run("calls repository Save exactly once on success", func(t *testing.T) {
		repo := &fakeFamilyRepo{}
		h := command.NewCreateFamilyCommandHandler(repo, fakeAuth{userId: currentUser})
		_ = h.Handle(ctx, command.CreateFamilyCommand{Name: "My Family", CreatorName: "Owner"})
		if repo.saveCalls != 1 {
			t.Fatalf("expected Save to be called once, got %d", repo.saveCalls)
		}
	})

	t.Run("sets owner user id to current user on created family", func(t *testing.T) {
		repo := &fakeFamilyRepo{}
		h := command.NewCreateFamilyCommandHandler(repo, fakeAuth{userId: currentUser})
		res := h.Handle(ctx, command.CreateFamilyCommand{Name: "My Family", CreatorName: "Owner"})
		fam := result.Match(res, func(v family.Family) family.Family { return v },
			func(err error) family.Family { return nil })
		if fam == nil || fam.Owner().UserId() != currentUser {
			t.Fatalf("expected owner user id %s, got %v", currentUser, fam)
		}
	})

	t.Run("adds exactly one member to the created family", func(t *testing.T) {
		repo := &fakeFamilyRepo{}
		h := command.NewCreateFamilyCommandHandler(repo, fakeAuth{userId: currentUser})
		res := h.Handle(ctx, command.CreateFamilyCommand{Name: "My Family", CreatorName: "Owner"})
		fam := result.Match(res, func(v family.Family) family.Family { return v },
			func(err error) family.Family { return nil })
		count := 0
		if fam != nil {
			for range fam.Members().It() {
				count++
			}
		}
		if count != 1 {
			t.Fatalf("expected exactly one member, got %d", count)
		}
	})

	t.Run("sets created member type to owner", func(t *testing.T) {
		repo := &fakeFamilyRepo{}
		h := command.NewCreateFamilyCommandHandler(repo, fakeAuth{userId: currentUser})
		res := h.Handle(ctx, command.CreateFamilyCommand{Name: "My Family", CreatorName: "Owner"})
		fam := result.Match(res, func(v family.Family) family.Family { return v },
			func(err error) family.Family { return nil })
		var memberType family.MemberType
		if fam != nil {
			for m := range fam.Members().It() {
				memberType = m.Type()
				break
			}
		}
		if memberType != family.OwnerMemberType {
			t.Fatalf("expected member type %s, got %s", family.OwnerMemberType, memberType)
		}
	})

	t.Run("sets creator name as member name", func(t *testing.T) {
		repo := &fakeFamilyRepo{}
		h := command.NewCreateFamilyCommandHandler(repo, fakeAuth{userId: currentUser})
		const creatorName = "John Doe"
		res := h.Handle(ctx, command.CreateFamilyCommand{Name: "My Family", CreatorName: creatorName})
		fam := result.Match(res, func(v family.Family) family.Family { return v },
			func(err error) family.Family { return nil })
		var name string
		if fam != nil {
			for m := range fam.Members().It() {
				name = m.Name()
				break
			}
		}
		if name != creatorName {
			t.Fatalf("expected member name %q, got %q", creatorName, name)
		}
	})

	t.Run("sets member user id to current user", func(t *testing.T) {
		repo := &fakeFamilyRepo{}
		h := command.NewCreateFamilyCommandHandler(repo, fakeAuth{userId: currentUser})
		res := h.Handle(ctx, command.CreateFamilyCommand{Name: "My Family", CreatorName: "Owner"})
		fam := result.Match(res, func(v family.Family) family.Family { return v },
			func(err error) family.Family { return nil })
		var uid string
		if fam != nil {
			for m := range fam.Members().It() {
				if m.UserId() != nil {
					uid = *m.UserId()
				}
				break
			}
		}
		if uid != currentUser {
			t.Fatalf("expected member user id %q, got %q", currentUser, uid)
		}
	})

	t.Run("uses provided CreatedAt when present", func(t *testing.T) {
		repo := &fakeFamilyRepo{}
		h := command.NewCreateFamilyCommandHandler(repo, fakeAuth{userId: currentUser})
		customTime := time.Date(2024, 7, 21, 10, 9, 8, 0, time.UTC)
		res := h.Handle(ctx,
			command.CreateFamilyCommand{Name: "My Family", CreatorName: "Owner", CreatedAt: &customTime})
		fam := result.Match(res, func(v family.Family) family.Family { return v },
			func(err error) family.Family { return nil })
		if fam == nil || !fam.CreatedAt().Equal(customTime) {
			t.Fatalf("expected CreatedAt %v, got %v", customTime, fam)
		}
	})
}
