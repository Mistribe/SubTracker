package command

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/subscription"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type PatchSubscriptionCommand struct {
	Subscription subscription.Subscription
}

type PatchSubscriptionCommandHandler struct {
	subscriptionRepository ports.SubscriptionRepository
	familyRepository       ports.FamilyRepository
}

func NewPatchSubscriptionCommandHandler(
	subscriptionRepository ports.SubscriptionRepository,
	familyRepository ports.FamilyRepository) *PatchSubscriptionCommandHandler {
	return &PatchSubscriptionCommandHandler{
		subscriptionRepository: subscriptionRepository,
		familyRepository:       familyRepository,
	}
}

func (h PatchSubscriptionCommandHandler) Handle(
	ctx context.Context,
	cmd PatchSubscriptionCommand) result.Result[subscription.Subscription] {
	sub, err := h.subscriptionRepository.GetById(ctx, cmd.Subscription.Id())
	if err != nil {
		return result.Fail[subscription.Subscription](err)
	}
	if sub == nil {
		return h.createSubscription(ctx, cmd)
	}
	return h.patchSubscription(ctx, cmd, sub)
}

func (h PatchSubscriptionCommandHandler) createSubscription(
	ctx context.Context,
	cmd PatchSubscriptionCommand) result.Result[subscription.Subscription] {
	err := ensureRelatedEntityExists(ctx, h.familyRepository, cmd.Subscription)
	if err != nil {
		return result.Fail[subscription.Subscription](err)
	}

	err = cmd.Subscription.GetValidationErrors()
	if err != nil {
		return result.Fail[subscription.Subscription](err)
	}

	err = h.subscriptionRepository.Save(ctx, cmd.Subscription)
	if err != nil {
		return result.Fail[subscription.Subscription](err)
	}

	return result.Success(cmd.Subscription)
}

func (h PatchSubscriptionCommandHandler) patchSubscription(
	ctx context.Context,
	cmd PatchSubscriptionCommand,
	sub subscription.Subscription) result.Result[subscription.Subscription] {
	if err := ensureRelatedEntityExists(ctx, h.familyRepository, cmd.Subscription); err != nil {
		return result.Fail[subscription.Subscription](err)
	}

	sub.SetFriendlyName(cmd.Subscription.FriendlyName())
	sub.SetFreeTrial(cmd.Subscription.FreeTrial())
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
