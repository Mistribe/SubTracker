package command

import (
	"context"
	"time"

	"github.com/mistribe/subtracker/internal/domain/authorization"
	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type InviteMemberCommand struct {
	Email          *string
	Name           *string
	Type           *family.MemberType
	FamilyId       types.FamilyID
	FamilyMemberId types.FamilyMemberID
}

type InviteMemberResponse struct {
	Code           string
	FamilyId       types.FamilyID
	FamilyMemberId types.FamilyMemberID
}

type InviteMemberCommandHandler struct {
	familyRepository ports.FamilyRepository
	authorization    ports.Authorization
}

func NewInviteMemberCommandHandler(
	familyRepository ports.FamilyRepository,
	authorization ports.Authorization) *InviteMemberCommandHandler {
	return &InviteMemberCommandHandler{
		familyRepository: familyRepository,
		authorization:    authorization,
	}
}

func (h InviteMemberCommandHandler) Handle(
	ctx context.Context,
	cmd InviteMemberCommand) result.Result[InviteMemberResponse] {
	fam, err := h.familyRepository.GetById(ctx, cmd.FamilyId)
	if err != nil {
		return result.Fail[InviteMemberResponse](err)
	}
	if fam == nil {
		return result.Fail[InviteMemberResponse](family.ErrFamilyNotFound)
	}
	if err = h.authorization.Can(ctx, authorization.PermissionWrite).For(fam); err != nil {
		return result.Fail[InviteMemberResponse](err)
	}

	code, err := family.NewGenerateInvitationCode(fam.Id())
	if err != nil {
		return result.Fail[InviteMemberResponse](err)
	}

	familyMember := fam.GetMember(cmd.FamilyMemberId)
	if familyMember == nil {
		return h.inviteNewMember(ctx, cmd, familyMember, fam, code)

	}
	return h.inviteExistingMember(ctx, familyMember, code, fam)
}

func (h InviteMemberCommandHandler) inviteExistingMember(
	ctx context.Context,
	familyMember family.Member,
	code string,
	fam family.Family) result.Result[InviteMemberResponse] {
	if familyMember.UserId() != nil {
		return result.Fail[InviteMemberResponse](family.ErrCannotInviteUser)
	}
	if !familyMember.SetInvitationCode(&code) {
		return result.Fail[InviteMemberResponse](family.ErrCannotInviteUser)
	}

	if err := fam.UpdateMember(familyMember); err != nil {
		return result.Fail[InviteMemberResponse](err)
	}

	if err := fam.GetValidationErrors(); err != nil {
		return result.Fail[InviteMemberResponse](err)
	}

	if err := h.familyRepository.Save(ctx, fam); err != nil {
		return result.Fail[InviteMemberResponse](err)
	}

	return result.Success(InviteMemberResponse{
		Code:           code,
		FamilyId:       fam.Id(),
		FamilyMemberId: familyMember.Id(),
	})
}

func (h InviteMemberCommandHandler) inviteNewMember(
	ctx context.Context,
	cmd InviteMemberCommand,
	familyMember family.Member,
	fam family.Family,
	invitationCode string) result.Result[InviteMemberResponse] {
	var memberName string
	if cmd.Name != nil {
		memberName = *cmd.Name
	} else {
		memberName = "Invited User"
	}
	var memberType family.MemberType
	if cmd.Type != nil {
		memberType = *cmd.Type
	} else {
		memberType = family.AdultMemberType
	}
	familyMember = family.NewMember(
		types.NewFamilyMemberID(),
		fam.Id(),
		memberName,
		memberType,
		&invitationCode,
		time.Now(),
		time.Now(),
	)
	if err := fam.AddMember(familyMember); err != nil {
		return result.Fail[InviteMemberResponse](err)
	}

	if err := fam.GetValidationErrors(); err != nil {
		return result.Fail[InviteMemberResponse](err)
	}

	if err := h.familyRepository.Save(ctx, fam); err != nil {
		return result.Fail[InviteMemberResponse](err)
	}

	return result.Success(InviteMemberResponse{
		Code:           invitationCode,
		FamilyId:       fam.Id(),
		FamilyMemberId: familyMember.Id(),
	})
}
