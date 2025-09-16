package dto

import (
	"github.com/mistribe/subtracker/internal/shared"
	"github.com/mistribe/subtracker/pkg/x/herd"
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

func NewPaginatedResponseModel[TValue any, TOut any](
	p shared.PaginatedResponse[TValue],
	mapper func(TValue) TOut) PaginatedResponseModel[TOut] {
	return PaginatedResponseModel[TOut]{
		Data:   herd.Select(p.Data(), mapper),
		Length: p.Length(),
		Total:  p.Total(),
	}
}
