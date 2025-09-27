package command

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type AcceptInvitationCommand struct {
	InvitationCode string
	FamilyId       types.FamilyID
	FamilyMemberId types.FamilyMemberID
}

type AcceptInvitationCommandHandler struct {
	familyRepository ports.FamilyRepository
	authService      ports.Authentication
}

func NewAcceptInvitationCommandHandler(
	familyRepository ports.FamilyRepository,
	authService ports.Authentication) *AcceptInvitationCommandHandler {
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

	connectedAccount := h.authService.MustGetConnectedAccount(ctx)
	userID := connectedAccount.UserID()
	member.SetUserId(&userID)

	if err = fam.UpdateMember(member); err != nil {
		return result.Fail[bool](err)
	}
	if vErr := fam.GetValidationErrors(); vErr != nil {
		return result.Fail[bool](err)
	}
	if err = h.familyRepository.Save(ctx, fam); err != nil {
		return result.Fail[bool](err)
	}

	return result.Success(true)
}
