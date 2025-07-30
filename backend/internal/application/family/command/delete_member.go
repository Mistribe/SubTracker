package command

import (
	"context"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/family"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type DeleteFamilyMemberCommand struct {
	Id       uuid.UUID
	FamilyId uuid.UUID
}

type DeleteFamilyMemberCommandHandler struct {
	repository family.Repository
}

func NewDeleteFamilyMemberCommandHandler(repository family.Repository) *DeleteFamilyMemberCommandHandler {
	return &DeleteFamilyMemberCommandHandler{repository: repository}
}

func (h DeleteFamilyMemberCommandHandler) Handle(
	ctx context.Context,
	command DeleteFamilyMemberCommand) result.Result[bool] {
	fam, err := h.repository.GetById(ctx, command.FamilyId)
	if err != nil {
		return result.Fail[bool](err)
	}

	if fam == nil {
		return result.Fail[bool](family.ErrFamilyNotFound)
	}
	return h.deleteMember(ctx, command, fam)
}

func (h DeleteFamilyMemberCommandHandler) deleteMember(
	ctx context.Context, command DeleteFamilyMemberCommand,
	fam family.Family) result.Result[bool] {
	if err := ensureOwnerIsEditor(ctx, fam.OwnerId()); err != nil {
		return result.Fail[bool](err)
	}
	mbr := fam.GetMember(command.Id)
	if mbr == nil {
		return result.Fail[bool](family.ErrFamilyMemberNotFound)
	}

	if err := fam.RemoveMember(mbr); err != nil {
		return result.Fail[bool](err)
	}

	if err := fam.GetValidationErrors(); err != nil {
		return result.Fail[bool](err)
	}

	if err := h.repository.Save(ctx, fam); err != nil {
		return result.Fail[bool](err)
	}

	return result.Success(true)
}
