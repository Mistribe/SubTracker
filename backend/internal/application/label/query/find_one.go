package query

import (
	"context"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/application/core/option"
	"github.com/oleexo/subtracker/internal/application/core/result"
	"github.com/oleexo/subtracker/internal/domain/label"
)

type FindOneQuery struct {
	ID uuid.UUID
}

func NewFindOneQuery(id uuid.UUID) FindOneQuery {
	return FindOneQuery{
		ID: id,
	}
}

type FindOneQueryHandler struct {
	repository label.Repository
}

func NewFindOneQueryHandler(repository label.Repository) *FindOneQueryHandler {
	return &FindOneQueryHandler{repository: repository}
}

func (h FindOneQueryHandler) Handle(ctx context.Context, query FindOneQuery) result.Result[option.Option[label.Label]] {
	lb, err := h.repository.Get(ctx, query.ID)
	if err != nil {
		return result.Fail[option.Option[label.Label]](err)
	}

	return result.Success(lb)
}
