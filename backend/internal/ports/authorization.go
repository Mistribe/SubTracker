package ports

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/user"
)

type Authorization interface {
	// EnsureWithinLimit validates if usage+delta is within the plan limit for the feature.
	// Returns ErrPlanLimitExceeded when exceeded.
	EnsureWithinLimit(ctx context.Context, userId string, feature user.Feature, delta int) error
}
