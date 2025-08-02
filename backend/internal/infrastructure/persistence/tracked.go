package persistence

import (
	"context"

	"github.com/oleexo/subtracker/internal/infrastructure/persistence/sql"
	"github.com/oleexo/subtracker/pkg/slicesx"
)

func saveTrackedSlice[TEntity comparable](
	ctx context.Context,
	db *DatabaseContext,
	s *slicesx.Tracked[TEntity],
	addedFunc func(ctx context.Context, queries *sql.Queries, entities []TEntity) error,
	updatedFunc func(ctx context.Context, queries *sql.Queries, entity TEntity) error,
	removedFunc func(ctx context.Context, queries *sql.Queries, entity TEntity) error,
) error {
	if !s.HasChanges() {
		return nil
	}
	var addedEntities []TEntity
	for add := range s.Added() {
		addedEntities = append(addedEntities, add)
	}
	queries := db.GetQueries(ctx)
	if len(addedEntities) > 0 {
		if err := addedFunc(ctx, queries, addedEntities); err != nil {
			return err
		}
	}
	for remove := range s.Removed() {
		if err := removedFunc(ctx, queries, remove); err != nil {
			return err
		}
	}

	for update := range s.Updated() {
		if err := updatedFunc(ctx, queries, update); err != nil {
			return err
		}
	}

	return nil
}
