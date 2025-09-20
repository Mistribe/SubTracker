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

func newMemberWithInviteDecline(famID uuid.UUID, name string, code *string) family.Member {
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

func newFamilyWithMembersDecline(id uuid.UUID, ownerID string, name string, members []family.Member) family.Family {
    return family.NewFamily(
        id,
        ownerID,
        name,
        members,
        time.Now(),
        time.Now(),
    )
}

func TestDeclineInvitationCommandHandler(t *testing.T) {
    ctx := context.Background()

    t.Run("GetById returns error", func(t *testing.T) {
        repo := ports.NewMockFamilyRepository(t)
        h := command.NewDeclineInvitationCommandHandler(repo)

        famID := uuid.New()
        cmd := command.DeclineInvitationCommand{FamilyId: famID}

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
        h := command.NewDeclineInvitationCommandHandler(repo)

        famID := uuid.New()
        cmd := command.DeclineInvitationCommand{FamilyId: famID}

        repo.EXPECT().GetById(mock.Anything, famID).Return(nil, nil)

        res := h.Handle(ctx, cmd)

        receivedErr := result.Match(res, func(_ bool) error { return nil }, func(err error) error { return err })
        require.Error(t, receivedErr)
        assert.True(t, errors.Is(receivedErr, family.ErrFamilyNotFound))
        repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
    })

    t.Run("Member not found returns ErrFamilyMemberNotFound", func(t *testing.T) {
        repo := ports.NewMockFamilyRepository(t)
        h := command.NewDeclineInvitationCommandHandler(repo)

        famID := uuid.New()
        otherMember := newMemberWithInviteDecline(famID, "John", nil)
        fam := newFamilyWithMembersDecline(famID, "owner1", "Doe family", []family.Member{otherMember})
        cmd := command.DeclineInvitationCommand{FamilyId: famID, FamilyMemberId: uuid.New()}

        repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)

        res := h.Handle(ctx, cmd)

        receivedErr := result.Match(res, func(_ bool) error { return nil }, func(err error) error { return err })
        require.Error(t, receivedErr)
        assert.True(t, errors.Is(receivedErr, family.ErrFamilyMemberNotFound))
        repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
    })

    t.Run("Invitation code is nil returns ErrBadInvitationCode", func(t *testing.T) {
        repo := ports.NewMockFamilyRepository(t)
        h := command.NewDeclineInvitationCommandHandler(repo)

        famID := uuid.New()
        m := newMemberWithInviteDecline(famID, "John", nil)
        fam := newFamilyWithMembersDecline(famID, "owner1", "Doe family", []family.Member{m})
        cmd := command.DeclineInvitationCommand{FamilyId: famID, FamilyMemberId: m.Id(), InvitationCode: "expected"}

        repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)

        res := h.Handle(ctx, cmd)

        receivedErr := result.Match(res, func(_ bool) error { return nil }, func(err error) error { return err })
        require.Error(t, receivedErr)
        assert.True(t, errors.Is(receivedErr, family.ErrBadInvitationCode))
        repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
    })

    t.Run("Invitation code mismatch returns ErrBadInvitationCode", func(t *testing.T) {
        repo := ports.NewMockFamilyRepository(t)
        h := command.NewDeclineInvitationCommandHandler(repo)

        famID := uuid.New()
        code := "expected"
        m := newMemberWithInviteDecline(famID, "John", &code)
        fam := newFamilyWithMembersDecline(famID, "owner1", "Doe family", []family.Member{m})
        cmd := command.DeclineInvitationCommand{FamilyId: famID, FamilyMemberId: m.Id(), InvitationCode: "wrong"}

        repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)

        res := h.Handle(ctx, cmd)

        receivedErr := result.Match(res, func(_ bool) error { return nil }, func(err error) error { return err })
        require.Error(t, receivedErr)
        assert.True(t, errors.Is(receivedErr, family.ErrBadInvitationCode))
        repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
    })

    t.Run("Validation error after update -> fail and no save", func(t *testing.T) {
        repo := ports.NewMockFamilyRepository(t)
        h := command.NewDeclineInvitationCommandHandler(repo)

        famID := uuid.New()
        code := "expected"
        m := newMemberWithInviteDecline(famID, "John", &code)
        // invalid family name to trigger validation error
        fam := newFamilyWithMembersDecline(famID, "owner1", "x", []family.Member{m})
        cmd := command.DeclineInvitationCommand{FamilyId: famID, FamilyMemberId: m.Id(), InvitationCode: code}

        repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)

        res := h.Handle(ctx, cmd)

        // We only assert that it's faulted and Save not called
        require.True(t, res.IsFaulted())
        repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
    })

    t.Run("Repository Save returns error", func(t *testing.T) {
        repo := ports.NewMockFamilyRepository(t)
        h := command.NewDeclineInvitationCommandHandler(repo)

        famID := uuid.New()
        code := "expected"
        m := newMemberWithInviteDecline(famID, "John", &code)
        fam := newFamilyWithMembersDecline(famID, "owner1", "Doe family", []family.Member{m})
        cmd := command.DeclineInvitationCommand{FamilyId: famID, FamilyMemberId: m.Id(), InvitationCode: code}

        repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
        repo.EXPECT().Save(mock.Anything, mock.Anything).Return(errors.New("save failed"))

        res := h.Handle(ctx, cmd)

        require.True(t, res.IsFaulted())
    })

    t.Run("Successful decline -> success true, invitation cleared, save called", func(t *testing.T) {
        repo := ports.NewMockFamilyRepository(t)
        h := command.NewDeclineInvitationCommandHandler(repo)

        famID := uuid.New()
        code := "expected"
        m := newMemberWithInviteDecline(famID, "John", &code)
        fam := newFamilyWithMembersDecline(famID, "owner1", "Doe family", []family.Member{m})
        cmd := command.DeclineInvitationCommand{FamilyId: famID, FamilyMemberId: m.Id(), InvitationCode: code}

        repo.EXPECT().GetById(mock.Anything, famID).Return(fam, nil)
        repo.EXPECT().Save(mock.Anything, mock.Anything).Return(nil)

        res := h.Handle(ctx, cmd)

        require.True(t, res.IsSuccess())
        // Ensure the member now has no invitation code
        updated := fam.GetMember(m.Id())
        require.NotNil(t, updated)
        assert.Nil(t, updated.InvitationCode())
    })
}
