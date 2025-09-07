package ports

import (
	"context"

	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/domain/user"
)

type UserRepository interface {
	GetUserFamilies(ctx context.Context, userId string) ([]uuid.UUID, error)
	GetUser(ctx context.Context, userId string) (user.User, error)
	Save(ctx context.Context, profile user.User) error
	CreateDefault(ctx context.Context, userId string) (user.User, error)
}
