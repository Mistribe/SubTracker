package command

import (
	"context"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/application/core/result"
	"github.com/oleexo/subtracker/internal/domain/family"
	"github.com/oleexo/subtracker/internal/domain/label"
	"github.com/oleexo/subtracker/internal/domain/subscription"
)

type CreateSubscriptionCommand struct {
	Subscription subscription.Subscription
}

type CreateSubscriptionCommandHandler struct {
	subscriptionRepository subscription.Repository
	labelRepository        label.Repository
	familyRepository       family.Repository
}

func NewCreateSubscriptionCommandHandler(subscriptionRepository subscription.Repository,
	labelRepository label.Repository,
	familyRepository family.Repository) *CreateSubscriptionCommandHandler {
	return &CreateSubscriptionCommandHandler{
		subscriptionRepository: subscriptionRepository,
		labelRepository:        labelRepository,
		familyRepository:       familyRepository,
	}
}

func (h CreateSubscriptionCommandHandler) Handle(ctx context.Context,
	command CreateSubscriptionCommand) result.Result[subscription.Subscription] {
	if err := h.assertNotExists(ctx, command.Subscription.Id()); err != nil {
		return result.Fail[subscription.Subscription](err)
	}
	if err := h.ensureLabelsExists(ctx, command.Subscription.Labels()); err != nil {
		return result.Fail[subscription.Subscription](err)
	}
	if err := h.ensureFamilyMemberExists(ctx, command.Subscription.FamilyMembers()); err != nil {
		return result.Fail[subscription.Subscription](err)
	}

	if err := h.subscriptionRepository.Save(ctx, &command.Subscription); err != nil {
		return result.Fail[subscription.Subscription](err)
	}

	return result.Success(command.Subscription)
}

func (h CreateSubscriptionCommandHandler) ensureFamilyMemberExists(ctx context.Context,
	familyMembers []uuid.UUID) error {
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

func (h CreateSubscriptionCommandHandler) ensureLabelsExists(ctx context.Context, labels []uuid.UUID) error {
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

func (h CreateSubscriptionCommandHandler) assertNotExists(ctx context.Context, id uuid.UUID) error {
	sub, err := h.subscriptionRepository.Get(ctx, id)
	if err != nil {
		return err
	}
	if sub != nil {
		return subscription.ErrSubscriptionAlreadyExists
	}

	return nil
}
