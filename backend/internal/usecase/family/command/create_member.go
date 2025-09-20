package command

import (
	"context"

	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/domain/user"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type CreateFamilyMemberCommand struct {
	FamilyId uuid.UUID
	Member   family.Member
}

type CreateFamilyMemberCommandHandler struct {
	familyRepository ports.FamilyRepository
	authorization    ports.Authorization
}

func NewCreateFamilyMemberCommandHandler(
	familyRepository ports.FamilyRepository,
	authorization ports.Authorization) *CreateFamilyMemberCommandHandler {
	return &CreateFamilyMemberCommandHandler{
		familyRepository: familyRepository,
		authorization:    authorization,
	}
}

func (h CreateFamilyMemberCommandHandler) Handle(
	ctx context.Context,
	command CreateFamilyMemberCommand) result.Result[family.Family] {
	fam, err := h.familyRepository.GetById(ctx, command.FamilyId)
	if err != nil {
		return result.Fail[family.Family](err)
	}
	if fam == nil {
		return result.Fail[family.Family](family.ErrFamilyNotFound)

	}

	if err = h.authorization.Can(ctx, user.PermissionWrite).For(fam); err != nil {
		return result.Fail[family.Family](err)
	}

	if err = h.authorization.EnsureWithinLimit(ctx, user.FeatureFamilyMembers, 1); err != nil {
		return result.Fail[family.Family](err)
	}

	return h.addFamilyMemberToFamily(ctx, command, fam)
}

func (h CreateFamilyMemberCommandHandler) addFamilyMemberToFamily(
	ctx context.Context,
	command CreateFamilyMemberCommand, fam family.Family) result.Result[family.Family] {

	if err := fam.AddMember(command.Member); err != nil {
		return result.Success(fam)
	}

	if err := fam.GetValidationErrors(); err != nil {
		return result.Fail[family.Family](err)
	}

	if err := h.familyRepository.Save(ctx, fam); err != nil {
		return result.Fail[family.Family](err)
	}

	return result.Success(fam)

}
