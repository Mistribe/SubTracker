package command

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/domain/subscription"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
)

func ensureRelatedEntityExists(
	ctx context.Context,
	familyRepository ports.FamilyRepository,
	newSubscription subscription.Subscription,
) error {
	if newSubscription.Owner().Type() == types.FamilyOwnerType {
		familyId := newSubscription.Owner().FamilyId()

		var members []types.FamilyMemberID
		if newSubscription.FamilyUsers() != nil && newSubscription.FamilyUsers().Len() > 0 {
			members = newSubscription.FamilyUsers().Values()
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
	familyRepository ports.FamilyRepository,
	familyId types.FamilyID,
	members []types.FamilyMemberID) error {
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
