package label

import (
	"context"

	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/internal/domain/entity"
)

type QueryParameters struct {
	entity.QueryParameters

	Owners []auth.OwnerType
}

func NewQueryParameters(limit, offset int, owners []auth.OwnerType) QueryParameters {
	if len(owners) == 0 {
		owners = append(owners, auth.SystemOwner, auth.PersonalOwner, auth.FamilyOwner)
	}
	return QueryParameters{
		QueryParameters: entity.QueryParameters{
			Limit:  limit,
			Offset: offset,
		},
		Owners: owners,
	}
}

type Repository interface {
	entity.Repository[Label]

	GetDefaults(ctx context.Context) ([]Label, error)
	GetAll(ctx context.Context, parameters QueryParameters) ([]Label, error)
	GetAllCount(ctx context.Context, parameters QueryParameters) (int64, error)
}
