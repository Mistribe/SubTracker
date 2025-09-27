package query

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type SeeInvitationQuery struct {
	FamilyId       types.FamilyID
	FamilyMemberId types.FamilyMemberID
	InvitationCode string
}

type SeeInvitationQueryResponse struct {
	Family            family.Family
	InvitedInasmuchAs types.FamilyMemberID
	UserId            types.UserID
}

type SeeInvitationQueryHandler struct {
	familyRepository ports.FamilyRepository
	authService      ports.Authentication
}

func NewSeeInvitationQueryHandler(
	familyRepository ports.FamilyRepository,
	authService ports.Authentication) *SeeInvitationQueryHandler {
	return &SeeInvitationQueryHandler{
		familyRepository: familyRepository,
		authService:      authService,
	}
}

func (h SeeInvitationQueryHandler) Handle(
	ctx context.Context,
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

	connectedAccount := h.authService.MustGetConnectedAccount(ctx)

	return result.Success(SeeInvitationQueryResponse{
		Family:            fam,
		InvitedInasmuchAs: member.Id(),
		UserId:            connectedAccount.UserID(),
	})

}
