package user

import (
    "context"
    "github.com/oleexo/subtracker/internal/domain/family"
)

func EnsureOwnership(ctx context.Context,
    familyRepository family.Repository,
    owner Owner) error {
    userId, ok := FromContext(ctx)
    if !ok {
        return ErrNotAuthorized
    }
    switch owner.Type() {
    case SystemOwner:
        return ErrNotAuthorized
    case FamilyOwner:
        isMember, err := familyRepository.IsUserMemberOfFamily(ctx, owner.FamilyId(), userId)
        if err != nil {
            return err
        }
        if !isMember {
            return ErrNotAuthorized
        }
    case PersonalOwner:
        if owner.UserId() != userId {
            return ErrNotAuthorized
        }
    }
    return nil
}
