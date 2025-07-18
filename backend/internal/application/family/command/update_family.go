package command

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/oleexo/subtracker/internal/domain/family"
	"github.com/oleexo/subtracker/pkg/langext/option"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type UpdateFamilyCommand struct {
	Id               uuid.UUID
	Name             string
	UpdatedAt        option.Option[time.Time]
	HaveJointAccount bool
}

type UpdateFamilyCommandHandler interface {
	Handle(ctx context.Context, cmd UpdateFamilyCommand) result.Result[family.Family]
}

type updateFamilyCommandHandler struct {
	repository family.Repository
}

func NewUpdateFamilyCommandHandler(repository family.Repository) UpdateFamilyCommandHandler {
	return &updateFamilyCommandHandler{
		repository: repository,
	}
}

func (h *updateFamilyCommandHandler) Handle(ctx context.Context, cmd UpdateFamilyCommand) result.Result[family.Family] {
	famOpt, err := h.repository.GetById(ctx, cmd.Id)
	if err != nil {
		return result.Fail[family.Family](err)
	}

	return option.Match(famOpt, func(fam family.Family) result.Result[family.Family] {
		return h.updateFamily(ctx, cmd, fam)
	}, func() result.Result[family.Family] {
		return result.Fail[family.Family](family.ErrFamilyNotFound)
	})

}

func (h *updateFamilyCommandHandler) updateFamily(ctx context.Context, cmd UpdateFamilyCommand, fam family.Family) result.Result[family.Family] {
	if err := ensureOwnerIsEditor(ctx, fam.OwnerId()); err != nil {
		return result.Fail[family.Family](err)
	}
	fam.SetName(cmd.Name)
	fam.SetHaveJointAccount(cmd.HaveJointAccount)

	cmd.UpdatedAt.IfSome(func(updatedAt time.Time) {
		fam.SetUpdatedAt(updatedAt)
	})
	cmd.UpdatedAt.IfNone(func() {
		fam.SetUpdatedAt(time.Now())
	})

	if err := fam.Validate(); err != nil {
		return result.Fail[family.Family](err)
	}

	if err := h.repository.Save(ctx, &fam); err != nil {
		return result.Fail[family.Family](err)
	}

	return result.Success(fam)
}
