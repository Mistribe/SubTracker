package query

import (
	"context"

	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/domain/auth"
	"github.com/mistribe/subtracker/internal/domain/label"
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
	labelRepository label.Repository
	authService     auth.Service
}

func NewFindOneQueryHandler(
	labelRepository label.Repository,
	authService auth.Service) *FindOneQueryHandler {
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
