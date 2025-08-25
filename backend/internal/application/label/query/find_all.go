package query

import (
	"context"

	"github.com/mistribe/subtracker/internal/application/core"
	"github.com/mistribe/subtracker/internal/domain/auth"
	"github.com/mistribe/subtracker/internal/domain/label"
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
	labelRepository label.Repository
	authService     auth.Service
}

func NewFindAllQueryHandler(
	labelRepository label.Repository,
	authService auth.Service) *FindAllQueryHandler {
	return &FindAllQueryHandler{
		labelRepository: labelRepository,
		authService:     authService,
	}
}

func (h FindAllQueryHandler) Handle(
	ctx context.Context,
	query FindAllQuery) result.Result[core.PaginatedResponse[label.Label]] {
	userId := h.authService.MustGetUserId(ctx)
	params := label.NewQueryParameters(query.SearchText, query.Limit, query.Offset)
	lbs, count, err := h.labelRepository.GetAll(ctx, userId, params)
	if err != nil {
		return result.Fail[core.PaginatedResponse[label.Label]](err)
	}
	return result.Success(core.NewPaginatedResponse(lbs, count))
}
