package command

import (
	"context"

	"github.com/oleexo/subtracker/internal/domain/family"
	"github.com/oleexo/subtracker/pkg/langext/option"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type CreateFamilyCommand struct {
	Family family.Family
}

type CreateFamilyCommandHandler struct {
	repository family.Repository
}

func NewCreateFamilyCommandHandler(repository family.Repository) *CreateFamilyCommandHandler {
	return &CreateFamilyCommandHandler{
		repository: repository,
	}
}

func (h CreateFamilyCommandHandler) Handle(ctx context.Context, cmd CreateFamilyCommand) result.Result[family.Family] {
	if err := cmd.Family.Validate(); err != nil {
		return result.Fail[family.Family](err)
	}
	famOpt, err := h.repository.GetById(ctx, cmd.Family.Id())
	if err != nil {
		return result.Fail[family.Family](err)
	}
	return option.Match(famOpt, func(fam family.Family) result.Result[family.Family] {
		if fam.Equal(cmd.Family) {
			return result.Success(cmd.Family)
		}
		return result.Fail[family.Family](family.ErrFamilyAlreadyExists)
	}, func() result.Result[family.Family] {
		return h.createFamily(ctx, cmd)
	})
}

func (h CreateFamilyCommandHandler) createFamily(
	ctx context.Context,
	cmd CreateFamilyCommand) result.Result[family.Family] {
	if err := h.repository.Save(ctx, &cmd.Family); err != nil {
		return result.Fail[family.Family](err)
	}

	return result.Success(cmd.Family)
}
