package query

import (
	"context"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/internal/domain/provider"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type FindAllQuery struct {
	Limit      int64
	Offset     int64
	SearchText string
}

func NewFindAllQuery(searchText string, limit, offset int64) FindAllQuery {
	return FindAllQuery{
		Limit:      limit,
		Offset:     offset,
		SearchText: searchText,
	}
}

type FindAllQueryHandler struct {
	providerRepository provider.Repository
	authService        auth.Service
}

func NewFindAllQueryHandler(
	providerRepository provider.Repository,
	authService auth.Service) *FindAllQueryHandler {
	return &FindAllQueryHandler{
		providerRepository: providerRepository,
		authService:        authService,
	}
}

func (h FindAllQueryHandler) Handle(
	ctx context.Context,
	query FindAllQuery) result.Result[core.PaginatedResponse[provider.Provider]] {
	userId := h.authService.MustGetUserId(ctx)
	providers, count, err := h.providerRepository.GetAllForUser(ctx, userId,
		provider.NewQueryParameters(query.SearchText, query.Limit, query.Offset))
	if err != nil {
		return result.Fail[core.PaginatedResponse[provider.Provider]](err)
	}

	return result.Success(core.NewPaginatedResponse(providers, count))
}
