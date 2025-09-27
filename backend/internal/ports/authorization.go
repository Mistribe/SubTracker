package ports

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/authorization"
	"github.com/mistribe/subtracker/internal/domain/types"
)

type Authorization interface {
	// Can checks if a user with the given userId has the specified permission within the provided scope.
	// Returns a boolean indicating authorization success and an error in case of failure.
	Can(ctx context.Context, permission authorization.Permission) PermissionRequest
}

// PermissionRequest defines an interface to handle permissions for entities that have ownership information.
type PermissionRequest interface {
	For(entity EntityWithOwnership) error
}

// EntityWithOwnership represents an entity that has an owner, typically used to determine ownership and access control.
// The Owner method retrieves the associated auth.Owner of the entity.
type EntityWithOwnership interface {
	Owner() types.Owner
}
