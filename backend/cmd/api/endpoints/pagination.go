package endpoints

import (
	"github.com/mistribe/subtracker/internal/application/core"
	"github.com/mistribe/subtracker/pkg/slicesx"
)

// PaginatedResponseModel
// @name PaginateResponse
type PaginatedResponseModel[TValue any] struct {
	// Data contains the list of items for the current page
	Data []TValue `json:"data" binding:"required"`
	// Length represents the number of items in the current page
	Length int `json:"length" binding:"required"`
	// Total represents the total number of items available
	Total int64 `json:"total" binding:"required"`
}

func newPaginatedResponseModel[TValue any, TOut any](
	p core.PaginatedResponse[TValue],
	mapper func(TValue) TOut) PaginatedResponseModel[TOut] {
	return PaginatedResponseModel[TOut]{
		Data:   slicesx.Select(p.Data(), mapper),
		Length: p.Length(),
		Total:  p.Total(),
	}
}
