package endpoints

import (
	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/pkg/ext"
)

type paginatedResponseModel[TValue any] struct {
	Data   []TValue `json:"data"`
	Length int      `json:"length"`
	Total  int64    `json:"total"`
}

func newPaginatedResponseModel[TValue any, TOut any](p core.PaginatedResponse[TValue],
	mapper func(TValue) TOut) paginatedResponseModel[TOut] {
	return paginatedResponseModel[TOut]{
		Data:   ext.Map(p.Data(), mapper),
		Length: p.Length(),
		Total:  p.Total(),
	}
}
