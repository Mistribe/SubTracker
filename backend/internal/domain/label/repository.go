package label

import (
	"context"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/internal/domain/entity"
)

type QueryParameters struct {
	entity.QueryParameters

	Owners   auth.OwnerTypes
	FamilyId *uuid.UUID
}

func NewQueryParameters(limit, offset int32,
	owners []auth.OwnerType,
	familyId *uuid.UUID) QueryParameters {
	if len(owners) == 0 {
		owners = append(owners, auth.SystemOwnerType, auth.PersonalOwnerType, auth.FamilyOwnerType)
	}
	return QueryParameters{
		QueryParameters: entity.QueryParameters{
			Limit:  limit,
			Offset: offset,
		},
		Owners:   owners,
		FamilyId: familyId,
	}
}

type Repository interface {
	entity.Repository[Label]

	GetSystemLabels(ctx context.Context) ([]Label, error)
	GetAll(ctx context.Context, parameters QueryParameters) ([]Label, int64, error)
}
