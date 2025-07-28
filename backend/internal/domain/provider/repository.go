package provider

import (
	"context"

	"github.com/oleexo/subtracker/internal/domain/entity"
)

type Repository interface {
	entity.Repository[Provider]

	GetAll(ctx context.Context, parameters entity.QueryParameters) ([]Provider, error)
	GetAllCount(ctx context.Context) (int64, error)
}
