package query

import (
	"context"

	"github.com/oleexo/subtracker/internal/domain/label"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type DefaultLabelQuery struct {
}

type DefaultLabelQueryHandler struct {
	repository label.Repository
}

func NewDefaultLabelQueryHandler(repository label.Repository) *DefaultLabelQueryHandler {
	return &DefaultLabelQueryHandler{repository: repository}
}

func (h DefaultLabelQueryHandler) Handle(ctx context.Context, query DefaultLabelQuery) result.Result[[]label.Label] {
	lbls, err := h.repository.GetSystemLabels(ctx)
	if err != nil {
		return result.Fail[[]label.Label](err)
	}
	if len(lbls) > 0 {
		return result.Success(lbls)
	}

	return result.Fail[[]label.Label](label.ErrMissingDefaultLabel)
}
