package command

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/family"
	"github.com/oleexo/subtracker/pkg/langext/option"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type UpdateFamilyMemberCommand struct {
	Id        uuid.UUID
	FamilyId  uuid.UUID
	Name      string
	IsKid     bool
	UpdatedAt option.Option[time.Time]
}

type UpdateFamilyMemberCommandHandler struct {
	repository family.Repository
}

func NewUpdateFamilyMemberCommandHandler(repository family.Repository) *UpdateFamilyMemberCommandHandler {
	return &UpdateFamilyMemberCommandHandler{repository: repository}
}

func (h UpdateFamilyMemberCommandHandler) Handle(
	ctx context.Context,
	command UpdateFamilyMemberCommand) result.Result[family.Family] {
	fam, err := h.repository.GetById(ctx, command.FamilyId)
	if err != nil {
		return result.Fail[family.Family](err)
	}

	if fam == nil {
		return result.Fail[family.Family](family.ErrFamilyNotFound)

	}

	mbr := fam.GetMember(command.Id)
	if mbr == nil {
		result.Fail[family.Family](family.ErrFamilyMemberNotFound)
	}
	return h.updateFamilyMember(ctx, command, fam, mbr)
}

func (h UpdateFamilyMemberCommandHandler) updateFamilyMember(
	ctx context.Context,
	command UpdateFamilyMemberCommand,
	fam family.Family,
	mbr family.Member) result.Result[family.Family] {
	if err := ensureOwnerIsEditor(ctx, fam.OwnerId()); err != nil {
		return result.Fail[family.Family](err)
	}
	mbr.SetName(command.Name)
	if command.IsKid {
		mbr.SetAsKid()
	} else {
		mbr.SetAsAdult()
	}

	command.UpdatedAt.IfSome(func(updatedAt time.Time) {
		mbr.SetUpdatedAt(updatedAt)
	})
	command.UpdatedAt.IfNone(func() {
		mbr.SetUpdatedAt(time.Now())
	})

	if err := fam.UpdateMember(mbr); err != nil {
		return result.Fail[family.Family](err)
	}

	if err := fam.GetValidationErrors(); err != nil {
		return result.Fail[family.Family](err)
	}

	if err := h.repository.Save(ctx, fam); err != nil {
		return result.Fail[family.Family](err)
	}

	return result.Success(fam)
}
