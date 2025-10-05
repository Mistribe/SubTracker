package ports

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/domain/types"
)

type FamilyRepository interface {
	Repository[types.FamilyID, family.Family]

	GetAccountFamily(ctx context.Context, userId types.UserID) (family.Family, error)
	MemberExists(ctx context.Context, familyId types.FamilyID, members ...types.FamilyMemberID) (bool, error)
	IsUserMemberOfFamily(ctx context.Context, familyId types.FamilyID, userId types.UserID) (bool, error)
}
