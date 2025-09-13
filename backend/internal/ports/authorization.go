package ports

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/auth"
	"github.com/mistribe/subtracker/internal/domain/user"
)

type Authorization interface {
	// EnsureWithinLimit validates if usage+delta is within the plan limit for the feature.
	// Returns ErrPlanLimitExceeded when exceeded.
	EnsureWithinLimit(ctx context.Context, feature user.Feature, delta int) error

	// Can checks if a user with the given userId has the specified permission within the provided scope.
	// Returns a boolean indicating authorization success and an error in case of failure.
	Can(ctx context.Context, permission user.Permission) PermissionRequest
}

type PermissionRequest interface {
	For(entity EntityWithOwnership) error
}

type EntityWithOwnership interface {
	Owner() auth.Owner
}
