package command

import (
	"context"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/label"
)

func ensureLabelExists(ctx context.Context,
	labelRepository label.Repository,
	labels []uuid.UUID) error {
	labelExists, err := labelRepository.Exists(ctx, labels...)
	if err != nil {
		return err
	}
	if !labelExists {
		return label.ErrLabelNotFound
	}
	return nil
}
