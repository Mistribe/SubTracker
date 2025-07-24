package endpoints

import (
	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/pkg/ext"
)

type paginatedResponseModel[TValue any] struct {
	Data   []TValue `json:"data" binding:"required"`
	Length int      `json:"length" binding:"required"`
	Total  int64    `json:"total" binding:"required"`
}

func newPaginatedResponseModel[TValue any, TOut any](
	p core.PaginatedResponse[TValue],
	mapper func(TValue) TOut) paginatedResponseModel[TOut] {
	return paginatedResponseModel[TOut]{
		Data:   ext.Map(p.Data(), mapper),
		Length: p.Length(),
		Total:  p.Total(),
	}
}
