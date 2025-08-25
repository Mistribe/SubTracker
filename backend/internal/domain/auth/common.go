package auth

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/family"
)

func EnsureOwnership(ctx context.Context,
	familyRepository family.Repository,
	owner Owner) error {
	userId, ok := GetUserIdFromContext(ctx)
	if !ok {
		return ErrNotAuthorized
	}
	switch owner.Type() {
	case SystemOwnerType:
		return ErrNotAuthorized
	case FamilyOwnerType:
		isMember, err := familyRepository.IsUserMemberOfFamily(ctx, owner.FamilyId(), userId)
		if err != nil {
			return err
		}
		if !isMember {
			return ErrNotAuthorized
		}
	case PersonalOwnerType:
		if owner.UserId() != userId {
			return ErrNotAuthorized
		}
	}
	return nil
}
