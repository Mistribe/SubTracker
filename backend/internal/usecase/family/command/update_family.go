package command

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/langext/option"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type UpdateFamilyCommand struct {
	Id        uuid.UUID
	Name      string
	UpdatedAt option.Option[time.Time]
}

type UpdateFamilyCommandHandler interface {
	Handle(ctx context.Context, cmd UpdateFamilyCommand) result.Result[family.Family]
}

type updateFamilyCommandHandler struct {
	repository ports.FamilyRepository
}

func NewUpdateFamilyCommandHandler(repository ports.FamilyRepository) UpdateFamilyCommandHandler {
	return &updateFamilyCommandHandler{
		repository: repository,
	}
}

func (h *updateFamilyCommandHandler) Handle(ctx context.Context, cmd UpdateFamilyCommand) result.Result[family.Family] {
	fam, err := h.repository.GetById(ctx, cmd.Id)
	if err != nil {
		return result.Fail[family.Family](err)
	}
	if fam == nil {
		return result.Fail[family.Family](family.ErrFamilyNotFound)

	}
	return h.updateFamily(ctx, cmd, fam)

}

func (h *updateFamilyCommandHandler) updateFamily(
	ctx context.Context, cmd UpdateFamilyCommand,
	fam family.Family) result.Result[family.Family] {
	if err := ensureOwnerIsEditor(ctx, fam.OwnerId()); err != nil {
		return result.Fail[family.Family](err)
	}
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

	if err := h.repository.Save(ctx, fam); err != nil {
		return result.Fail[family.Family](err)
	}

	return result.Success(fam)
}
