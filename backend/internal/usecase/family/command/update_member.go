package command

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/domain/user"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/langext/option"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type UpdateFamilyMemberCommand struct {
	Id        uuid.UUID
	FamilyId  uuid.UUID
	Name      string
	Type      family.MemberType
	UpdatedAt option.Option[time.Time]
}

type UpdateFamilyMemberCommandHandler struct {
	familyRepository ports.FamilyRepository
	authorization    ports.Authorization
}

func NewUpdateFamilyMemberCommandHandler(
	familyRepository ports.FamilyRepository,
	authorization ports.Authorization) *UpdateFamilyMemberCommandHandler {
	return &UpdateFamilyMemberCommandHandler{
		familyRepository: familyRepository,
		authorization:    authorization,
	}
}

func (h UpdateFamilyMemberCommandHandler) Handle(
	ctx context.Context,
	command UpdateFamilyMemberCommand) result.Result[family.Family] {
	fam, err := h.familyRepository.GetById(ctx, command.FamilyId)
	if err != nil {
		return result.Fail[family.Family](err)
	}

	if fam == nil {
		return result.Fail[family.Family](family.ErrFamilyNotFound)
	}

	mbr := fam.GetMember(command.Id)
	if mbr == nil {
		return result.Fail[family.Family](family.ErrFamilyMemberNotFound)
	}

	if err = h.authorization.Can(ctx, user.PermissionWrite).For(fam); err != nil {
		return result.Fail[family.Family](err)
	}

	return h.updateFamilyMember(ctx, command, fam, mbr)
}

func (h UpdateFamilyMemberCommandHandler) updateFamilyMember(
	ctx context.Context,
	cmd UpdateFamilyMemberCommand,
	fam family.Family,
	mbr family.Member) result.Result[family.Family] {
	mbr.SetName(cmd.Name)
	mbr.SetType(cmd.Type)

	cmd.UpdatedAt.IfSome(func(updatedAt time.Time) {
		mbr.SetUpdatedAt(updatedAt)
	})
	cmd.UpdatedAt.IfNone(func() {
		mbr.SetUpdatedAt(time.Now())
	})

	if err := fam.UpdateMember(mbr); err != nil {
		return result.Fail[family.Family](err)
	}

	if err := fam.GetValidationErrors(); err != nil {
		return result.Fail[family.Family](err)
	}

	if err := h.familyRepository.Save(ctx, fam); err != nil {
		return result.Fail[family.Family](err)
	}

	return result.Success(fam)
}
