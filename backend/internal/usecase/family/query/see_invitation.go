package query

import (
	"context"

	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type SeeInvitationQuery struct {
	FamilyId       uuid.UUID
	FamilyMemberId uuid.UUID
	InvitationCode string
}

type SeeInvitationQueryResponse struct {
	Family            family.Family
	InvitedInasmuchAs uuid.UUID
	UserId            string
}

type SeeInvitationQueryHandler struct {
	familyRepository ports.FamilyRepository
	authService      ports.AuthService
}

func NewSeeInvitationQueryHandler(familyRepository ports.FamilyRepository,
	authService ports.AuthService) *SeeInvitationQueryHandler {
	return &SeeInvitationQueryHandler{
		familyRepository: familyRepository,
		authService:      authService,
	}
}

func (h SeeInvitationQueryHandler) Handle(ctx context.Context,
	query SeeInvitationQuery) result.Result[SeeInvitationQueryResponse] {
	fam, err := h.familyRepository.GetById(ctx, query.FamilyId)
	if err != nil {
		return result.Fail[SeeInvitationQueryResponse](err)
	}

	if fam == nil {
		return result.Fail[SeeInvitationQueryResponse](family.ErrFamilyNotFound)
	}

	member := fam.GetMember(query.FamilyMemberId)
	if member == nil {
		return result.Fail[SeeInvitationQueryResponse](family.ErrFamilyMemberNotFound)
	}

	if member.InvitationCode() == nil {
		return result.Fail[SeeInvitationQueryResponse](family.ErrBadInvitationCode)
	}

	if *member.InvitationCode() != query.InvitationCode {
		return result.Fail[SeeInvitationQueryResponse](family.ErrBadInvitationCode)
	}

	userId := h.authService.MustGetUserId(ctx)

	return result.Success(SeeInvitationQueryResponse{
		Family:            fam,
		InvitedInasmuchAs: member.Id(),
		UserId:            userId,
	})

}
