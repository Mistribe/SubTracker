package ports

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/auth"
	"github.com/mistribe/subtracker/internal/domain/user"
	"github.com/mistribe/subtracker/internal/shared"
)

type Authorization interface {
	// EnsureWithinLimit validates if usage+delta is within the plan limit for the feature.
	// Returns ErrPlanLimitExceeded when exceeded.
	EnsureWithinLimit(ctx context.Context, feature user.Feature, delta int) error

	// Can checks if a user with the given userId has the specified permission within the provided scope.
	// Returns a boolean indicating authorization success and an error in case of failure.
	Can(ctx context.Context, permission user.Permission) PermissionRequest

	// GetCurrentLimit retrieves the current usage and limit information for the specified feature in the context of the user.
	// Returns a LimitInfo object containing usage details and an error if fetching fails or the feature is unavailable.
	GetCurrentLimit(ctx context.Context, feature user.Feature) (shared.Limit, error)
	GetCurrentLimits(ctx context.Context, category user.Category) (shared.Limits, error)
}

// PermissionRequest defines an interface to handle permissions for entities that have ownership information.
type PermissionRequest interface {
	For(entity EntityWithOwnership) error
}

// EntityWithOwnership represents an entity that has an owner, typically used to determine ownership and access control.
// The Owner method retrieves the associated auth.Owner of the entity.
type EntityWithOwnership interface {
	Owner() auth.Owner
}

// LimitInfo represents an interface to handle limitations, including retrieval of limits and current usage information.
type LimitInfo interface {
	// Remaining retrieves the remaining limit, returning a pointer to an int64.
	Remaining() int64
	// Limit retrieves the maximum limit, returning a pointer to an int64.
	Limit() int64
	// HasLimit determines whether a limit is set, returning a boolean.
	HasLimit() bool
	// CurrentCount retrieves the current usage count, returning an int64.
	CurrentCount() int64
}
