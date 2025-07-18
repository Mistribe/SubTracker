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
	command DeleteFamilyMemberCommand) result.Result[result.Unit] {
	famOpt, err := h.repository.GetById(ctx, command.FamilyId)
	if err != nil {
		return result.Fail[result.Unit](err)
	}

	return option.Match(famOpt, func(fam family.Family) result.Result[result.Unit] {
		return h.deleteMember(ctx, command, fam)
	}, func() result.Result[result.Unit] {
		return result.Fail[result.Unit](family.ErrFamilyNotFound)
	})
}

func (h DeleteFamilyMemberCommandHandler) deleteMember(ctx context.Context, command DeleteFamilyMemberCommand, fam family.Family) result.Result[result.Unit] {
	if err := ensureOwnerIsEditor(ctx, fam.OwnerId()); err != nil {
		return result.Fail[result.Unit](err)
	}
	if !fam.ContainsMember(command.Id) {
		return result.Fail[result.Unit](family.ErrFamilyMemberNotFound)
	}

	if err := h.repository.DeleteMember(ctx, command.Id); err != nil {
		return result.Fail[result.Unit](err)
	}

	return result.Void()
}
