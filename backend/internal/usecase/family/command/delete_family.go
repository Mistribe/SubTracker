package command

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/authorization"
	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type DeleteFamilyCommand struct {
	FamilyId types.FamilyID
}

type DeleteFamilyCommandHandler struct {
	familyRepository ports.FamilyRepository
	authorization    ports.Authorization
}

func NewDeleteFamilyCommandHandler(
	familyRepository ports.FamilyRepository,
	authorization ports.Authorization) *DeleteFamilyCommandHandler {
	return &DeleteFamilyCommandHandler{
		familyRepository: familyRepository,
		authorization:    authorization,
	}
}

func (h DeleteFamilyCommandHandler) Handle(
	ctx context.Context,
	command DeleteFamilyCommand) result.Result[bool] {
	fam, err := h.familyRepository.GetById(ctx, command.FamilyId)
	if err != nil {
		return result.Fail[bool](err)
	}
	if fam == nil {
		return result.Fail[bool](family.ErrFamilyNotFound)
	}
	if err = h.authorization.Can(ctx, authorization.PermissionDelete).For(fam); err != nil {
		return result.Fail[bool](err)
	}
	ok, err := h.familyRepository.Delete(ctx, command.FamilyId)
	if err != nil {
		return result.Fail[bool](err)
	}
	if !ok {
		return result.Fail[bool](family.ErrFamilyNotFound)
	}

	return result.Success(true)
}
