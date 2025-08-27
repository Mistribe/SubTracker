package command

import (
	"context"

	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/domain/label"
	"github.com/mistribe/subtracker/internal/ports"
)

func ensureLabelExists(ctx context.Context,
	labelRepository ports.LabelRepository,
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
