package family

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/pkg/langext/option"
)

type Repository interface {

	// GetById retrieves a Family entity by its unique identifier. Returns an Option[Family] and an error, if any occurs.
	GetById(ctx context.Context, id uuid.UUID) (option.Option[Family], error)

	// GetOwn retrieves the family associated with the current user from the repository. It returns an Option[Family] or an error.
	GetOwn(ctx context.Context) (option.Option[Family], error)

	// GetAll retrieves all families from the repository. Returns a slice of Family and an error if retrieval fails.
	GetAll(ctx context.Context) ([]Family, error)

	// GetAllMembers retrieves all members associated with the specified family ID from the repository.
	GetAllMembers(ctx context.Context, id uuid.UUID) ([]Member, error)

	// Save persists a Family entity to the repository. Returns an error if the operation fails.
	Save(ctx context.Context, family *Family) error

	// SaveMember persists the given Member entity to the data store and returns an error if the operation fails.
	SaveMember(ctx context.Context, member *Member) error

	// Delete removes a family entity identified by the given UUID from the repository. Returns an error on failure.
	Delete(ctx context.Context, id uuid.UUID) error

	// DeleteMember removes a member with the given UUID from the family context. Returns an error if the deletion fails.
	DeleteMember(ctx context.Context, id uuid.UUID) error

	// MemberExists checks if all provided member UUIDs exist within a specified family UUID context. Returns true if they exist.
	MemberExists(ctx context.Context, familyId uuid.UUID, members ...uuid.UUID) (bool, error)

	// MarkAsUpdated updates the updatedAt timestamp for a specific family identified by familyId in the repository.
	MarkAsUpdated(ctx context.Context, familyId uuid.UUID, updatedAt time.Time) error
}
