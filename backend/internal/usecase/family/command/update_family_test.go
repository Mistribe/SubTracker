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
    "github.com/mistribe/subtracker/pkg/langext/option"
    "github.com/mistribe/subtracker/pkg/langext/result"
)

func TestUpdateFamilyCommandHandler(t *testing.T) {
    ctx := context.Background()

    t.Run("GetById returns error", func(t *testing.T) {
        repo := ports.NewMockFamilyRepository(t)
        auth := stubAuthorization{}
        h := command.NewUpdateFamilyCommandHandler(repo, auth)

        famID := uuid.New()
        cmd := command.UpdateFamilyCommand{Id: famID, Name: "New Name", UpdatedAt: option.None[time.Time]()} // UpdatedAt won't be used

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
        auth := stubAuthorization{}
        h := command.NewUpdateFamilyCommandHandler(repo, auth)

        famID := uuid.New()
        cmd := command.UpdateFamilyCommand{Id: famID, Name: "New Name", UpdatedAt: option.None[time.Time]()} 

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
        auth := stubAuthorization{req: permReq}
        h := command.NewUpdateFamilyCommandHandler(repo, auth)

        famID := uuid.New()
        fam := makeFamilyWithMembers(famID, nil)
        cmd := command.UpdateFamilyCommand{Id: famID, Name: "New Name", UpdatedAt: option.None[time.Time]()} 

        expected := errors.New("forbidden")
        repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
        permReq.EXPECT().For(fam).Return(expected)

        res := h.Handle(ctx, cmd)

        err := result.Match(res, func(_ family.Family) error { return nil }, func(err error) error { return err })
        require.Error(t, err)
        assert.True(t, errors.Is(err, expected))
        repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
    })

    t.Run("Validation error after update -> fault and no Save", func(t *testing.T) {
        repo := ports.NewMockFamilyRepository(t)
        permReq := ports.NewMockPermissionRequest(t)
        auth := stubAuthorization{req: permReq}
        h := command.NewUpdateFamilyCommandHandler(repo, auth)

        famID := uuid.New()
        // Start with valid family, then set invalid name (too short) in command
        m := makeMember(famID, "John")
        fam := makeFamilyWithMembers(famID, []family.Member{m})
        cmd := command.UpdateFamilyCommand{Id: famID, Name: "x", UpdatedAt: option.None[time.Time]()} 

        repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
        permReq.EXPECT().For(fam).Return(nil)

        res := h.Handle(ctx, cmd)

        require.True(t, res.IsFaulted())
        repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
    })

    t.Run("Repository Save error propagates", func(t *testing.T) {
        repo := ports.NewMockFamilyRepository(t)
        permReq := ports.NewMockPermissionRequest(t)
        auth := stubAuthorization{req: permReq}
        h := command.NewUpdateFamilyCommandHandler(repo, auth)

        famID := uuid.New()
        m := makeMember(famID, "John")
        fam := makeFamilyWithMembers(famID, []family.Member{m})
        newName := "Updated Family"
        cmd := command.UpdateFamilyCommand{Id: famID, Name: newName, UpdatedAt: option.Some(time.Now())}

        repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
        permReq.EXPECT().For(fam).Return(nil)
        repo.EXPECT().Save(mock.Anything, mock.Anything).Return(errors.New("save failed"))

        res := h.Handle(ctx, cmd)

        require.True(t, res.IsFaulted())
    })

    t.Run("Success with UpdatedAt provided (Some) sets exact time and name", func(t *testing.T) {
        repo := ports.NewMockFamilyRepository(t)
        permReq := ports.NewMockPermissionRequest(t)
        auth := stubAuthorization{req: permReq}
        h := command.NewUpdateFamilyCommandHandler(repo, auth)

        famID := uuid.New()
        m := makeMember(famID, "John")
        fam := makeFamilyWithMembers(famID, []family.Member{m})
        newName := "Updated Family"
        ts := time.Date(2024, 12, 25, 10, 30, 0, 0, time.UTC)
        cmd := command.UpdateFamilyCommand{Id: famID, Name: newName, UpdatedAt: option.Some(ts)}

        repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
        permReq.EXPECT().For(fam).Return(nil)
        repo.EXPECT().Save(mock.Anything, mock.Anything).Return(nil)

        res := h.Handle(ctx, cmd)

        require.True(t, res.IsSuccess())
        got := result.Match(res, func(v family.Family) family.Family { return v }, func(err error) family.Family { t.Fatalf("unexpected error: %v", err); return nil })
        require.NotNil(t, got)
        assert.Equal(t, newName, got.Name())
        assert.Equal(t, ts, got.UpdatedAt())
    })

    t.Run("Success with UpdatedAt None sets time to now-ish", func(t *testing.T) {
        repo := ports.NewMockFamilyRepository(t)
        permReq := ports.NewMockPermissionRequest(t)
        auth := stubAuthorization{req: permReq}
        h := command.NewUpdateFamilyCommandHandler(repo, auth)

        famID := uuid.New()
        m := makeMember(famID, "John")
        fam := makeFamilyWithMembers(famID, []family.Member{m})
        newName := "Updated Family"
        cmd := command.UpdateFamilyCommand{Id: famID, Name: newName, UpdatedAt: option.None[time.Time]()}

        repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
        permReq.EXPECT().For(fam).Return(nil)
        repo.EXPECT().Save(mock.Anything, mock.Anything).Return(nil)

        before := time.Now()
        res := h.Handle(ctx, cmd)
        after := time.Now()

        require.True(t, res.IsSuccess())
        got := result.Match(res, func(v family.Family) family.Family { return v }, func(err error) family.Family { t.Fatalf("unexpected error: %v", err); return nil })
        require.NotNil(t, got)
        assert.Equal(t, newName, got.Name())
        // UpdatedAt should be between before and after
        assert.True(t, !got.UpdatedAt().Before(before) && !got.UpdatedAt().After(after), "updatedAt not within expected window: %v not in [%v, %v]", got.UpdatedAt(), before, after)
    })
}
