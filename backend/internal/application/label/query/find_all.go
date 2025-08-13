package query

import (
	"context"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/internal/domain/label"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type FindAllQuery struct {
	Limit      int32
	Offset     int32
	SearchText string
}

func NewFindAllQuery(
	searchText string,
	limit int32,
	offset int32) FindAllQuery {
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
