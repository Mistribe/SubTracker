package command

import (
	"context"

	"github.com/oleexo/subtracker/internal/domain/family"
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
	command CreateFamilyMemberCommand) result.Result[family.Member] {
	existingMember, err := h.repository.Get(ctx, command.Member.Id())
	if err != nil {
		return result.Fail[family.Member](err)
	}
	if existingMember.IsSome() {
		return result.Fail[family.Member](family.ErrMemberAlreadyExists)
	}

	if err := command.Member.Validate(); err != nil {
		return result.Fail[family.Member](err)
	}
	// todo make a new member not from input
	if err := h.repository.Save(ctx, &command.Member); err != nil {
		return result.Fail[family.Member](err)
	}

	return result.Success(command.Member)
}
