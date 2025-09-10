package query

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/label"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/shared"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type FindAllQuery struct {
	Limit      int64
	Offset     int64
	SearchText string
}

func NewFindAllQuery(
	searchText string,
	limit, offset int64) FindAllQuery {
	return FindAllQuery{
		SearchText: searchText,
		Limit:      limit,
		Offset:     offset,
	}
}

type FindAllQueryHandler struct {
	labelRepository ports.LabelRepository
	authService     ports.AuthenticationService
}

func NewFindAllQueryHandler(
	labelRepository ports.LabelRepository,
	authService ports.AuthenticationService) *FindAllQueryHandler {
	return &FindAllQueryHandler{
		labelRepository: labelRepository,
		authService:     authService,
	}
}

func (h FindAllQueryHandler) Handle(
	ctx context.Context,
	query FindAllQuery) result.Result[shared.PaginatedResponse[label.Label]] {
	userId := h.authService.MustGetUserId(ctx)
	params := ports.NewLabelQueryParameters(query.SearchText, query.Limit, query.Offset)
	lbs, count, err := h.labelRepository.GetAll(ctx, userId, params)
	if err != nil {
		return result.Fail[shared.PaginatedResponse[label.Label]](err)
	}
	return result.Success(shared.NewPaginatedResponse(lbs, count))
}
