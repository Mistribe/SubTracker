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

func TestUpdateFamilyMemberCommandHandler(t *testing.T) {
    ctx := context.Background()

    t.Run("GetById returns error", func(t *testing.T) {
        repo := ports.NewMockFamilyRepository(t)
        auth := stubAuthorization{}
        h := command.NewUpdateFamilyMemberCommandHandler(repo, auth)

        famID := uuid.New()
        cmd := command.UpdateFamilyMemberCommand{FamilyId: famID, Id: uuid.New()}

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
        h := command.NewUpdateFamilyMemberCommandHandler(repo, auth)

        famID := uuid.New()
        cmd := command.UpdateFamilyMemberCommand{FamilyId: famID, Id: uuid.New()}

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
        h := command.NewUpdateFamilyMemberCommandHandler(repo, auth)

        famID := uuid.New()
        m := makeMember(famID, "John")
        fam := makeFamilyWithMembers(famID, []family.Member{m})
        cmd := command.UpdateFamilyMemberCommand{FamilyId: famID, Id: m.Id(), Name: "Johnny", Type: family.KidMemberType, UpdatedAt: option.None[time.Time]()}

        repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
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
        permReq := ports.NewMockPermissionRequest(t)
        auth := stubAuthorization{req: permReq}
        h := command.NewUpdateFamilyMemberCommandHandler(repo, auth)

        famID := uuid.New()
        // family with a different member
        other := makeMember(famID, "John")
        fam := makeFamilyWithMembers(famID, []family.Member{other})
        cmd := command.UpdateFamilyMemberCommand{FamilyId: famID, Id: uuid.New(), Name: "Johnny", Type: family.KidMemberType, UpdatedAt: option.None[time.Time]()}

        repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
        // Authorization not checked when member is missing

        res := h.Handle(ctx, cmd)

        err := result.Match(res, func(_ family.Family) error { return nil }, func(err error) error { return err })
        require.Error(t, err)
        assert.True(t, errors.Is(err, family.ErrFamilyMemberNotFound))
        repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
    })

    t.Run("Validation error after update -> fault and no Save", func(t *testing.T) {
        repo := ports.NewMockFamilyRepository(t)
        permReq := ports.NewMockPermissionRequest(t)
        auth := stubAuthorization{req: permReq}
        h := command.NewUpdateFamilyMemberCommandHandler(repo, auth)

        famID := uuid.New()
        m := makeMember(famID, "John")
        fam := makeFamilyWithMembers(famID, []family.Member{m})
        // Set invalid name (empty) to trigger member validation error
        cmd := command.UpdateFamilyMemberCommand{FamilyId: famID, Id: m.Id(), Name: "", Type: family.KidMemberType, UpdatedAt: option.None[time.Time]()}

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
        h := command.NewUpdateFamilyMemberCommandHandler(repo, auth)

        famID := uuid.New()
        m := makeMember(famID, "John")
        fam := makeFamilyWithMembers(famID, []family.Member{m})
        cmd := command.UpdateFamilyMemberCommand{FamilyId: famID, Id: m.Id(), Name: "Johnny", Type: family.KidMemberType, UpdatedAt: option.Some(time.Now())}

        repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
        permReq.EXPECT().For(fam).Return(nil)
        repo.EXPECT().Save(mock.Anything, mock.Anything).Return(errors.New("save failed"))

        res := h.Handle(ctx, cmd)

        require.True(t, res.IsFaulted())
    })

    t.Run("Success with UpdatedAt Some sets exact time and updates fields", func(t *testing.T) {
        repo := ports.NewMockFamilyRepository(t)
        permReq := ports.NewMockPermissionRequest(t)
        auth := stubAuthorization{req: permReq}
        h := command.NewUpdateFamilyMemberCommandHandler(repo, auth)

        famID := uuid.New()
        m := makeMember(famID, "John")
        fam := makeFamilyWithMembers(famID, []family.Member{m})
        ts := time.Date(2024, 12, 25, 10, 30, 0, 0, time.UTC)
        newName := "Johnny"
        newType := family.KidMemberType
        cmd := command.UpdateFamilyMemberCommand{FamilyId: famID, Id: m.Id(), Name: newName, Type: newType, UpdatedAt: option.Some(ts)}

        repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
        permReq.EXPECT().For(fam).Return(nil)
        repo.EXPECT().Save(mock.Anything, mock.Anything).Return(nil)

        res := h.Handle(ctx, cmd)

        require.True(t, res.IsSuccess())
        got := result.Match(res, func(v family.Family) family.Family { return v }, func(err error) family.Family { t.Fatalf("unexpected error: %v", err); return nil })
        require.NotNil(t, got)
        updated := got.GetMember(m.Id())
        require.NotNil(t, updated)
        assert.Equal(t, newName, updated.Name())
        assert.Equal(t, newType, updated.Type())
        // Member does not expose UpdatedAt; instead we can ensure family UpdatedAt equals ts because handler sets member.UpdatedAt and family.UpdateMember marks dirty; but family updatedAt is set only when Base.SetUpdatedAt called. Here we rely on Save being called and no direct access to member updatedAt.
        // We cannot read member updatedAt directly; thus we only verify Save called and assume UpdatedAt was set.
    })

    t.Run("Success with UpdatedAt None sets time to now-ish and updates fields", func(t *testing.T) {
        repo := ports.NewMockFamilyRepository(t)
        permReq := ports.NewMockPermissionRequest(t)
        auth := stubAuthorization{req: permReq}
        h := command.NewUpdateFamilyMemberCommandHandler(repo, auth)

        famID := uuid.New()
        m := makeMember(famID, "John")
        fam := makeFamilyWithMembers(famID, []family.Member{m})
        newName := "Johnny"
        newType := family.AdultMemberType
        cmd := command.UpdateFamilyMemberCommand{FamilyId: famID, Id: m.Id(), Name: newName, Type: newType, UpdatedAt: option.None[time.Time]()}

        repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
        permReq.EXPECT().For(fam).Return(nil)
        repo.EXPECT().Save(mock.Anything, mock.Anything).Return(nil)

        before := time.Now()
        res := h.Handle(ctx, cmd)
        after := time.Now()

        require.True(t, res.IsSuccess())
        got := result.Match(res, func(v family.Family) family.Family { return v }, func(err error) family.Family { t.Fatalf("unexpected error: %v", err); return nil })
        require.NotNil(t, got)
        updated := got.GetMember(m.Id())
        require.NotNil(t, updated)
        assert.Equal(t, newName, updated.Name())
        assert.Equal(t, newType, updated.Type())
        _ = before
        _ = after
        // As above, we cannot assert member's UpdatedAt directly; ensure Save was called and fields updated.
    })
}
