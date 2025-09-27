package command

import (
	"context"

	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/domain/authorization"
	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type RevokeMemberCommand struct {
	FamilyId       uuid.UUID
	FamilyMemberId uuid.UUID
}

type RevokeMemberCommandHandler struct {
	familyRepository ports.FamilyRepository
	authorization    ports.Authorization
}

func NewRevokeMemberCommandHandler(
	familyRepository ports.FamilyRepository,
	authorization ports.Authorization) *RevokeMemberCommandHandler {
	return &RevokeMemberCommandHandler{
		familyRepository: familyRepository,
		authorization:    authorization,
	}
}

func (h RevokeMemberCommandHandler) Handle(ctx context.Context, cmd RevokeMemberCommand) result.Result[bool] {
	fam, err := h.familyRepository.GetById(ctx, cmd.FamilyId)
	if err != nil {
		return result.Fail[bool](err)
	}

	if fam == nil {
		return result.Fail[bool](family.ErrFamilyNotFound)
	}

	if err = h.authorization.Can(ctx, authorization.PermissionWrite).For(fam); err != nil {
		return result.Fail[bool](err)
	}

	member := fam.GetMember(cmd.FamilyMemberId)
	if member == nil {
		return result.Fail[bool](family.ErrFamilyMemberNotFound)
	}

	return h.revokeMember(ctx, member, fam)
}

func (h RevokeMemberCommandHandler) revokeMember(
	ctx context.Context,
	member family.Member,
	fam family.Family) result.Result[bool] {
	if member.UserId() == nil {
		return result.Success(true)
	}

	member.SetUserId(nil)
	if err := fam.UpdateMember(member); err != nil {
		return result.Fail[bool](err)
	}

	if err := fam.GetValidationErrors(); err != nil {
		return result.Fail[bool](err)
	}

	if err := h.familyRepository.Save(ctx, fam); err != nil {
		return result.Fail[bool](err)
	}

	return result.Success(true)
}
