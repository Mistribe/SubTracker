package query

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/label"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type FindOneQuery struct {
	LabelID types.LabelID
}

func NewFindOneQuery(labelID types.LabelID) FindOneQuery {
	return FindOneQuery{
		LabelID: labelID,
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
	connectedAccount := h.authService.MustGetConnectedAccount(ctx)
	lbl, err := h.labelRepository.GetByIdForUser(ctx, connectedAccount.UserID(), query.LabelID)
	if err != nil {
		return result.Fail[label.Label](err)
	}

	if lbl == nil {
		return result.Fail[label.Label](label.ErrLabelNotFound)

	}
	return result.Success(lbl)
}
