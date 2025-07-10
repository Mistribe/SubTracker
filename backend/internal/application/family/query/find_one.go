package query

import (
	"context"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/application/core/option"
	"github.com/oleexo/subtracker/internal/application/core/result"
	"github.com/oleexo/subtracker/internal/domain/family"
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
	query FindOneQuery) result.Result[option.Option[family.Member]] {
	member, err := h.repository.Get(ctx, query.Id)
	if err != nil {
		return result.Fail[option.Option[family.Member]](err)
	}

	return result.Success(member)
}
