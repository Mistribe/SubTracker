package command

import (
	"context"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/family"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type DeleteFamilyCommand struct {
	FamilyId uuid.UUID
}

type DeleteFamilyCommandHandler struct {
	repository family.Repository
}

func NewDeleteFamilyCommandHandler(repository family.Repository) *DeleteFamilyCommandHandler {
	return &DeleteFamilyCommandHandler{repository: repository}
}

func (h DeleteFamilyCommandHandler) Handle(
	ctx context.Context,
	command DeleteFamilyCommand) result.Result[bool] {
	fam, err := h.repository.GetById(ctx, command.FamilyId)
	if err != nil {
		return result.Fail[bool](err)
	}
	if err := ensureOwnerIsEditor(ctx, fam.OwnerId()); err != nil {
		return result.Fail[bool](err)
	}
	ok, err := h.repository.Delete(ctx, command.FamilyId)
	if err != nil {
		return result.Fail[bool](err)
	}
	if !ok {
		return result.Fail[bool](family.ErrFamilyNotFound)
	}

	return result.Success(true)
}
