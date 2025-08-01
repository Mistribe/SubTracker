package command

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/internal/domain/family"
	"github.com/oleexo/subtracker/pkg/ext"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type CreateFamilyCommand struct {
	FamilyId    *uuid.UUID
	Name        string
	CreatorName string
	CreatedAt   *time.Time
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
	if cmd.FamilyId != nil {
		fam, err := h.familyRepository.GetById(ctx, *cmd.FamilyId)
		if err != nil {
			return result.Fail[family.Family](err)
		}

		if fam != nil {
			return result.Fail[family.Family](family.ErrFamilyAlreadyExists)
		}
	} else {
		newId, err := uuid.NewV7()
		if err != nil {
			return result.Fail[family.Family](err)
		}
		cmd.FamilyId = &newId
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

	createdAt := ext.ValueOrDefault(cmd.CreatedAt, time.Now())

	fam := family.NewFamily(
		*cmd.FamilyId,
		userId,
		cmd.Name,
		[]family.Member{},
		createdAt,
		createdAt,
	)

	creator := family.NewMember(
		memberId,
		*cmd.FamilyId,
		cmd.CreatorName,
		family.OwnerMemberType,
		createdAt,
		createdAt,
	)
	creator.SetUserId(&userId)

	if err := fam.AddMember(creator); err != nil {
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
