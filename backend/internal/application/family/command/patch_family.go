package command

import (
	"context"

	"github.com/oleexo/subtracker/internal/domain/family"
	"github.com/oleexo/subtracker/pkg/langext/option"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type PatchFamilyCommand struct {
	Family family.Family
}

type PatchFamilyCommandHandler struct {
	repository family.Repository
}

func NewPatchFamilyCommandHandler(familyRepository family.Repository) *PatchFamilyCommandHandler {
	return &PatchFamilyCommandHandler{repository: familyRepository}
}

func (h PatchFamilyCommandHandler) Handle(ctx context.Context, cmd PatchFamilyCommand) result.Result[family.Family] {
	famOpt, err := h.repository.GetById(ctx, cmd.Family.Id())
	if err != nil {
		return result.Fail[family.Family](err)
	}
	return option.Match(famOpt, func(fam family.Family) result.Result[family.Family] {
		return h.patchFamily(ctx, cmd, fam)
	}, func() result.Result[family.Family] {
		return h.createFamily(ctx, cmd)
	})
}

func (h PatchFamilyCommandHandler) createFamily(
	ctx context.Context,
	cmd PatchFamilyCommand) result.Result[family.Family] {
	if err := cmd.Family.GetValidationErrors(); err != nil {
		return result.Fail[family.Family](err)
	}
	if err := h.repository.Save(ctx, cmd.Family); err != nil {
		return result.Fail[family.Family](err)
	}

	return result.Success(cmd.Family)
}

func (h PatchFamilyCommandHandler) patchFamily(
	ctx context.Context,
	cmd PatchFamilyCommand,
	fam family.Family) result.Result[family.Family] {
	if err := ensureOwnerIsEditor(ctx, fam.OwnerId()); err != nil {
		return result.Fail[family.Family](err)
	}
	fam.SetName(cmd.Family.Name())
	fam.SetUpdatedAt(cmd.Family.UpdatedAt())

	for mbr := range fam.Members().It() {
		if !cmd.Family.Members().Contains(mbr) {
			fam.Members().Remove(mbr)
		}
	}

	for mbr := range cmd.Family.Members().It() {
		existingMember, ok := fam.Members().Get(mbr)
		if !ok {
			fam.Members().Add(mbr)
		} else {
			if mbr.IsKid() {
				existingMember.SetAsKid()
			} else {
				existingMember.SetAsAdult()
			}
			existingMember.SetName(mbr.Name())
			existingMember.SetEmail(mbr.Email())
			existingMember.SetUpdatedAt(mbr.UpdatedAt())

			fam.Members().Update(existingMember)
		}
	}

	if err := fam.GetValidationErrors(); err != nil {
		return result.Fail[family.Family](err)
	}

	if err := h.repository.Save(ctx, fam); err != nil {
		return result.Fail[family.Family](err)
	}

	return result.Success(fam)
}
