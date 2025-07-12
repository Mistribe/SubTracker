package query

import (
	"context"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/label"
	"github.com/oleexo/subtracker/pkg/langext/option"
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
	lb, err := h.repository.Get(ctx, query.ID)
	if err != nil {
		return result.Fail[label.Label](err)
	}

	return option.Match(lb, func(lb label.Label) result.Result[label.Label] {
		return result.Success(lb)
	}, func() result.Result[label.Label] {
		return result.Fail[label.Label](label.ErrLabelNotFound)
	})
}
