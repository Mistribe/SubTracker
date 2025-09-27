package ports

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/provider"
	"github.com/mistribe/subtracker/internal/domain/types"
)

type ProviderQueryParameters struct {
	QueryParameters

	SearchText string
}

func NewProviderQueryParameters(
	searchText string,
	limit, offset int64) ProviderQueryParameters {
	return ProviderQueryParameters{
		QueryParameters: QueryParameters{
			Limit:  limit,
			Offset: offset,
		},
		SearchText: searchText,
	}
}

type ProviderRepository interface {
	Repository[types.ProviderID, provider.Provider]

	GetByIdForUser(ctx context.Context, userId types.UserID, providerId types.ProviderID) (provider.Provider, error)
	GetAll(ctx context.Context, parameters ProviderQueryParameters) ([]provider.Provider, int64, error)
	GetAllForUser(ctx context.Context, userId types.UserID, parameters ProviderQueryParameters) ([]provider.Provider,
		int64,
		error)
	GetSystemProviders(ctx context.Context) ([]provider.Provider, int64, error)
}
