package authorization

import (
	"github.com/mistribe/subtracker/internal/domain/authorization"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
)

type permissionRequest struct {
	userId     types.UserID
	userRole   types.Role
	permission authorization.Permission
	userFamily *types.FamilyID
	error      error
}

func (r permissionRequest) For(entity ports.EntityWithOwnership) error {
	if r.error != nil {
		return r.error
	}
	owner := entity.Owner()

	if r.userRole == types.RoleAdmin {
		return nil
	}

	switch owner.Type() {
	case types.SystemOwnerType:
		return authorization.ErrUnauthorized
	case types.FamilyOwnerType:
		if r.userFamily == nil {
			return authorization.ErrUnauthorized
		}
		if owner.FamilyId() != *r.userFamily {
			return authorization.ErrUnauthorized
		}
		return nil
	case types.PersonalOwnerType:
		if owner.UserId() != r.userId {
			return authorization.ErrUnauthorized
		}
		return nil
	}

	return authorization.ErrUnauthorized
}
