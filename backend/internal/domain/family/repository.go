package family

import (
	"context"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/pkg/langext/option"
)

type Repository interface {
	Get(ctx context.Context, id uuid.UUID) (option.Option[Member], error)
	GetAll(ctx context.Context) ([]Member, error)
	Save(ctx context.Context, member Member) error
	Delete(ctx context.Context, id uuid.UUID) error
	Exists(ctx context.Context, members ...uuid.UUID) (bool, error)
}
