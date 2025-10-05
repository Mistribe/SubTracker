package command

import (
	"context"
	"time"

	"github.com/mistribe/subtracker/internal/domain/authorization"
	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/langext/option"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type UpdateFamilyCommand struct {
	FamilyID  types.FamilyID
	Name      string
	UpdatedAt option.Option[time.Time]
}

type UpdateFamilyCommandHandler struct {
	familyRepository ports.FamilyRepository
	authorization    ports.Authorization
}

func NewUpdateFamilyCommandHandler(
	familyRepository ports.FamilyRepository,
	authorization ports.Authorization) *UpdateFamilyCommandHandler {
	return &UpdateFamilyCommandHandler{
		familyRepository: familyRepository,
		authorization:    authorization,
	}
}

func (h *UpdateFamilyCommandHandler) Handle(ctx context.Context, cmd UpdateFamilyCommand) result.Result[family.Family] {
	fam, err := h.familyRepository.GetById(ctx, cmd.FamilyID)
	if err != nil {
		return result.Fail[family.Family](err)
	}
	if fam == nil {
		return result.Fail[family.Family](family.ErrFamilyNotFound)

	}

	if err = h.authorization.Can(ctx, authorization.PermissionWrite).For(fam); err != nil {
		return result.Fail[family.Family](err)
	}
	return h.updateFamily(ctx, cmd, fam)

}

func (h *UpdateFamilyCommandHandler) updateFamily(
	ctx context.Context, cmd UpdateFamilyCommand,
	fam family.Family) result.Result[family.Family] {
	fam.SetName(cmd.Name)

	cmd.UpdatedAt.IfSome(func(updatedAt time.Time) {
		fam.SetUpdatedAt(updatedAt)
	})
	cmd.UpdatedAt.IfNone(func() {
		fam.SetUpdatedAt(time.Now())
	})

	if err := fam.GetValidationErrors(); err != nil {
		return result.Fail[family.Family](err)
	}

	if err := h.familyRepository.Save(ctx, fam); err != nil {
		return result.Fail[family.Family](err)
	}

	return result.Success(fam)
}
