package query

import (
	"context"

	"github.com/oleexo/subtracker/internal/domain/family"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type FindUserFamilyQuery struct {
	UserId string
}

type FindUserFamilyQueryHandler struct {
	familyRepository family.Repository
}

func NewFindOneQueryHandler(familyRepository family.Repository) *FindUserFamilyQueryHandler {
	return &FindUserFamilyQueryHandler{familyRepository: familyRepository}
}

func (h FindUserFamilyQueryHandler) Handle(
	ctx context.Context,
	query FindUserFamilyQuery) result.Result[family.Family] {
	member, err := h.familyRepository.GetUserFamily(ctx, query.UserId)
	if err != nil {
		return result.Fail[family.Family](err)
	}

	if member == nil {
		return result.Fail[family.Family](family.ErrFamilyNotFound)
	}
	return result.Success(member)
}
