package command

import (
	"context"

	"github.com/oleexo/subtracker/internal/domain/family"
	"github.com/oleexo/subtracker/internal/domain/subscription"
	"github.com/oleexo/subtracker/pkg/langext/option"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type CreateSubscriptionCommand struct {
	Subscription subscription.Subscription
}

type CreateSubscriptionCommandHandler struct {
	subscriptionRepository subscription.Repository
	familyRepository       family.Repository
}

func NewCreateSubscriptionCommandHandler(
	subscriptionRepository subscription.Repository,
	familyRepository family.Repository) *CreateSubscriptionCommandHandler {
	return &CreateSubscriptionCommandHandler{
		subscriptionRepository: subscriptionRepository,
		familyRepository:       familyRepository,
	}
}

func (h CreateSubscriptionCommandHandler) Handle(
	ctx context.Context,
	cmd CreateSubscriptionCommand) result.Result[subscription.Subscription] {
	opt, err := h.subscriptionRepository.GetById(ctx, cmd.Subscription.Id())
	if err != nil {
		return result.Fail[subscription.Subscription](err)
	}

	return option.Match(opt, func(in subscription.Subscription) result.Result[subscription.Subscription] {
		if in.Equal(cmd.Subscription) {
			return result.Success(in)
		}
		return result.Fail[subscription.Subscription](subscription.ErrSubscriptionAlreadyExists)
	}, func() result.Result[subscription.Subscription] {
		return h.createSubscription(ctx, cmd)
	})

}

func (h CreateSubscriptionCommandHandler) createSubscription(
	ctx context.Context,
	cmd CreateSubscriptionCommand) result.Result[subscription.Subscription] {
	err := ensureRelatedEntityExists(ctx, h.familyRepository, cmd.Subscription)
	if err != nil {
		return result.Fail[subscription.Subscription](err)
	}

	if err = cmd.Subscription.GetValidationErrors(); err != nil {
		return result.Fail[subscription.Subscription](err)
	}

	err = h.subscriptionRepository.Save(ctx, cmd.Subscription)
	if err != nil {
		return result.Fail[subscription.Subscription](err)
	}

	return result.Success(cmd.Subscription)
}
