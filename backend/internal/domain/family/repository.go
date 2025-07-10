package family

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Exists(ctx context.Context, members ...uuid.UUID) (bool, error)
}
