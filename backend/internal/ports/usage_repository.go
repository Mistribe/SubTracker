package ports

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/billing"
	"github.com/mistribe/subtracker/internal/domain/types"
)

// UsageRepository tracks usage counters for features in a given period.
// Implementations must ensure that Increment is atomic and safe under concurrency,
// typically via row-level locks within a transaction.
type UsageRepository interface {
	// Get returns the usage counter for the account/feature/period if it exists.
	Get(ctx context.Context, userID types.UserID, feature billing.Feature) (billing.UsageCounter, bool, error)
	GetAll(ctx context.Context, userID types.UserID) ([]billing.UsageCounter, error)
}
