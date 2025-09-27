package command

import (
	"context"
	"time"

	"github.com/mistribe/subtracker/internal/domain/authorization"
	"github.com/mistribe/subtracker/internal/domain/currency"
	"github.com/mistribe/subtracker/internal/domain/subscription"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/langext/option"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type UpdateSubscriptionCommand struct {
	SubscriptionID   types.SubscriptionID
	FriendlyName     *string
	FreeTrial        subscription.FreeTrial
	ProviderID       types.ProviderID
	Price            currency.Amount
	Owner            types.Owner
	Payer            subscription.Payer
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
}

func NewUpdateSubscriptionCommandHandler(
	subscriptionRepository ports.SubscriptionRepository,
	familyRepository ports.FamilyRepository,
	authorization ports.Authorization) *UpdateSubscriptionCommandHandler {
	return &UpdateSubscriptionCommandHandler{
		subscriptionRepository: subscriptionRepository,
		familyRepository:       familyRepository,
		authorization:          authorization,
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

	sub.SetFriendlyName(cmd.FriendlyName)
	sub.SetFreeTrial(cmd.FreeTrial)
	sub.SetPrice(cmd.Price)
	sub.SetOwner(cmd.Owner)
	sub.SetPayer(cmd.Payer)
	sub.SetFamilyUsers(cmd.FamilyUsers)
	sub.SetStartDate(cmd.StartDate)
	sub.SetEndDate(cmd.EndDate)
	sub.SetRecurrency(cmd.Recurrency)
	sub.SetCustomRecurrency(cmd.CustomRecurrency)
	sub.SetUpdatedAt(cmd.UpdatedAt.ValueOrDefault(time.Now()))

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
