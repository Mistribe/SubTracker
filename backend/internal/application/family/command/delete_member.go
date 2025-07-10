package command

import (
	"context"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/application/core/result"
	"github.com/oleexo/subtracker/internal/domain/family"
)

type DeleteFamilyMemberCommand struct {
	Id uuid.UUID
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
	existingMember, err := h.repository.Get(ctx, command.Id)
	if err != nil {
		return result.Fail[bool](err)
	}
	if existingMember.IsNone() {
		return result.Fail[bool](family.ErrFamilyMemberNotFound)
	}

	if err := h.repository.Delete(ctx, command.Id); err != nil {
		return result.Fail[bool](err)
	}

	return result.Success(true)
}
