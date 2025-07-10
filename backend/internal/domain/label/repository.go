package label

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Get(ctx context.Context, id uuid.UUID) (*Label, error)
	GetAll(ctx context.Context) ([]Label, error)
	Save(ctx context.Context, label Label) error
	Exists(ctx context.Context, ids ...uuid.UUID) (bool, error)
}
