package command

import (
	"context"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/internal/domain/family"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type CreateFamilyCommand struct {
	Family family.Family
}

type CreateFamilyCommandHandler struct {
	familyRepository family.Repository
	authService      auth.Service
}

func NewCreateFamilyCommandHandler(familyRepository family.Repository,
	authService auth.Service) *CreateFamilyCommandHandler {
	return &CreateFamilyCommandHandler{
		familyRepository: familyRepository,
		authService:      authService,
	}
}

func (h CreateFamilyCommandHandler) Handle(ctx context.Context, cmd CreateFamilyCommand) result.Result[family.Family] {
	fam, err := h.familyRepository.GetById(ctx, cmd.Family.Id())
	if err != nil {
		return result.Fail[family.Family](err)
	}

	if fam != nil {
		if fam.Equal(cmd.Family) {
			return result.Success(cmd.Family)
		}
		return result.Fail[family.Family](family.ErrFamilyAlreadyExists)
	}
	return h.createFamily(ctx, cmd)
}

func (h CreateFamilyCommandHandler) createFamily(
	ctx context.Context,
	cmd CreateFamilyCommand) result.Result[family.Family] {

	userId := h.authService.MustGetUserId(ctx)
	memberId, err := uuid.NewV7()
	if err != nil {
		return result.Fail[family.Family](err)
	}

	creator := family.NewMember(
		memberId,
		cmd.Family.Id(),
		"You",
		false,
		cmd.Family.CreatedAt(),
		cmd.Family.UpdatedAt(),
	)
	creator.SetUserId(&userId)

	if err := cmd.Family.AddMember(creator); err != nil {
		return result.Fail[family.Family](err)
	}
	if err := cmd.Family.GetValidationErrors(); err != nil {
		return result.Fail[family.Family](err)
	}
	if err := h.familyRepository.Save(ctx, cmd.Family); err != nil {
		return result.Fail[family.Family](err)
	}

	return result.Success(cmd.Family)
}
