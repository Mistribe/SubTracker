package ports

import (
	"context"

	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/domain/auth"
	"github.com/mistribe/subtracker/internal/domain/user"
)

type AuthenticationService interface {
	MustGetUserId(ctx context.Context) string
	MustGetFamilies(ctx context.Context) []uuid.UUID
	IsInFamily(ctx context.Context, familyId uuid.UUID) bool
	IsOwner(ctx context.Context, owner auth.Owner) (bool, error)
	MustGetUserRole(ctx context.Context) user.Role
}
