package command

import (
	"context"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/pkg/langext/option"

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
	famOpt, err := h.repository.GetById(ctx, command.FamilyId)
	if err != nil {
		return result.Fail[bool](err)
	}

	return option.Match(famOpt, func(fam family.Family) result.Result[bool] {
		return h.deleteMember(ctx, command, fam)
	}, func() result.Result[bool] {
		return result.Fail[bool](family.ErrFamilyNotFound)
	})
}

func (h DeleteFamilyMemberCommandHandler) deleteMember(
	ctx context.Context, command DeleteFamilyMemberCommand,
	fam family.Family) result.Result[bool] {
	if err := ensureOwnerIsEditor(ctx, fam.OwnerId()); err != nil {
		return result.Fail[bool](err)
	}
	if !fam.ContainsMember(command.Id) {
		return result.Fail[bool](family.ErrFamilyMemberNotFound)
	}

	if err := fam.GetValidationErrors(); err != nil {
		return result.Fail[bool](err)
	}

	if err := h.repository.Save(ctx, fam); err != nil {
		return result.Fail[bool](err)
	}

	return result.Success(true)
}
