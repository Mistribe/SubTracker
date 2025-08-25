package command

import (
	"context"

	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/domain/auth"
	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type AcceptInvitationCommand struct {
	InvitationCode string
	FamilyId       uuid.UUID
	FamilyMemberId uuid.UUID
}

type AcceptInvitationCommandHandler struct {
	familyRepository family.Repository
	authService      auth.Service
}

func NewAcceptInvitationCommandHandler(
	familyRepository family.Repository,
	authService auth.Service) *AcceptInvitationCommandHandler {
	return &AcceptInvitationCommandHandler{familyRepository: familyRepository, authService: authService}
}

func (h AcceptInvitationCommandHandler) Handle(ctx context.Context, cmd AcceptInvitationCommand) result.Result[bool] {
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

	userId := h.authService.MustGetUserId(ctx)
	member.SetUserId(&userId)

	if err = fam.UpdateMember(member); err != nil {
		return result.Fail[bool](err)
	}
	if err = fam.GetValidationErrors(); err != nil {
		return result.Fail[bool](err)
	}
	if err = h.familyRepository.Save(ctx, fam); err != nil {
		return result.Fail[bool](err)
	}

	return result.Success(true)
}
