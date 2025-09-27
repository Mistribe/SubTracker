package ports

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/entity"
)

type QueryParameters struct {
	Limit  int64
	Offset int64
}

func NewQueryParameters(limit, offset int64) QueryParameters {
	return QueryParameters{
		Limit:  limit,
		Offset: offset,
	}
}

type Repository[TKey comparable, TEntity entity.Entity[TKey]] interface {
	GetById(ctx context.Context, entityId TKey) (TEntity, error)
	Save(ctx context.Context, entities ...TEntity) error
	Delete(ctx context.Context, entityId TKey) (bool, error)
	Exists(ctx context.Context, ids ...TKey) (bool, error)
}
