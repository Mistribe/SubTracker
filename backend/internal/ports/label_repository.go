package ports

import (
	"context"

	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/domain/label"
)

type LabelQueryParameters struct {
	QueryParameters

	SearchText string
}

func NewLabelQueryParameters(
	searchText string,
	limit, offset int64) LabelQueryParameters {
	return LabelQueryParameters{
		QueryParameters: QueryParameters{
			Limit:  limit,
			Offset: offset,
		},
		SearchText: searchText,
	}
}

type LabelRepository interface {
	Repository[label.Label]

	GetSystemLabels(ctx context.Context) ([]label.Label, error)
	GetAll(ctx context.Context, userId string, parameters LabelQueryParameters) ([]label.Label, int64, error)
	GetByIdForUser(ctx context.Context, userId string, id uuid.UUID) (label.Label, error)
}
