package provider

import (
	"context"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/entity"
)

type Repository interface {
	entity.Repository[Provider]

	GetAll(ctx context.Context, parameters entity.QueryParameters) ([]Provider, int64, error)
	GetSystemProviders(ctx context.Context) ([]Provider, error)
	DeletePlan(ctx context.Context, planId uuid.UUID) (bool, error)
	DeletePrice(ctx context.Context, id uuid.UUID) (bool, error)
}
