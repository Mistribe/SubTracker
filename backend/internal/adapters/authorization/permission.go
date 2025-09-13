package authorization

import (
	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/domain/auth"
	"github.com/mistribe/subtracker/internal/domain/user"
	"github.com/mistribe/subtracker/internal/ports"
)

type permissionRequest struct {
	userId     string
	userRole   user.Role
	permission user.Permission
	userFamily *uuid.UUID
	error      error
}

func (r permissionRequest) For(entity ports.EntityWithOwnership) error {
	if r.error != nil {
		return r.error
	}
	owner := entity.Owner()

	if r.userRole == user.RoleAdmin {
		return nil
	}

	switch owner.Type() {
	case auth.SystemOwnerType:
		return user.ErrUnauthorized
	case auth.FamilyOwnerType:
		if r.userFamily == nil {
			return user.ErrUnauthorized
		}
		if owner.FamilyId() != *r.userFamily {
			return user.ErrUnauthorized
		}
		return nil
	case auth.PersonalOwnerType:
		if owner.UserId() != r.userId {
			return user.ErrUnauthorized
		}
	}

	return user.ErrUnauthorized
}
