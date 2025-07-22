package command

import (
	"context"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/family"
	"github.com/oleexo/subtracker/internal/domain/label"
	"github.com/oleexo/subtracker/internal/domain/subscription"
	"github.com/oleexo/subtracker/pkg/langext/option"
)

func createSubscription(
	ctx context.Context,
	labelRepository label.Repository,
	familyRepository family.Repository,
	newSubscription subscription.Subscription,
) (subscription.Subscription, error) {
	if err := ensureLabelsExists(ctx, labelRepository, newSubscription.Labels().Values()); err != nil {
		return newSubscription, err
	}
	if err := ensureFamilyMemberExists(ctx, familyRepository, newSubscription.FamilyId(),
		newSubscription.FamilyMembers().Values()); err != nil {
		return newSubscription, err
	}

	return newSubscription, nil
}

func ensureLabelsExists(
	ctx context.Context,
	labelRepository label.Repository,
	labels []uuid.UUID) error {
	if len(labels) == 0 {
		return nil
	}
	exists, err := labelRepository.Exists(ctx, labels...)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return label.ErrLabelNotFound
}

func ensureFamilyMemberExists(
	ctx context.Context,
	familyRepository family.Repository,
	familyId option.Option[uuid.UUID],
	members []uuid.UUID) error {
	return option.Match(familyId, func(familyId uuid.UUID) error {
		if len(members) == 0 {
			return nil
		}
		exists, err := familyRepository.MemberExists(ctx, familyId, members...)
		if err != nil {
			return err
		}
		if exists {
			return nil
		}
		return family.ErrFamilyMemberNotFound
	}, func() error {
		return nil
	})
}
