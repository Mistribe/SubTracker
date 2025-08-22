package family

import (
	"context"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/entity"
)

type Repository interface {
	entity.Repository[Family]

	GetUserFamily(ctx context.Context, userId string) (Family, error)
	// MemberExists checks if all provided member UUIDs exist within a specified family UUID context. Returns true if they exist.
	MemberExists(ctx context.Context, familyId uuid.UUID, members ...uuid.UUID) (bool, error)
	IsUserMemberOfFamily(ctx context.Context, familyId uuid.UUID, userId string) (bool, error)
}
