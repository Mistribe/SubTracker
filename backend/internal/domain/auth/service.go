package auth

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	MustGetUserId(ctx context.Context) string
	MustGetFamilies(ctx context.Context) []uuid.UUID
	IsInFamily(ctx context.Context, familyId uuid.UUID) bool
	IsOwner(ctx context.Context, owner Owner) (bool, error)
}
