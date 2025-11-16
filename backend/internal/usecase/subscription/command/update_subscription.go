package command

import (
	"context"
	"time"

	"github.com/mistribe/subtracker/internal/domain/authorization"
	"github.com/mistribe/subtracker/internal/domain/currency"
	"github.com/mistribe/subtracker/internal/domain/subscription"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/shared"
	"github.com/mistribe/subtracker/pkg/langext/option"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type UpdateSubscriptionCommand struct {
	SubscriptionID   types.SubscriptionID
	FriendlyName     *string
	FreeTrial        subscription.FreeTrial
	ProviderID       types.ProviderID
	Price            currency.Amount
	Owner            types.OwnerType
	PayerType        *subscription.PayerType
	PayerMemberId    *types.FamilyMemberID
	FamilyUsers      []types.FamilyMemberID
	Labels           []types.LabelID
	StartDate        time.Time
	EndDate          *time.Time
	Recurrency       subscription.RecurrencyType
	CustomRecurrency *int32
	UpdatedAt        option.Option[time.Time]
}

type UpdateSubscriptionCommandHandler struct {
	subscriptionRepository ports.SubscriptionRepository
	familyRepository       ports.FamilyRepository
	authorization          ports.Authorization
	ownerFactory           shared.OwnerFactory
}

func NewUpdateSubscriptionCommandHandler(
	subscriptionRepository ports.SubscriptionRepository,
	familyRepository ports.FamilyRepository,
	ownerFactory shared.OwnerFactory,
	authorization ports.Authorization) *UpdateSubscriptionCommandHandler {
	return &UpdateSubscriptionCommandHandler{
		subscriptionRepository: subscriptionRepository,
		familyRepository:       familyRepository,
		authorization:          authorization,
		ownerFactory:           ownerFactory,
	}
}

func (h UpdateSubscriptionCommandHandler) Handle(
	ctx context.Context,
	cmd UpdateSubscriptionCommand) result.Result[subscription.Subscription] {
	sub, err := h.subscriptionRepository.GetById(ctx, cmd.SubscriptionID)
	if err != nil {
		return result.Fail[subscription.Subscription](err)
	}
	if sub == nil {
		return result.Fail[subscription.Subscription](subscription.ErrSubscriptionNotFound)
	}

	err = h.authorization.Can(ctx, authorization.PermissionWrite).For(sub)
	if err != nil {
		return result.Fail[subscription.Subscription](err)
	}

	return h.updateSubscription(ctx, cmd, sub)
}

func (h UpdateSubscriptionCommandHandler) updateSubscription(
	ctx context.Context,
	cmd UpdateSubscriptionCommand,
	sub subscription.Subscription) result.Result[subscription.Subscription] {

	if cmd.FriendlyName != nil {
		sub.SetFriendlyName(cmd.FriendlyName)
	}
	if cmd.FreeTrial != nil {
		sub.SetFreeTrial(cmd.FreeTrial)
	}
	if sub.Price() != nil && cmd.Price != nil {
		sub.SetPrice(cmd.Price)
	}
	if cmd.Owner != sub.Owner().Type() {
		owner, err := h.ownerFactory.Resolve(ctx, cmd.Owner)
		if err != nil {
			return result.Fail[subscription.Subscription](err)
		}
		sub.SetOwner(owner)
	}
	if cmd.PayerType == nil && sub.Payer() != nil {
		sub.SetPayer(nil)
	} else if cmd.PayerType != nil &&
		(*cmd.PayerType != sub.Payer().Type() ||
			cmd.PayerMemberId != nil && *cmd.PayerMemberId != sub.Payer().MemberId()) {
		payer, err := subscription.NewPayerFromOwner(sub.Owner(), *cmd.PayerType, cmd.PayerMemberId)
		if err != nil {
			return result.Fail[subscription.Subscription](err)
		}
		sub.SetPayer(payer)
	}
	if cmd.FamilyUsers != nil {
		sub.SetFamilyUsers(cmd.FamilyUsers)
	}
	if !cmd.StartDate.IsZero() {
		sub.SetStartDate(cmd.StartDate)
	}
	if cmd.EndDate != nil {
		sub.SetEndDate(cmd.EndDate)
	}
	var zeroRecurrency subscription.RecurrencyType
	if cmd.Recurrency != zeroRecurrency {
		sub.SetRecurrency(cmd.Recurrency)
	}
	if cmd.CustomRecurrency != nil {
		sub.SetCustomRecurrency(cmd.CustomRecurrency)
	}
	var updatedAt time.Time
	if cmd.UpdatedAt != nil && cmd.UpdatedAt.IsSome() {
		updatedAt = *cmd.UpdatedAt.Value()
	} else {
		updatedAt = time.Now()
	}
	sub.SetUpdatedAt(updatedAt)

	if err := ensureRelatedEntityExists(ctx, h.familyRepository, sub); err != nil {
		return result.Fail[subscription.Subscription](err)
	}

	if err := sub.GetValidationErrors(); err != nil {
		return result.Fail[subscription.Subscription](err)
	}

	if err := h.subscriptionRepository.Save(ctx, sub); err != nil {
		return result.Fail[subscription.Subscription](err)
	}

	return result.Success(sub)
}
