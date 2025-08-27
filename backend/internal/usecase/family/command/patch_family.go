package command

import (
	"context"

	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type PatchFamilyCommand struct {
	Family family.Family
}

type PatchFamilyCommandHandler struct {
	familyRepository ports.FamilyRepository
	authService      ports.AuthService
}

func NewPatchFamilyCommandHandler(
	familyRepository ports.FamilyRepository,
	authService ports.AuthService) *PatchFamilyCommandHandler {
	return &PatchFamilyCommandHandler{
		familyRepository: familyRepository,
		authService:      authService,
	}
}

func (h PatchFamilyCommandHandler) Handle(ctx context.Context, cmd PatchFamilyCommand) result.Result[family.Family] {
	fam, err := h.familyRepository.GetById(ctx, cmd.Family.Id())
	if err != nil {
		return result.Fail[family.Family](err)
	}
	if fam == nil {
		return h.createFamily(ctx, cmd)
	}
	return h.patchFamily(ctx, cmd, fam)
}

func (h PatchFamilyCommandHandler) createFamily(
	ctx context.Context,
	cmd PatchFamilyCommand) result.Result[family.Family] {
	userId := h.authService.MustGetUserId(ctx)
	memberId, err := uuid.NewV7()
	if err != nil {
		return result.Fail[family.Family](err)
	}

	creator := family.NewMember(
		memberId,
		cmd.Family.Id(),
		"You",
		family.OwnerMemberType,
		nil,
		cmd.Family.CreatedAt(),
		cmd.Family.UpdatedAt(),
	)
	creator.SetUserId(&userId)

	cmd.Family.AddMember(creator)

	if err := cmd.Family.GetValidationErrors(); err != nil {
		return result.Fail[family.Family](err)
	}
	if err := h.familyRepository.Save(ctx, cmd.Family); err != nil {
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
			existingMember.SetType(mbr.Type())
			existingMember.SetName(mbr.Name())
			existingMember.SetUpdatedAt(mbr.UpdatedAt())

			fam.Members().Update(existingMember)
		}
	}

	if err := fam.GetValidationErrors(); err != nil {
		return result.Fail[family.Family](err)
	}

	if err := h.familyRepository.Save(ctx, fam); err != nil {
		return result.Fail[family.Family](err)
	}

	return result.Success(fam)
}
