package dto

import (
	"github.com/mistribe/subtracker/internal/shared"
)

// Limit represents a structure containing information about remaining quota, total limit, and whether a limit is set.
type Limit struct {
	Category string
	Feature  string

	// Remaining indicates the remaining quota available for use.
	Remaining int64

	// Limit represents a structure containing quota limits, remaining values, and whether a limit is applied.
	Limit int64

	// HasLimit indicates whether a limit is imposed.
	HasLimit bool
}

type Limits []Limit

func NewLimits(limits []shared.Limit) Limits {
	results := make([]Limit, len(limits))
	for i, limit := range limits {
		results[i] = Limit{
			Category:  limit.Category(),
			Feature:   limit.Feature(),
			Remaining: limit.Remaining(),
			Limit:     limit.Limit(),
			HasLimit:  limit.HasLimit(),
		}
	}
	return results
}
