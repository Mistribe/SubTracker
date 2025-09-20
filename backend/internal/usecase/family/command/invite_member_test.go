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


func TestInviteMemberCommandHandler(t *testing.T) {
    ctx := context.Background()

    t.Run("GetById returns error", func(t *testing.T) {
        repo := ports.NewMockFamilyRepository(t)
        permReq := ports.NewMockPermissionRequest(t)
        auth := stubAuthorization{req: permReq}
        h := command.NewInviteMemberCommandHandler(repo, auth)

        famID := uuid.New()
        cmd := command.InviteMemberCommand{FamilyId: famID, FamilyMemberId: uuid.New()}

        expected := errors.New("db error")
        repo.EXPECT().GetById(mock.Anything, famID).Return(nil, expected)

        res := h.Handle(ctx, cmd)
        require.True(t, res.IsFaulted())
    })

    t.Run("Family not found returns ErrFamilyNotFound", func(t *testing.T) {
        repo := ports.NewMockFamilyRepository(t)
        permReq := ports.NewMockPermissionRequest(t)
        auth := stubAuthorization{req: permReq}
        h := command.NewInviteMemberCommandHandler(repo, auth)

        famID := uuid.New()
        cmd := command.InviteMemberCommand{FamilyId: famID, FamilyMemberId: uuid.New()}

        repo.EXPECT().GetById(mock.Anything, famID).Return(nil, nil)

        res := h.Handle(ctx, cmd)
        err := result.Match(res, func(_ command.InviteMemberResponse) error { return nil }, func(err error) error { return err })
        require.Error(t, err)
        assert.True(t, errors.Is(err, family.ErrFamilyNotFound))
    })

    t.Run("Authorization denies returns error", func(t *testing.T) {
        repo := ports.NewMockFamilyRepository(t)
        permReq := ports.NewMockPermissionRequest(t)
        auth := stubAuthorization{req: permReq}
        h := command.NewInviteMemberCommandHandler(repo, auth)

        famID := uuid.New()
        fam := makeFamilyWithMembers(famID, nil)
        cmd := command.InviteMemberCommand{FamilyId: famID, FamilyMemberId: uuid.New()}

        repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
        expected := errors.New("forbidden")
        permReq.EXPECT().For(fam).Return(expected)

        res := h.Handle(ctx, cmd)
        err := result.Match(res, func(_ command.InviteMemberResponse) error { return nil }, func(err error) error { return err })
        require.Error(t, err)
        assert.True(t, errors.Is(err, expected))
        repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
    })

    t.Run("Existing member with user already set returns ErrCannotInviteUser", func(t *testing.T) {
        repo := ports.NewMockFamilyRepository(t)
        permReq := ports.NewMockPermissionRequest(t)
        auth := stubAuthorization{req: permReq}
        h := command.NewInviteMemberCommandHandler(repo, auth)

        famID := uuid.New()
        m := makeMember(famID, "John")
        userID := "user-1"
        m.SetUserId(&userID)
        fam := makeFamilyWithMembers(famID, []family.Member{m})
        cmd := command.InviteMemberCommand{FamilyId: famID, FamilyMemberId: m.Id()}

        repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
        permReq.EXPECT().For(fam).Return(nil)

        res := h.Handle(ctx, cmd)

        err := result.Match(res, func(_ command.InviteMemberResponse) error { return nil }, func(err error) error { return err })
        require.Error(t, err)
        assert.True(t, errors.Is(err, family.ErrCannotInviteUser))
        repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
    })

    t.Run("Existing member invite success updates code and saves", func(t *testing.T) {
        repo := ports.NewMockFamilyRepository(t)
        permReq := ports.NewMockPermissionRequest(t)
        auth := stubAuthorization{req: permReq}
        h := command.NewInviteMemberCommandHandler(repo, auth)

        famID := uuid.New()
        m := makeMember(famID, "John")
        fam := makeFamilyWithMembers(famID, []family.Member{m})
        cmd := command.InviteMemberCommand{FamilyId: famID, FamilyMemberId: m.Id()}

        repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
        permReq.EXPECT().For(fam).Return(nil)
        repo.EXPECT().Save(mock.Anything, mock.Anything).Return(nil)

        res := h.Handle(ctx, cmd)

        require.True(t, res.IsSuccess())
        out := result.Match(res, func(v command.InviteMemberResponse) command.InviteMemberResponse { return v }, func(err error) command.InviteMemberResponse { t.Fatalf("unexpected error: %v", err); return command.InviteMemberResponse{} })
        require.NotEmpty(t, out.Code)
        assert.Equal(t, famID, out.FamilyId)
        assert.Equal(t, m.Id(), out.FamilyMemberId)
        updated := fam.GetMember(m.Id())
        require.NotNil(t, updated)
        require.NotNil(t, updated.InvitationCode())
        assert.Equal(t, out.Code, *updated.InvitationCode())
    })

    t.Run("New member invite success with default name and type", func(t *testing.T) {
        repo := ports.NewMockFamilyRepository(t)
        permReq := ports.NewMockPermissionRequest(t)
        auth := stubAuthorization{req: permReq}
        h := command.NewInviteMemberCommandHandler(repo, auth)

        famID := uuid.New()
        fam := makeFamilyWithMembers(famID, nil)
        cmd := command.InviteMemberCommand{FamilyId: famID, FamilyMemberId: uuid.New(), Name: nil, Type: nil}

        repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
        permReq.EXPECT().For(fam).Return(nil)
        repo.EXPECT().Save(mock.Anything, mock.Anything).Return(nil)

        res := h.Handle(ctx, cmd)

        require.True(t, res.IsSuccess())
        out := result.Match(res, func(v command.InviteMemberResponse) command.InviteMemberResponse { return v }, func(err error) command.InviteMemberResponse { t.Fatalf("unexpected error: %v", err); return command.InviteMemberResponse{} })
        require.NotEmpty(t, out.Code)
        // The handler creates a new member; ensure it's present and has defaults
        created := fam.GetMember(out.FamilyMemberId)
        require.NotNil(t, created)
        assert.Equal(t, "Invited User", created.Name())
        assert.Equal(t, family.AdultMemberType, created.Type())
        require.NotNil(t, created.InvitationCode())
        assert.Equal(t, out.Code, *created.InvitationCode())
    })

    t.Run("Existing member: validation error prevents Save", func(t *testing.T) {
        repo := ports.NewMockFamilyRepository(t)
        permReq := ports.NewMockPermissionRequest(t)
        auth := stubAuthorization{req: permReq}
        h := command.NewInviteMemberCommandHandler(repo, auth)

        famID := uuid.New()
        // invalid family name (length < 3) triggers validation errors
        m := makeMember(famID, "John")
        fam := family.NewFamily(famID, "owner-1", "x", []family.Member{m}, time.Now(), time.Now())
        cmd := command.InviteMemberCommand{FamilyId: famID, FamilyMemberId: m.Id()}

        repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
        permReq.EXPECT().For(fam).Return(nil)

        res := h.Handle(ctx, cmd)

        require.True(t, res.IsFaulted())
        repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
    })

    t.Run("Existing member: repository Save error propagates", func(t *testing.T) {
        repo := ports.NewMockFamilyRepository(t)
        permReq := ports.NewMockPermissionRequest(t)
        auth := stubAuthorization{req: permReq}
        h := command.NewInviteMemberCommandHandler(repo, auth)

        famID := uuid.New()
        m := makeMember(famID, "John")
        fam := makeFamilyWithMembers(famID, []family.Member{m})
        cmd := command.InviteMemberCommand{FamilyId: famID, FamilyMemberId: m.Id()}

        repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
        permReq.EXPECT().For(fam).Return(nil)
        repo.EXPECT().Save(mock.Anything, mock.Anything).Return(errors.New("save failed"))

        res := h.Handle(ctx, cmd)

        require.True(t, res.IsFaulted())
    })
}
