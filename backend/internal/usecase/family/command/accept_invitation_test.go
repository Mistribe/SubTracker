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

func newMemberWithInvite(famID uuid.UUID, name string, code *string) family.Member {
    return family.NewMember(
        uuid.New(),
        famID,
        name,
        family.AdultMemberType,
        code,
        time.Now(),
        time.Now(),
    )
}

func newFamilyWithMembers(id uuid.UUID, ownerID string, name string, members []family.Member) family.Family {
    return family.NewFamily(
        id,
        ownerID,
        name,
        members,
        time.Now(),
        time.Now(),
    )
}

func TestAcceptInvitationCommandHandler(t *testing.T) {
    ctx := context.Background()

    t.Run("GetById returns error", func(t *testing.T) {
        repo := ports.NewMockFamilyRepository(t)
        auth := ports.NewMockAuthentication(t)
        h := command.NewAcceptInvitationCommandHandler(repo, auth)

        famID := uuid.New()
        cmd := command.AcceptInvitationCommand{FamilyId: famID}

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
        auth := ports.NewMockAuthentication(t)
        h := command.NewAcceptInvitationCommandHandler(repo, auth)

        famID := uuid.New()
        cmd := command.AcceptInvitationCommand{FamilyId: famID}

        repo.EXPECT().GetById(mock.Anything, famID).Return(nil, nil)

        res := h.Handle(ctx, cmd)

        receivedErr := result.Match(res, func(_ bool) error { return nil }, func(err error) error { return err })
        require.Error(t, receivedErr)
        assert.True(t, errors.Is(receivedErr, family.ErrFamilyNotFound))
        repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
    })

    t.Run("Member not found returns ErrFamilyMemberNotFound", func(t *testing.T) {
        repo := ports.NewMockFamilyRepository(t)
        auth := ports.NewMockAuthentication(t)
        h := command.NewAcceptInvitationCommandHandler(repo, auth)

        famID := uuid.New()
        otherMember := newMemberWithInvite(famID, "John", nil)
        fam := newFamilyWithMembers(famID, "owner1", "Doe family", []family.Member{otherMember})
        cmd := command.AcceptInvitationCommand{FamilyId: famID, FamilyMemberId: uuid.New()}

        repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)

        res := h.Handle(ctx, cmd)

        receivedErr := result.Match(res, func(_ bool) error { return nil }, func(err error) error { return err })
        require.Error(t, receivedErr)
        assert.True(t, errors.Is(receivedErr, family.ErrFamilyMemberNotFound))
        repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
    })

    t.Run("Invitation code is nil returns ErrBadInvitationCode", func(t *testing.T) {
        repo := ports.NewMockFamilyRepository(t)
        auth := ports.NewMockAuthentication(t)
        h := command.NewAcceptInvitationCommandHandler(repo, auth)

        famID := uuid.New()
        m := newMemberWithInvite(famID, "John", nil)
        fam := newFamilyWithMembers(famID, "owner1", "Doe family", []family.Member{m})
        cmd := command.AcceptInvitationCommand{FamilyId: famID, FamilyMemberId: m.Id(), InvitationCode: "expected"}

        repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)

        res := h.Handle(ctx, cmd)

        receivedErr := result.Match(res, func(_ bool) error { return nil }, func(err error) error { return err })
        require.Error(t, receivedErr)
        assert.True(t, errors.Is(receivedErr, family.ErrBadInvitationCode))
        repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
    })

    t.Run("Invitation code mismatch returns ErrBadInvitationCode", func(t *testing.T) {
        repo := ports.NewMockFamilyRepository(t)
        auth := ports.NewMockAuthentication(t)
        h := command.NewAcceptInvitationCommandHandler(repo, auth)

        famID := uuid.New()
        code := "expected"
        m := newMemberWithInvite(famID, "John", &code)
        fam := newFamilyWithMembers(famID, "owner1", "Doe family", []family.Member{m})
        cmd := command.AcceptInvitationCommand{FamilyId: famID, FamilyMemberId: m.Id(), InvitationCode: "wrong"}

        repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)

        res := h.Handle(ctx, cmd)

        receivedErr := result.Match(res, func(_ bool) error { return nil }, func(err error) error { return err })
        require.Error(t, receivedErr)
        assert.True(t, errors.Is(receivedErr, family.ErrBadInvitationCode))
        repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
    })

    t.Run("Validation error after update -> fail and no save", func(t *testing.T) {
        repo := ports.NewMockFamilyRepository(t)
        auth := ports.NewMockAuthentication(t)
        h := command.NewAcceptInvitationCommandHandler(repo, auth)

        famID := uuid.New()
        code := "expected"
        m := newMemberWithInvite(famID, "John", &code)
        // invalid family name to trigger validation error
        fam := newFamilyWithMembers(famID, "owner1", "x", []family.Member{m})
        cmd := command.AcceptInvitationCommand{FamilyId: famID, FamilyMemberId: m.Id(), InvitationCode: code}

        repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
        auth.EXPECT().MustGetUserId(mock.Anything).Return("user-123")

        res := h.Handle(ctx, cmd)

        // We only assert that it's faulted and Save not called
        require.True(t, res.IsFaulted())
        repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
    })

    t.Run("Repository Save returns error", func(t *testing.T) {
        repo := ports.NewMockFamilyRepository(t)
        auth := ports.NewMockAuthentication(t)
        h := command.NewAcceptInvitationCommandHandler(repo, auth)

        famID := uuid.New()
        code := "expected"
        m := newMemberWithInvite(famID, "John", &code)
        fam := newFamilyWithMembers(famID, "owner1", "Doe family", []family.Member{m})
        cmd := command.AcceptInvitationCommand{FamilyId: famID, FamilyMemberId: m.Id(), InvitationCode: code}

        repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
        auth.EXPECT().MustGetUserId(mock.Anything).Return("user-123")
        repo.EXPECT().Save(mock.Anything, mock.Anything).Return(errors.New("save failed"))

        res := h.Handle(ctx, cmd)

        require.True(t, res.IsFaulted())
    })

    t.Run("Successful acceptance -> success true, userId set, save called", func(t *testing.T) {
        repo := ports.NewMockFamilyRepository(t)
        auth := ports.NewMockAuthentication(t)
        h := command.NewAcceptInvitationCommandHandler(repo, auth)

        famID := uuid.New()
        code := "expected"
        m := newMemberWithInvite(famID, "John", &code)
        fam := newFamilyWithMembers(famID, "owner1", "Doe family", []family.Member{m})
        cmd := command.AcceptInvitationCommand{FamilyId: famID, FamilyMemberId: m.Id(), InvitationCode: code}

        repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
        auth.EXPECT().MustGetUserId(mock.Anything).Return("user-123")
        repo.EXPECT().Save(mock.Anything, mock.Anything).Return(nil)

        res := h.Handle(ctx, cmd)

        require.True(t, res.IsSuccess())
        // Ensure the member now has the user id set
        updated := fam.GetMember(m.Id())
        require.NotNil(t, updated)
        require.NotNil(t, updated.UserId())
        assert.Equal(t, "user-123", *updated.UserId())
    })
}
