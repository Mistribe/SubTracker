package ports

import (
	"context"

	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/domain/user"
)

type UserRepository interface {
	GetUserFamilies(ctx context.Context, userId string) ([]uuid.UUID, error)
	GetUserProfile(ctx context.Context, userId string) (user.Profile, error)
	SaveProfile(ctx context.Context, profile user.Profile) error
}
