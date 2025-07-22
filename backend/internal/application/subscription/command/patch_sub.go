package command

import (
	"context"

	"github.com/oleexo/subtracker/internal/domain/family"
	"github.com/oleexo/subtracker/internal/domain/label"
	"github.com/oleexo/subtracker/internal/domain/subscription"
	"github.com/oleexo/subtracker/pkg/langext/option"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type PatchSubscriptionCommand struct {
	Subscription subscription.Subscription
}

type PatchSubscriptionCommandHandler struct {
	subscriptionRepository subscription.Repository
	labelRepository        label.Repository
	familyRepository       family.Repository
}

func NewPatchSubscriptionCommandHandler(
	subscriptionRepository subscription.Repository,
	labelRepository label.Repository,
	familyRepository family.Repository) *PatchSubscriptionCommandHandler {
	return &PatchSubscriptionCommandHandler{
		subscriptionRepository: subscriptionRepository, labelRepository: labelRepository,
		familyRepository: familyRepository,
	}
}

func (h PatchSubscriptionCommandHandler) Handle(
	ctx context.Context,
	cmd PatchSubscriptionCommand) result.Result[subscription.Subscription] {
	if err := cmd.Subscription.Validate(); err != nil {
		return result.Fail[subscription.Subscription](err)
	}
	subOpt, err := h.subscriptionRepository.Get(ctx, cmd.Subscription.Id())
	if err != nil {
		return result.Fail[subscription.Subscription](err)
	}
	return option.Match(subOpt, func(fam subscription.Subscription) result.Result[subscription.Subscription] {
		return h.patchSubscription(ctx, cmd, fam)
	}, func() result.Result[subscription.Subscription] {
		return h.createSubscription(ctx, cmd)
	})
}

func (h PatchSubscriptionCommandHandler) createSubscription(
	ctx context.Context,
	cmd PatchSubscriptionCommand) result.Result[subscription.Subscription] {
	newSub, err := createSubscription(ctx, h.labelRepository, h.familyRepository, cmd.Subscription)
	if err != nil {
		return result.Fail[subscription.Subscription](err)
	}

	err = h.subscriptionRepository.Save(ctx, &newSub)
	if err != nil {
		return result.Fail[subscription.Subscription](err)
	}

	return result.Success(newSub)
}

func (h PatchSubscriptionCommandHandler) patchSubscription(
	ctx context.Context,
	cmd PatchSubscriptionCommand,
	sub subscription.Subscription) result.Result[subscription.Subscription] {
	if err := ensureLabelsExists(ctx, h.labelRepository, cmd.Subscription.Labels().Values()); err != nil {
		return result.Fail[subscription.Subscription](err)
	}

	if err := ensureFamilyMemberExists(ctx, h.familyRepository, cmd.Subscription.FamilyId(),
		cmd.Subscription.FamilyMembers().Values()); err != nil {
		return result.Fail[subscription.Subscription](err)
	}

	sub.SetName(cmd.Subscription.Name())
	sub.SetLabels(cmd.Subscription.Labels().Values())
	sub.SetFamilyMembers(cmd.Subscription.FamilyMembers().Values())
	sub.SetPayer(cmd.Subscription.Payer())
	sub.SetPayedByJointAccount(cmd.Subscription.PayedByJointAccount())
	sub.SetFamilyId(cmd.Subscription.FamilyId())
	sub.SetUpdatedAt(cmd.Subscription.UpdatedAt())

	for payment := range sub.Payments().It() {
		if !cmd.Subscription.Payments().Contains(payment) {
			sub.RemovePayment(payment.Id())
		}
	}

	for payment := range cmd.Subscription.Payments().It() {
		existingPayment, ok := sub.Payments().Get(payment)
		if !ok {
			sub.AddPayment(payment)
		} else {
			existingPayment.SetPrice(payment.Price())
			existingPayment.SetStartDate(payment.StartDate())
			existingPayment.SetEndDate(payment.EndDate())
			existingPayment.SetMonths(payment.Months())
			existingPayment.SetCurrency(payment.Currency())
			existingPayment.SetUpdatedAt(payment.UpdatedAt())
			sub.Payments().Update(existingPayment)
		}
	}

	if err := sub.Validate(); err != nil {
		return result.Fail[subscription.Subscription](err)
	}

	if err := h.subscriptionRepository.Save(ctx, &sub); err != nil {
		return result.Fail[subscription.Subscription](err)
	}

	return result.Success(sub)
}
