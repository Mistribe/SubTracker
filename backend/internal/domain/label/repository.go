package label

import (
	"context"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/application/core/option"
)

type Repository interface {
	Get(ctx context.Context, id uuid.UUID) (option.Option[Label], error)
	GetAll(ctx context.Context, withDefault bool) ([]Label, error)
	Save(ctx context.Context, label Label) error
	Delete(ctx context.Context, id uuid.UUID) error
	Exists(ctx context.Context, ids ...uuid.UUID) (bool, error)
}
