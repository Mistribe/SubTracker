package ports

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/label"
	"github.com/mistribe/subtracker/internal/domain/types"
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
	Repository[types.LabelID, label.Label]

	GetSystemLabels(ctx context.Context) ([]label.Label, error)
	GetAll(ctx context.Context, userId types.UserID, parameters LabelQueryParameters) ([]label.Label, int64, error)
	GetByIdForUser(ctx context.Context, userId types.UserID, id types.LabelID) (label.Label, error)
}
