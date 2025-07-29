package command

import (
	"context"

	"github.com/oleexo/subtracker/internal/domain/family"
	"github.com/oleexo/subtracker/internal/domain/subscription"
	"github.com/oleexo/subtracker/pkg/langext/option"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type UpdateSubscriptionCommand struct {
	Subscription subscription.Subscription
}

type UpdateSubscriptionCommandHandler struct {
	subscriptionRepository subscription.Repository
	familyRepository       family.Repository
}

func NewUpdateSubscriptionCommandHandler(
	subscriptionRepository subscription.Repository,
	familyRepository family.Repository) *UpdateSubscriptionCommandHandler {
	return &UpdateSubscriptionCommandHandler{
		subscriptionRepository: subscriptionRepository,
		familyRepository:       familyRepository,
	}
}

func (h UpdateSubscriptionCommandHandler) Handle(
	ctx context.Context,
	cmd UpdateSubscriptionCommand) result.Result[subscription.Subscription] {
	subOpt, err := h.subscriptionRepository.GetById(ctx, cmd.Subscription.Id())
	if err != nil {
		return result.Fail[subscription.Subscription](err)
	}
	return option.Match(subOpt, func(sub subscription.Subscription) result.Result[subscription.Subscription] {
		return h.updateSubscription(ctx, cmd, sub)
	}, func() result.Result[subscription.Subscription] {
		return result.Fail[subscription.Subscription](subscription.ErrSubscriptionNotFound)
	})
}

func (h UpdateSubscriptionCommandHandler) updateSubscription(
	ctx context.Context,
	cmd UpdateSubscriptionCommand,
	sub subscription.Subscription) result.Result[subscription.Subscription] {
	if err := ensureRelatedEntityExists(ctx, h.familyRepository, cmd.Subscription); err != nil {
		return result.Fail[subscription.Subscription](err)
	}

	sub.SetFriendlyName(cmd.Subscription.FriendlyName())
	sub.SetFreeTrialDays(cmd.Subscription.FreeTrialDays())
	sub.SetOwner(cmd.Subscription.Owner())
	sub.SetPayer(cmd.Subscription.Payer())
	sub.SetServiceUsers(cmd.Subscription.ServiceUsers().Values())
	sub.SetStartDate(cmd.Subscription.StartDate())
	sub.SetEndDate(cmd.Subscription.EndDate())
	sub.SetRecurrency(cmd.Subscription.Recurrency())
	sub.SetCustomRecurrency(cmd.Subscription.CustomRecurrency())
	sub.SetUpdatedAt(cmd.Subscription.UpdatedAt())

	if err := sub.GetValidationErrors(); err != nil {
		return result.Fail[subscription.Subscription](err)
	}

	if err := h.subscriptionRepository.Save(ctx, sub); err != nil {
		return result.Fail[subscription.Subscription](err)
	}

	return result.Success(sub)
}
