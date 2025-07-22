package persistence

import (
	"context"

	"github.com/oleexo/subtracker/pkg/slicesx"
)

func saveTrackedSlice[T comparable](
	ctx context.Context,
	s *slicesx.Tracked[T],
	createFunc func(context.Context, T) error,
	updateFunc func(context.Context, T) error,
	deleteFunc func(context.Context, T) error,
) error {
	for add := range s.Added() {
		if err := createFunc(ctx, add); err != nil {
			return err
		}
	}
	for remove := range s.Removed() {
		if err := deleteFunc(ctx, remove); err != nil {
			return err
		}
	}

	for update := range s.Updated() {
		if err := updateFunc(ctx, update); err != nil {
			return err
		}
	}

	return nil
}
