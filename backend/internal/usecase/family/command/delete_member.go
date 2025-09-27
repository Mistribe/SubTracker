package command

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/authorization"
	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type DeleteFamilyMemberCommand struct {
	FamilyMemberID types.FamilyMemberID
	FamilyID       types.FamilyID
}

type DeleteFamilyMemberCommandHandler struct {
	familyRepository ports.FamilyRepository
	authorization    ports.Authorization
}

func NewDeleteFamilyMemberCommandHandler(
	familyRepository ports.FamilyRepository,
	authorization ports.Authorization) *DeleteFamilyMemberCommandHandler {
	return &DeleteFamilyMemberCommandHandler{
		familyRepository: familyRepository,
		authorization:    authorization,
	}
}

func (h DeleteFamilyMemberCommandHandler) Handle(
	ctx context.Context,
	command DeleteFamilyMemberCommand) result.Result[bool] {
	fam, err := h.familyRepository.GetById(ctx, command.FamilyID)
	if err != nil {
		return result.Fail[bool](err)
	}

	if fam == nil {
		return result.Fail[bool](family.ErrFamilyNotFound)
	}

	if err = h.authorization.Can(ctx, authorization.PermissionDelete).For(fam); err != nil {
		return result.Fail[bool](err)
	}
	return h.deleteMember(ctx, command, fam)
}

func (h DeleteFamilyMemberCommandHandler) deleteMember(
	ctx context.Context, command DeleteFamilyMemberCommand,
	fam family.Family) result.Result[bool] {
	mbr := fam.GetMember(command.FamilyMemberID)
	if mbr == nil {
		return result.Fail[bool](family.ErrFamilyMemberNotFound)
	}

	if err := fam.RemoveMember(mbr); err != nil {
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
