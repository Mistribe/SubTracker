package entity

import (
	"context"

	"github.com/google/uuid"
)

type QueryParameters struct {
	Limit  int
	Offset int
}

func NewQueryParameters(limit, offset int) QueryParameters {
	return QueryParameters{
		Limit:  limit,
		Offset: offset,
	}
}

type Repository[TEntity Entity] interface {
	GetById(ctx context.Context, entityId uuid.UUID) (TEntity, error)
	Save(ctx context.Context, entity TEntity) error
	Delete(ctx context.Context, entityId uuid.UUID) (bool, error)
	Exists(ctx context.Context, ids ...uuid.UUID) (bool, error)
}
