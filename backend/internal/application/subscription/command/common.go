package command

import (
	"context"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/family"
	"github.com/oleexo/subtracker/internal/domain/subscription"
	"github.com/oleexo/subtracker/internal/domain/user"
)

func ensureRelatedEntityExists(
	ctx context.Context,
	familyRepository family.Repository,
	newSubscription subscription.Subscription,
) error {
	if newSubscription.Owner().Type() == user.FamilyOwner {
		familyId := newSubscription.Owner().FamilyId()

		var members []uuid.UUID
		if newSubscription.ServiceUsers() != nil && newSubscription.ServiceUsers().Len() > 0 {
			members = newSubscription.ServiceUsers().Values()
		}
		if newSubscription.Payer() != nil && newSubscription.Payer().Type() == subscription.FamilyMemberPayer {
			members = append(members, newSubscription.Payer().MemberId())
		}
		if err := ensureFamilyMemberExists(ctx, familyRepository, familyId, members); err != nil {
			return err
		}
	}

	return nil
}

func ensureFamilyMemberExists(
	ctx context.Context,
	familyRepository family.Repository,
	familyId uuid.UUID,
	members []uuid.UUID) error {
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
}
