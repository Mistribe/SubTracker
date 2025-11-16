package export

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
)

// LabelResolver resolves label IDs to their names
type LabelResolver interface {
	// ResolveLabelNames takes a slice of label IDs and returns their names
	ResolveLabelNames(ctx context.Context, labelIds []types.LabelID) ([]string, error)
}

type labelResolver struct {
	labelRepo ports.LabelRepository
}

// NewLabelResolver creates a new label resolver
func NewLabelResolver(labelRepo ports.LabelRepository) LabelResolver {
	return &labelResolver{
		labelRepo: labelRepo,
	}
}

func (r *labelResolver) ResolveLabelNames(ctx context.Context, labelIds []types.LabelID) ([]string, error) {
	if len(labelIds) == 0 {
		return []string{}, nil
	}

	names := make([]string, 0, len(labelIds))

	for _, labelId := range labelIds {
		label, err := r.labelRepo.GetById(ctx, labelId)
		if err != nil {
			// Skip labels that cannot be found (they may have been deleted)
			continue
		}
		names = append(names, label.Name())
	}

	return names, nil
}
