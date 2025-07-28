package family

import (
	"context"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/entity"
)

type Repository interface {
	entity.Repository[Family]

	GetAll(ctx context.Context, parameters entity.QueryParameters) ([]Family, error)
	GetAllCount(ctx context.Context) (int64, error)

	// GetAllMembers retrieves all members associated with the specified family ID from the repository.
	GetAllMembers(ctx context.Context, id uuid.UUID) ([]Member, error)

	// MemberExists checks if all provided member UUIDs exist within a specified family UUID context. Returns true if they exist.
	MemberExists(ctx context.Context, familyId uuid.UUID, members ...uuid.UUID) (bool, error)
}
