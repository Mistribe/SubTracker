package ports

import (
	"context"

	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/domain/provider"
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
	Repository[provider.Provider]

	GetByIdForUser(ctx context.Context, userId string, providerId uuid.UUID) (provider.Provider, error)
	GetAll(ctx context.Context, parameters ProviderQueryParameters) ([]provider.Provider, int64, error)
	GetAllForUser(ctx context.Context, userId string, parameters ProviderQueryParameters) ([]provider.Provider, int64,
		error)
	GetSystemProviders(ctx context.Context) ([]provider.Provider, int64, error)
	DeletePlan(ctx context.Context, planId uuid.UUID) (bool, error)
	DeletePrice(ctx context.Context, id uuid.UUID) (bool, error)
}
