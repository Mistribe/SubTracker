package query

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/domain/user"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/shared"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type FindUserFamilyQuery struct {
	UserId string
}

type FindUserFamilyQueryResponse struct {
	Family family.Family
	Limits shared.Limits
}

type FindUserFamilyQueryHandler struct {
	familyRepository ports.FamilyRepository
	authorization    ports.Authorization
}

func NewFindOneQueryHandler(
	familyRepository ports.FamilyRepository,
	authorization ports.Authorization) *FindUserFamilyQueryHandler {
	return &FindUserFamilyQueryHandler{
		familyRepository: familyRepository,
		authorization:    authorization,
	}
}

func (h FindUserFamilyQueryHandler) Handle(
	ctx context.Context,
	query FindUserFamilyQuery) result.Result[FindUserFamilyQueryResponse] {
	fam, err := h.familyRepository.GetUserFamily(ctx, query.UserId)
	if err != nil {
		return result.Fail[FindUserFamilyQueryResponse](err)
	}

	if fam == nil {
		return result.Fail[FindUserFamilyQueryResponse](family.ErrFamilyNotFound)
	}
	limits, err := h.authorization.GetCurrentLimits(ctx, user.CategoryFamily)
	if err != nil {
		return result.Fail[FindUserFamilyQueryResponse](err)
	}
	return result.Success(FindUserFamilyQueryResponse{
		Family: fam,
		Limits: limits,
	})
}
