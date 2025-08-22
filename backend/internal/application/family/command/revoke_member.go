package command

import (
	"context"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/internal/domain/family"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type RevokeMemberCommand struct {
	FamilyId       uuid.UUID
	FamilyMemberId uuid.UUID
}

type RevokeMemberCommandHandler struct {
	familyRepository family.Repository
	authService      auth.Service
}

func NewRevokeMemberCommandHandler(
	familyRepository family.Repository,
	authService auth.Service) *RevokeMemberCommandHandler {
	return &RevokeMemberCommandHandler{familyRepository: familyRepository, authService: authService}
}

func (h RevokeMemberCommandHandler) Handle(ctx context.Context, cmd RevokeMemberCommand) result.Result[bool] {
	fam, err := h.familyRepository.GetById(ctx, cmd.FamilyId)
	if err != nil {
		return result.Fail[bool](err)
	}

	if fam == nil {
		return result.Fail[bool](family.ErrFamilyNotFound)
	}

	userId := h.authService.MustGetUserId(ctx)
	if fam.OwnerId() != userId {
		return result.Fail[bool](family.ErrOnlyOwnerCanEditFamily)
	}

	member := fam.GetMember(cmd.FamilyMemberId)
	if member == nil {
		return result.Fail[bool](family.ErrFamilyMemberNotFound)
	}

	if member.UserId() == nil {
		return result.Success(true)
	}

	member.SetUserId(nil)
	if err = fam.UpdateMember(member); err != nil {
		return result.Fail[bool](err)
	}

	if err = fam.GetValidationErrors(); err != nil {
		return result.Fail[bool](err)
	}

	if err = h.familyRepository.Save(ctx, fam); err != nil {
		return result.Fail[bool](err)
	}

	return result.Success(true)
}
