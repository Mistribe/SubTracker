package billing

import "errors"

var (
	// ErrNotFound is a generic not-found error for repository lookups.
	ErrNotFound = errors.New("not found")

	ErrPlanNotFound = errors.New("plan not found")

	// ErrFeatureNotFound signals that a feature key could not be resolved.
	ErrFeatureNotFound = errors.New("feature not found")

	// ErrFeatureDisabled is returned when a boolean feature is not enabled.
	ErrFeatureDisabled = errors.New("feature disabled for account")

	// ErrInvalidFeatureType is returned when an operation targets the wrong feature type.
	ErrInvalidFeatureType = errors.New("invalid feature type")

	// ErrQuotaExceeded is returned when a quota operation would exceed the allowed limit.
	ErrQuotaExceeded           = errors.New("quota exceeded")
	ErrCannotGetQuotaOnFeature = errors.New("cannot get quota on feature")
)
