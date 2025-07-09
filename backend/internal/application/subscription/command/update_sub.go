package command

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/oleexo/subtracker/internal/application/core/option"
	"github.com/oleexo/subtracker/internal/application/core/result"
	"github.com/oleexo/subtracker/internal/domain/family"
	"github.com/oleexo/subtracker/internal/domain/label"
	"github.com/oleexo/subtracker/internal/domain/subscription"
)

type UpdateSubscriptionCommand struct {
	Id            uuid.UUID
	Name          string
	Labels        []uuid.UUID
	FamilyMembers []uuid.UUID
	Payer         option.Option[uuid.UUID]
	UpdatedAt     option.Option[time.Time]
}

type UpdateSubscriptionCommandHandler struct {
	subscriptionRepository subscription.Repository
	labelRepository        label.Repository
	familyRepository       family.Repository
}

func NewUpdateSubscriptionCommandHandler(subscriptionRepository subscription.Repository,
	labelRepository label.Repository,
	familyRepository family.Repository) *UpdateSubscriptionCommandHandler {
	return &UpdateSubscriptionCommandHandler{
		subscriptionRepository: subscriptionRepository,
		labelRepository:        labelRepository,
		familyRepository:       familyRepository,
	}
}

func (h UpdateSubscriptionCommandHandler) Handle(ctx context.Context, command UpdateSubscriptionCommand) result.Result[subscription.Subscription] {
	subOpt, err := h.subscriptionRepository.Get(ctx, command.Id)
	if err != nil {
		return result.Fail[subscription.Subscription](err)
	}
	return option.Match(subOpt, func(sub subscription.Subscription) result.Result[subscription.Subscription] {
		return h.updateSubscription(ctx, command, &sub)
	}, func() result.Result[subscription.Subscription] {
		return result.Fail[subscription.Subscription](subscription.ErrSubscriptionNotFound)
	})
}

func (h UpdateSubscriptionCommandHandler) updateSubscription(ctx context.Context,
	command UpdateSubscriptionCommand,
	sub *subscription.Subscription) result.Result[subscription.Subscription] {
	if err := h.ensureLabelsExists(ctx, command.Labels); err != nil {
		return result.Fail[subscription.Subscription](err)
	}

	if err := h.ensureFamilyMemberExists(ctx, command.FamilyMembers); err != nil {
		return result.Fail[subscription.Subscription](err)
	}

	sub.SetName(command.Name)
	sub.SetLabels(command.Labels)
	sub.SetFamilyMembers(command.FamilyMembers)
	sub.SetPayer(command.Payer)

	command.UpdatedAt.IfSome(func(updatedAt time.Time) {
		sub.SetUpdatedAt(updatedAt)
	})
	command.UpdatedAt.IfNone(func() {
		sub.SetUpdatedAt(time.Now())
	})

	if err := sub.Validate(); err != nil {
		return result.Fail[subscription.Subscription](err)
	}

	if err := h.subscriptionRepository.Save(ctx, sub); err != nil {
		return result.Fail[subscription.Subscription](err)
	}

	return result.Success(*sub)
}

func (h UpdateSubscriptionCommandHandler) ensureFamilyMemberExists(ctx context.Context, familyMembers []uuid.UUID) error {
	if len(familyMembers) == 0 {
		return nil
	}
	exists, err := h.familyRepository.Exists(ctx, familyMembers...)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return family.ErrFamilyMemberNotFound
}

func (h UpdateSubscriptionCommandHandler) ensureLabelsExists(ctx context.Context, labels []uuid.UUID) error {
	if len(labels) == 0 {
		return nil
	}
	exists, err := h.labelRepository.Exists(ctx, labels...)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return label.ErrLabelNotFound
}
