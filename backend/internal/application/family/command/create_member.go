package command

import (
	"context"

	"github.com/oleexo/subtracker/internal/domain/family"
	"github.com/oleexo/subtracker/pkg/langext/option"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type CreateFamilyMemberCommand struct {
	Member family.Member
}

type CreateFamilyMemberCommandHandler struct {
	repository family.Repository
}

func NewCreateFamilyMemberCommandHandler(repository family.Repository) *CreateFamilyMemberCommandHandler {
	return &CreateFamilyMemberCommandHandler{repository: repository}
}

func (h CreateFamilyMemberCommandHandler) Handle(
	ctx context.Context,
	command CreateFamilyMemberCommand) result.Result[family.Family] {
	if err := command.Member.Validate(); err != nil {
		return result.Fail[family.Family](err)
	}
	famOpt, err := h.repository.GetById(ctx, command.Member.FamilyId())
	if err != nil {
		return result.Fail[family.Family](err)
	}

	return option.Match(famOpt, func(fam family.Family) result.Result[family.Family] {
		return h.addFamilyMemberToFamily(ctx, command, fam)
	}, func() result.Result[family.Family] {
		return result.Fail[family.Family](family.ErrFamilyNotFound)
	})
}

func (h CreateFamilyMemberCommandHandler) addFamilyMemberToFamily(ctx context.Context, command CreateFamilyMemberCommand, fam family.Family) result.Result[family.Family] {
	if err := ensureOwnerIsEditor(ctx, fam.OwnerId()); err != nil {
		return result.Fail[family.Family](err)
	}
	if err := fam.AddMember(command.Member); err != nil {
		return result.Fail[family.Family](err)
	}

	if err := h.repository.SaveMember(ctx, &command.Member); err != nil {
		return result.Fail[family.Family](err)
	}

	return result.Success(fam)

}
