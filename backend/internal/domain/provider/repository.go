package provider

import (
	"context"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/entity"
)

type QueryParameters struct {
	entity.QueryParameters

	SearchText string
}

func NewQueryParameters(
	searchText string,
	limit, offset int64) QueryParameters {
	return QueryParameters{
		QueryParameters: entity.QueryParameters{
			Limit:  limit,
			Offset: offset,
		},
		SearchText: searchText,
	}
}

type Repository interface {
	entity.Repository[Provider]

	GetByIdForUser(ctx context.Context, userId string, providerId uuid.UUID) (Provider, error)
	GetAll(ctx context.Context, parameters QueryParameters) ([]Provider, int64, error)
	GetAllForUser(ctx context.Context, userId string, parameters QueryParameters) ([]Provider, int64, error)
	GetSystemProviders(ctx context.Context) ([]Provider, int64, error)
	DeletePlan(ctx context.Context, planId uuid.UUID) (bool, error)
	DeletePrice(ctx context.Context, id uuid.UUID) (bool, error)
}
