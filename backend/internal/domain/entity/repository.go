package entity

import (
	"context"

	"github.com/google/uuid"
)

type QueryParameters struct {
	Limit  int32
	Offset int32
}

func NewQueryParameters(limit, offset int32) QueryParameters {
	return QueryParameters{
		Limit:  limit,
		Offset: offset,
	}
}

type Repository[TEntity Entity] interface {
	GetById(ctx context.Context, entityId uuid.UUID) (TEntity, error)
	Save(ctx context.Context, entities ...TEntity) error
	Delete(ctx context.Context, entityId uuid.UUID) (bool, error)
	Exists(ctx context.Context, ids ...uuid.UUID) (bool, error)
}
