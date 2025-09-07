package ports

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/auth"
)

type Feature string

const (
	FeatureActiveSubscriptions Feature = "active_subscriptions"
	FeatureCustomLabels        Feature = "custom_labels"
	FeatureCustomProviders     Feature = "custom_providers"
)

// Scope can be personal or family-owned. Reuse your existing Owner concept.
type LimitScope struct {
	Owner auth.Owner
}

// PlanLimiter checks plan-based constraints before commands mutate state.
type PlanLimiter interface {
	// EnsureWithinLimit validates if usage+delta is within plan limit for the feature.
	// Returns ErrPlanLimitExceeded when exceeded.
	EnsureWithinLimit(ctx context.Context, userId string, scope LimitScope, feature Feature, delta int) error
}

// Counting dependencies needed by PlanLimiter.
type SubscriptionCounter interface {
	CountActive(ctx context.Context, scope LimitScope) (int, error)
}

type LabelCounter interface {
	CountCustom(ctx context.Context, scope LimitScope) (int, error)
}

type ProviderCounter interface {
	CountCustom(ctx context.Context, scope LimitScope) (int, error)
}

// Optional: maps an HTTP or application error code to “plan exceeded” for clients.
type PlanLimitExceeded interface {
	error
	Feature() Feature
	Limit() int
	Usage() int
}
