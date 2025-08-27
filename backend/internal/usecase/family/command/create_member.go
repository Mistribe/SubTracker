package command

import (
	"context"

	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type CreateFamilyMemberCommand struct {
	FamilyId uuid.UUID
	Member   family.Member
}

type CreateFamilyMemberCommandHandler struct {
	repository ports.FamilyRepository
}

func NewCreateFamilyMemberCommandHandler(repository ports.FamilyRepository) *CreateFamilyMemberCommandHandler {
	return &CreateFamilyMemberCommandHandler{repository: repository}
}

func (h CreateFamilyMemberCommandHandler) Handle(
	ctx context.Context,
	command CreateFamilyMemberCommand) result.Result[family.Family] {
	fam, err := h.repository.GetById(ctx, command.FamilyId)
	if err != nil {
		return result.Fail[family.Family](err)
	}
	if fam == nil {
		return result.Fail[family.Family](family.ErrFamilyNotFound)

	}

	return h.addFamilyMemberToFamily(ctx, command, fam)
}

func (h CreateFamilyMemberCommandHandler) addFamilyMemberToFamily(
	ctx context.Context,
	command CreateFamilyMemberCommand, fam family.Family) result.Result[family.Family] {
	if err := ensureOwnerIsEditor(ctx, fam.OwnerId()); err != nil {
		return result.Fail[family.Family](err)
	}
	if err := fam.AddMember(command.Member); err != nil {
		return result.Success(fam)
	}

	if err := fam.GetValidationErrors(); err != nil {
		return result.Fail[family.Family](err)
	}

	if err := h.repository.Save(ctx, fam); err != nil {
		return result.Fail[family.Family](err)
	}

	return result.Success(fam)

}
