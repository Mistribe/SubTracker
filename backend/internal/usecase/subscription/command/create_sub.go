package command

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/subscription"
	"github.com/mistribe/subtracker/internal/domain/user"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type CreateSubscriptionCommand struct {
	Subscription subscription.Subscription
}

type CreateSubscriptionCommandHandler struct {
	subscriptionRepository ports.SubscriptionRepository
	authorization          ports.Authorization
	familyRepository       ports.FamilyRepository
}

func NewCreateSubscriptionCommandHandler(
	subscriptionRepository ports.SubscriptionRepository,
	authorization ports.Authorization,
	familyRepository ports.FamilyRepository) *CreateSubscriptionCommandHandler {
	return &CreateSubscriptionCommandHandler{
		subscriptionRepository: subscriptionRepository,
		familyRepository:       familyRepository,
		authorization:          authorization,
	}
}

func (h CreateSubscriptionCommandHandler) Handle(
	ctx context.Context,
	cmd CreateSubscriptionCommand) result.Result[subscription.Subscription] {
	if err := h.authorization.Can(ctx, user.PermissionWrite).For(cmd.Subscription); err != nil {
		return result.Fail[subscription.Subscription](err)
	}
	sub, err := h.subscriptionRepository.GetById(ctx, cmd.Subscription.Id())
	if err != nil {
		return result.Fail[subscription.Subscription](err)
	}

	if sub != nil {
		if sub.Equal(cmd.Subscription) {
			return result.Success(sub)
		}
		return result.Fail[subscription.Subscription](subscription.ErrSubscriptionAlreadyExists)
	}

	err = h.authorization.EnsureWithinLimit(ctx, user.FeatureActiveSubscriptions, 1)
	if err != nil {
		return result.Fail[subscription.Subscription](err)
	}

	return h.createSubscription(ctx, cmd)
}

func (h CreateSubscriptionCommandHandler) createSubscription(
	ctx context.Context,
	cmd CreateSubscriptionCommand) result.Result[subscription.Subscription] {
	if err := ensureRelatedEntityExists(ctx, h.familyRepository, cmd.Subscription); err != nil {
		return result.Fail[subscription.Subscription](err)
	}

	if err := cmd.Subscription.GetValidationErrors(); err != nil {
		return result.Fail[subscription.Subscription](err)
	}

	if err := h.subscriptionRepository.Save(ctx, cmd.Subscription); err != nil {
		return result.Fail[subscription.Subscription](err)
	}

	return result.Success(cmd.Subscription)
}
