package command

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/langext/result"
	"github.com/mistribe/subtracker/pkg/x"
)

type CreateFamilyCommand struct {
	FamilyId    *uuid.UUID
	Name        string
	CreatorName string
	CreatedAt   *time.Time
}

type CreateFamilyCommandHandler struct {
	familyRepository ports.FamilyRepository
	authService      ports.Authentication
}

func NewCreateFamilyCommandHandler(
	familyRepository ports.FamilyRepository,
	authService ports.Authentication) *CreateFamilyCommandHandler {
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

	createdAt := x.ValueOrDefault(cmd.CreatedAt, time.Now())

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
		nil,
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
