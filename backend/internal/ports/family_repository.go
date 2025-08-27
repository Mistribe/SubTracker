package ports

import (
	"context"

	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/domain/family"
)

type FamilyRepository interface {
	Repository[family.Family]

	GetUserFamily(ctx context.Context, userId string) (family.Family, error)
	// MemberExists checks if all provided member UUIDs exist within a specified family UUID context. Returns true if they exist.
	MemberExists(ctx context.Context, familyId uuid.UUID, members ...uuid.UUID) (bool, error)
	IsUserMemberOfFamily(ctx context.Context, familyId uuid.UUID, userId string) (bool, error)
}
