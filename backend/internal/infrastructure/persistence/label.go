package persistence

import (
	"context"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/application/core/option"
	"github.com/oleexo/subtracker/internal/domain/label"
)

type LabelRepository struct {
	labels map[uuid.UUID]label.Label
}

func NewLabelRepository() *LabelRepository {
	return &LabelRepository{
		labels: make(map[uuid.UUID]label.Label),
	}
}

func (r LabelRepository) Get(ctx context.Context, id uuid.UUID) (option.Option[label.Label], error) {
	if lbl, ok := r.labels[id]; ok {
		return option.Some(lbl), nil
	}
	return option.None[label.Label](), nil
}

func (r LabelRepository) GetAll(ctx context.Context, withDefault bool) ([]label.Label, error) {
	labels := make([]label.Label, 0, len(r.labels))
	for _, lbl := range r.labels {
		if !withDefault && lbl.IsDefault() {
			continue
		}
		labels = append(labels, lbl)
	}
	return labels, nil
}

func (r LabelRepository) GetDefaults(ctx context.Context) ([]label.Label, error) {
	labels := make([]label.Label, 0, len(r.labels))
	for _, lbl := range r.labels {
		if lbl.IsDefault() {
			labels = append(labels, lbl)
		}
	}
	return labels, nil
}

func (r LabelRepository) Save(ctx context.Context, lbl label.Label) error {
	r.labels[lbl.Id()] = lbl
	return nil
}

func (r LabelRepository) Delete(ctx context.Context, id uuid.UUID) error {
	delete(r.labels, id)
	return nil
}

func (r LabelRepository) Exists(ctx context.Context, ids ...uuid.UUID) (bool, error) {
	for _, id := range ids {
		if _, exists := r.labels[id]; !exists {
			return false, nil
		}
	}
	return true, nil
}
