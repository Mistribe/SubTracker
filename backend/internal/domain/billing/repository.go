package billing

import (
	"context"
)

type Account interface {
	PlanID() PlanID
}

// UsageRepository tracks usage counters for features in a given period.
// Implementations must ensure that Increment is atomic and safe under concurrency,
// typically via row-level locks within a transaction.
type UsageRepository interface {
	// Get returns the usage counter for the account/feature/period if it exists.
	Get(ctx context.Context, account Account, featureID FeatureID) (UsageCounter, bool, error)
}
