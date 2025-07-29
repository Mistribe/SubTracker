package query

import (
	"context"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/label"
	"github.com/oleexo/subtracker/pkg/langext/result"
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

func (h FindOneQueryHandler) Handle(ctx context.Context, query FindOneQuery) result.Result[label.Label] {
	lbl, err := h.repository.GetById(ctx, query.ID)
	if err != nil {
		return result.Fail[label.Label](err)
	}

	if lbl == nil {
		return result.Fail[label.Label](label.ErrLabelNotFound)

	}
	return result.Success(lbl)
}
