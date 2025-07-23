package label

import (
	"context"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/pkg/langext/option"
)

type Repository interface {
	Get(ctx context.Context, id uuid.UUID) (option.Option[Label], error)
	GetAll(ctx context.Context, size int, page int, withDefault bool) ([]Label, error)
	GetDefaults(ctx context.Context) ([]Label, error)
	Save(ctx context.Context, label *Label) error
	Delete(ctx context.Context, id uuid.UUID) error
	Exists(ctx context.Context, ids ...uuid.UUID) (bool, error)
	GetAllCount(ctx context.Context, withDefault bool) (int64, error)
}
