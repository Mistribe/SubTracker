package label

import (
	"context"

	"github.com/oleexo/subtracker/internal/domain/entity"
)

type QueryParameters struct {
	entity.QueryParameters

	WithDefaults bool
}

func NewQueryParameters(limit, offset int, withDefaults bool) QueryParameters {
	return QueryParameters{
		QueryParameters: entity.QueryParameters{
			Limit:  limit,
			Offset: offset,
		}, WithDefaults: withDefaults,
	}
}

type Repository interface {
	entity.Repository[Label]

	GetDefaults(ctx context.Context) ([]Label, error)
	GetAll(ctx context.Context, parameters QueryParameters) ([]Label, error)
	GetAllCount(ctx context.Context, parameters QueryParameters) (int64, error)
}
