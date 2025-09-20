package query

import (
	"context"

	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/domain/label"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/langext/result"
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
	labelRepository ports.LabelRepository
	authService     ports.Authentication
}

func NewFindOneQueryHandler(
	labelRepository ports.LabelRepository,
	authService ports.Authentication) *FindOneQueryHandler {
	return &FindOneQueryHandler{
		labelRepository: labelRepository,
		authService:     authService,
	}
}

func (h FindOneQueryHandler) Handle(ctx context.Context, query FindOneQuery) result.Result[label.Label] {
	userId := h.authService.MustGetUserId(ctx)
	lbl, err := h.labelRepository.GetByIdForUser(ctx, userId, query.ID)
	if err != nil {
		return result.Fail[label.Label](err)
	}

	if lbl == nil {
		return result.Fail[label.Label](label.ErrLabelNotFound)

	}
	return result.Success(lbl)
}
