package shared

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
)

type OwnerFactory interface {
	Resolve(ctx context.Context, t types.OwnerType) (types.Owner, error)
}

type ownerFactory struct {
	auth             ports.Authentication
	familyRepository ports.FamilyRepository
}

func NewOwnerFactory(
	auth ports.Authentication,
	familyRepository ports.FamilyRepository,
) OwnerFactory {
	return &ownerFactory{
		auth:             auth,
		familyRepository: familyRepository,
	}
}

func (f *ownerFactory) Resolve(ctx context.Context, t types.OwnerType) (types.Owner, error) {
	acc := f.auth.MustGetConnectedAccount(ctx)
	userID := acc.UserID()

	switch t {
	case types.PersonalOwnerType:
		return types.NewOwner(types.PersonalOwnerType, nil, &userID), nil

	case types.FamilyOwnerType:
		fam, err := f.familyRepository.GetAccountFamily(ctx, userID)
		if err != nil {
			return nil, err
		}
		if fam == nil {
			return nil, family.ErrFamilyNotFound
		}
		familyID := fam.Id()
		return types.NewOwner(types.FamilyOwnerType, &familyID, nil), nil

	case types.SystemOwnerType:
		return types.SystemOwner, nil

	default:
		return nil, types.ErrUnknownOwnerType
	}
}
