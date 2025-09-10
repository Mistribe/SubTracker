package query

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/provider"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/shared"
	"github.com/mistribe/subtracker/pkg/langext/result"
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
	providerRepository ports.ProviderRepository
	authService        ports.AuthenticationService
}

func NewFindAllQueryHandler(
	providerRepository ports.ProviderRepository,
	authService ports.AuthenticationService) *FindAllQueryHandler {
	return &FindAllQueryHandler{
		providerRepository: providerRepository,
		authService:        authService,
	}
}

func (h FindAllQueryHandler) Handle(
	ctx context.Context,
	query FindAllQuery) result.Result[shared.PaginatedResponse[provider.Provider]] {
	userId := h.authService.MustGetUserId(ctx)
	providers, count, err := h.providerRepository.GetAllForUser(ctx, userId,
		ports.NewProviderQueryParameters(query.SearchText, query.Limit, query.Offset))
	if err != nil {
		return result.Fail[shared.PaginatedResponse[provider.Provider]](err)
	}

	return result.Success(shared.NewPaginatedResponse(providers, count))
}
