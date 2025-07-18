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
	Email     option.Option[string]
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
	famOpt, err := h.repository.GetById(ctx, command.FamilyId)
	if err != nil {
		return result.Fail[family.Family](err)
	}

	return option.Match(famOpt, func(fam family.Family) result.Result[family.Family] {
		return option.Match(fam.GetMember(command.Id),
			func(mbr family.Member) result.Result[family.Family] {
				return h.updateFamilyMember(ctx, command, fam, mbr)
			}, func() result.Result[family.Family] {
				return result.Fail[family.Family](family.ErrFamilyMemberNotFound)
			})

	}, func() result.Result[family.Family] {
		return result.Fail[family.Family](family.ErrFamilyNotFound)
	})
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
	mbr.SetEmail(command.Email)

	command.UpdatedAt.IfSome(func(updatedAt time.Time) {
		mbr.SetUpdatedAt(updatedAt)
	})
	command.UpdatedAt.IfNone(func() {
		mbr.SetUpdatedAt(time.Now())
	})

	if err := mbr.Validate(); err != nil {
		return result.Fail[family.Family](err)
	}

	if err := fam.UpdateMember(mbr); err != nil {
		return result.Fail[family.Family](err)
	}

	if err := fam.Validate(); err != nil {
		return result.Fail[family.Family](err)
	}

	if err := h.repository.SaveMember(ctx, &mbr); err != nil {
		return result.Fail[family.Family](err)
	}

	return result.Success(fam)
}
