package user

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	GetUserFamilies(ctx context.Context, userId string) ([]uuid.UUID, error)
	GetUserProfile(ctx context.Context, userId string) (Profile, error)
	SaveProfile(ctx context.Context, profile Profile) error
}
