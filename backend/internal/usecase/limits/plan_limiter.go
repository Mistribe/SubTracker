package limits

import (
	"context"
	"fmt"

	"github.com/mistribe/subtracker/internal/ports"
)

type planLimiter struct {
	userRepository ports.UserRepository
}

func NewPlanLimiter(userRepository ports.UserRepository) ports.PlanLimiter {
	return &planLimiter{
		userRepository: userRepository,
	}
}

type planExceeded struct {
	feature ports.Feature
	limit   int
	usage   int
}

func (e planExceeded) Error() string {
	return fmt.Sprintf("plan limit exceeded for %s: usage=%d limit=%d", e.feature, e.usage, e.limit)
}
func (e planExceeded) Feature() ports.Feature { return e.feature }
func (e planExceeded) Limit() int             { return e.limit }
func (e planExceeded) Usage() int             { return e.usage }

func (l *planLimiter) EnsureWithinLimit(ctx context.Context, userId string, scope ports.LimitScope,
	feature ports.Feature, delta int) error {
	// 1) Load user and plan
	u := l.userService.Get(ctx, userId) // or however you fetch user profile in your application layer
	plan := u.Plan()
	limits := plan.Limits()

	// 2) Resolve the numeric limit for the feature
	var planLimit *int
	switch feature {
	case ports.FeatureActiveSubscriptions:
		planLimit = limits.MaxActiveSubscriptions()
	case ports.FeatureCustomLabels:
		planLimit = limits.MaxCustomLabels()
	case ports.FeatureCustomProviders:
		planLimit = limits.MaxCustomProviders()
	default:
		return nil
	}

	// Unlimited
	if planLimit == nil {
		return nil
	}

	// 3) Count usage for scope
	var usage int
	var err error
	switch feature {
	case ports.FeatureActiveSubscriptions:
		usage, err = l.subscriptionCount.CountActive(ctx, scope)
	case ports.FeatureCustomLabels:
		usage, err = l.labelCount.CountCustom(ctx, scope)
	case ports.FeatureCustomProviders:
		usage, err = l.providerCount.CountCustom(ctx, scope)
	}
	if err != nil {
		return err
	}

	// 4) Compare usage+delta
	if usage+delta > *planLimit {
		return planExceeded{
			feature: feature,
			limit:   *planLimit,
			usage:   usage,
		}
	}

	return nil
}
