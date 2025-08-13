package label

import (
	"context"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/entity"
)

type QueryParameters struct {
	entity.QueryParameters

	SearchText string
	UserId     string
}

func NewQueryParameters(
	userId string,
	searchText string,
	limit, offset int32) QueryParameters {
	return QueryParameters{
		QueryParameters: entity.QueryParameters{
			Limit:  limit,
			Offset: offset,
		},
		SearchText: searchText,
		UserId:     userId,
	}
}

type Repository interface {
	entity.Repository[Label]

	GetSystemLabels(ctx context.Context) ([]Label, error)
	GetAll(ctx context.Context, parameters QueryParameters) ([]Label, int64, error)
	GetByIdForUser(ctx context.Context, userId string, id uuid.UUID) (Label, error)
}
