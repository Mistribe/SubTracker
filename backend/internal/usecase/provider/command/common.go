package command

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/label"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
)

func ensureLabelExists(ctx context.Context,
	labelRepository ports.LabelRepository,
	labels []types.LabelID) error {
	labelExists, err := labelRepository.Exists(ctx, labels...)
	if err != nil {
		return err
	}
	if !labelExists {
		return label.ErrLabelNotFound
	}
	return nil
}
