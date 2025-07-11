package persistence

import (
	"context"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/application/core/option"
	"github.com/oleexo/subtracker/internal/domain/label"
)

type LabelRepository struct {
}

func NewLabelRepository() *LabelRepository {
	return &LabelRepository{}
}

func (r LabelRepository) Get(ctx context.Context, id uuid.UUID) (option.Option[label.Label], error) {
	//TODO implement me
	panic("implement me")
}

func (r LabelRepository) GetAll(ctx context.Context) ([]label.Label, error) {
	//TODO implement me
	panic("implement me")
}

func (r LabelRepository) Save(ctx context.Context, label label.Label) error {
	//TODO implement me
	panic("implement me")
}

func (r LabelRepository) Delete(ctx context.Context, id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (r LabelRepository) Exists(ctx context.Context, ids ...uuid.UUID) (bool, error) {
	//TODO implement me
	panic("implement me")
}
