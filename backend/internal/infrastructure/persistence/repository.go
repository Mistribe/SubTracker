package persistence

import (
	"context"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/entity"
	"github.com/oleexo/subtracker/pkg/langext/option"
)

type Repository[TEntity entity.Entity] struct {
	db *DatabaseContext
}

func (r Repository[TEntity]) GetById(ctx context.Context, entityId uuid.UUID) (option.Option[TEntity], error) {
	//TODO implement me
	panic("implement me")
}

func (r Repository[TEntity]) GetAll(ctx context.Context, f ...entity.QueryParameterFunc) ([]TEntity, error) {
	//TODO implement me
	panic("implement me")
}

func (r Repository[TEntity]) GetAllCount(ctx context.Context) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (r Repository[TEntity]) Save(ctx context.Context, entity TEntity) error {
	//TODO implement me
	panic("implement me")
}

func (r Repository[TEntity]) Delete(ctx context.Context, entityId uuid.UUID) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (r Repository[TEntity]) Exists(ctx context.Context, ids ...uuid.UUID) (bool, error) {
	//TODO implement me
	panic("implement me")
}
