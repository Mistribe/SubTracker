package query

import (
	"context"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/family"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type FindOneQuery struct {
	Id uuid.UUID
}

type FindOneQueryHandler struct {
	repository family.Repository
}

func NewFindOneQueryHandler(repository family.Repository) *FindOneQueryHandler {
	return &FindOneQueryHandler{repository: repository}
}

func (h FindOneQueryHandler) Handle(
	ctx context.Context,
	query FindOneQuery) result.Result[family.Family] {
	member, err := h.repository.GetById(ctx, query.Id)
	if err != nil {
		return result.Fail[family.Family](err)
	}

	if member == nil {
		return result.Fail[family.Family](family.ErrFamilyNotFound)
	}
	return result.Success(member)
}
