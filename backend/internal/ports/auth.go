package ports

import (
	"context"

	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/domain/auth"
)

type AuthService interface {
	MustGetUserId(ctx context.Context) string
	MustGetFamilies(ctx context.Context) []uuid.UUID
	IsInFamily(ctx context.Context, familyId uuid.UUID) bool
	IsOwner(ctx context.Context, owner auth.Owner) (bool, error)
}
