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

type stubAuthorization struct{ req ports.PermissionRequest }

func (s stubAuthorization) EnsureWithinLimit(ctx context.Context, feature user.Feature, delta int) error {
    return nil
}
func (s stubAuthorization) Can(ctx context.Context, permission user.Permission) ports.PermissionRequest {
    return s.req
}
func (s stubAuthorization) GetCurrentLimit(ctx context.Context, feature user.Feature) (shared.Limit, error) {
    return nil, nil
}
func (s stubAuthorization) GetCurrentLimits(ctx context.Context, category user.Category) (shared.Limits, error) {
    return nil, nil
}

func makeFamily(famID uuid.UUID) family.Family {
    m := family.NewMember(uuid.New(), famID, "John", family.OwnerMemberType, nil, time.Now(), time.Now())
    return family.NewFamily(famID, "owner-1", "Doe family", []family.Member{m}, time.Now(), time.Now())
}

func TestDeleteFamilyCommandHandler(t *testing.T) {
    ctx := context.Background()

    t.Run("GetById returns error", func(t *testing.T) {
        repo := ports.NewMockFamilyRepository(t)
        auth := stubAuthorization{}
        h := command.NewDeleteFamilyCommandHandler(repo, auth)

        famID := uuid.New()
        cmd := command.DeleteFamilyCommand{FamilyId: famID}

        expectedErr := errors.New("db error")
        repo.EXPECT().GetById(mock.Anything, famID).Return(nil, expectedErr)

        res := h.Handle(ctx, cmd)

        receivedErr := result.Match(res, func(_ bool) error { return nil }, func(err error) error { return err })
        require.Error(t, receivedErr)
        assert.True(t, errors.Is(receivedErr, expectedErr))
        repo.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
    })

    t.Run("Family not found returns ErrFamilyNotFound", func(t *testing.T) {
        repo := ports.NewMockFamilyRepository(t)
        auth := stubAuthorization{}
        h := command.NewDeleteFamilyCommandHandler(repo, auth)

        famID := uuid.New()
        cmd := command.DeleteFamilyCommand{FamilyId: famID}

        repo.EXPECT().GetById(mock.Anything, famID).Return(nil, nil)

        res := h.Handle(ctx, cmd)

        receivedErr := result.Match(res, func(_ bool) error { return nil }, func(err error) error { return err })
        require.Error(t, receivedErr)
        assert.True(t, errors.Is(receivedErr, family.ErrFamilyNotFound))
        repo.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
    })

    t.Run("Authorization denies returns error", func(t *testing.T) {
        repo := ports.NewMockFamilyRepository(t)
        permReq := ports.NewMockPermissionRequest(t)
        auth := stubAuthorization{req: permReq}
        h := command.NewDeleteFamilyCommandHandler(repo, auth)

        famID := uuid.New()
        fam := makeFamily(famID)
        cmd := command.DeleteFamilyCommand{FamilyId: famID}

        expectedErr := errors.New("forbidden")
        repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
        permReq.EXPECT().For(fam).Return(expectedErr)

        res := h.Handle(ctx, cmd)

        receivedErr := result.Match(res, func(_ bool) error { return nil }, func(err error) error { return err })
        require.Error(t, receivedErr)
        assert.True(t, errors.Is(receivedErr, expectedErr))
        repo.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
    })

    t.Run("Repository Delete returns error", func(t *testing.T) {
        repo := ports.NewMockFamilyRepository(t)
        permReq := ports.NewMockPermissionRequest(t)
        auth := stubAuthorization{req: permReq}
        h := command.NewDeleteFamilyCommandHandler(repo, auth)

        famID := uuid.New()
        fam := makeFamily(famID)
        cmd := command.DeleteFamilyCommand{FamilyId: famID}

        repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
        permReq.EXPECT().For(fam).Return(nil)

        expectedErr := errors.New("delete failed")
        repo.EXPECT().Delete(mock.Anything, famID).Return(false, expectedErr)

        res := h.Handle(ctx, cmd)

        receivedErr := result.Match(res, func(_ bool) error { return nil }, func(err error) error { return err })
        require.Error(t, receivedErr)
        assert.True(t, errors.Is(receivedErr, expectedErr))
    })

    t.Run("Repository Delete returns false -> ErrFamilyNotFound", func(t *testing.T) {
        repo := ports.NewMockFamilyRepository(t)
        permReq := ports.NewMockPermissionRequest(t)
        auth := stubAuthorization{req: permReq}
        h := command.NewDeleteFamilyCommandHandler(repo, auth)

        famID := uuid.New()
        fam := makeFamily(famID)
        cmd := command.DeleteFamilyCommand{FamilyId: famID}

        repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
        permReq.EXPECT().For(fam).Return(nil)
        repo.EXPECT().Delete(mock.Anything, famID).Return(false, nil)

        res := h.Handle(ctx, cmd)

        receivedErr := result.Match(res, func(_ bool) error { return nil }, func(err error) error { return err })
        require.Error(t, receivedErr)
        assert.True(t, errors.Is(receivedErr, family.ErrFamilyNotFound))
    })

    t.Run("Successful deletion returns success true", func(t *testing.T) {
        repo := ports.NewMockFamilyRepository(t)
        permReq := ports.NewMockPermissionRequest(t)
        auth := stubAuthorization{req: permReq}
        h := command.NewDeleteFamilyCommandHandler(repo, auth)

        famID := uuid.New()
        fam := makeFamily(famID)
        cmd := command.DeleteFamilyCommand{FamilyId: famID}

        repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
        permReq.EXPECT().For(fam).Return(nil)
        repo.EXPECT().Delete(mock.Anything, famID).Return(true, nil)

        res := h.Handle(ctx, cmd)

        require.True(t, res.IsSuccess())
        ok := result.Match(res, func(v bool) bool { return v }, func(err error) bool { return false })
        assert.True(t, ok)
    })
}
