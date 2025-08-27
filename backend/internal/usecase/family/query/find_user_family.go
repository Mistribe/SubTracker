package query

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type FindUserFamilyQuery struct {
	UserId string
}

type FindUserFamilyQueryHandler struct {
	familyRepository ports.FamilyRepository
}

func NewFindOneQueryHandler(familyRepository ports.FamilyRepository) *FindUserFamilyQueryHandler {
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
