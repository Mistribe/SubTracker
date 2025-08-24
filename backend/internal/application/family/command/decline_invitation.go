package command

import (
	"context"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/family"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type DeclineInvitationCommand struct {
	InvitationCode string
	FamilyId       uuid.UUID
	FamilyMemberId uuid.UUID
}

type DeclineInvitationCommandHandler struct {
	familyRepository family.Repository
}

func NewDeclineInvitationCommandHandler(familyRepository family.Repository) *DeclineInvitationCommandHandler {
	return &DeclineInvitationCommandHandler{familyRepository: familyRepository}
}

func (h DeclineInvitationCommandHandler) Handle(ctx context.Context, cmd DeclineInvitationCommand) result.Result[bool] {
	fam, err := h.familyRepository.GetById(ctx, cmd.FamilyId)
	if err != nil {
		return result.Fail[bool](err)
	}

	if fam == nil {
		return result.Fail[bool](family.ErrFamilyNotFound)
	}

	member := fam.GetMember(cmd.FamilyMemberId)
	if member == nil {
		return result.Fail[bool](family.ErrFamilyMemberNotFound)
	}

	if member.InvitationCode() == nil {
		return result.Fail[bool](family.ErrBadInvitationCode)
	}

	if *member.InvitationCode() != cmd.InvitationCode {
		return result.Fail[bool](family.ErrBadInvitationCode)
	}

	return h.declineInvitation(member, fam)
}

func (h DeclineInvitationCommandHandler) declineInvitation(member family.Member,
	fam family.Family) result.Result[bool] {
	member.SetInvitationCode(nil)

	if err := fam.UpdateMember(member); err != nil {
		return result.Fail[bool](err)
	}

	if err := fam.GetValidationErrors(); err != nil {
		return result.Fail[bool](err)
	}

	if err := h.familyRepository.Save(context.Background(), fam); err != nil {
		return result.Fail[bool](err)
	}

	return result.Success(true)
}
