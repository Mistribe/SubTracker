package label

import (
	"context"

	"github.com/oleexo/subtracker/internal/domain/entity"
	"github.com/oleexo/subtracker/internal/domain/user"
)

type QueryParameters struct {
	entity.QueryParameters

	Owners []user.OwnerType
}

func NewQueryParameters(limit, offset int, owners []user.OwnerType) QueryParameters {
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
