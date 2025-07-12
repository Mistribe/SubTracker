package query

import (
	"context"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/family"
	"github.com/oleexo/subtracker/pkg/langext/option"
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
	query FindOneQuery) result.Result[family.Member] {
	member, err := h.repository.Get(ctx, query.Id)
	if err != nil {
		return result.Fail[family.Member](err)
	}

	return option.Match[family.Member, result.Result[family.Member]](member,
		func(in family.Member) result.Result[family.Member] {
			return result.Success(in)
		},
		func() result.Result[family.Member] {
			return result.Fail[family.Member](family.ErrFamilyMemberNotFound)
		},
	)
}
