package auth

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	GetUserFamilies(ctx context.Context, userId string) ([]uuid.UUID, error)
}
