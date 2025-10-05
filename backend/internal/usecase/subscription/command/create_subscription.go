package command

import (
	"context"
	"time"

	"github.com/mistribe/subtracker/internal/domain/authorization"
	"github.com/mistribe/subtracker/internal/domain/billing"
	"github.com/mistribe/subtracker/internal/domain/currency"
	"github.com/mistribe/subtracker/internal/domain/subscription"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/langext/option"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type CreateSubscriptionCommand struct {
	SubscriptionID   option.Option[types.SubscriptionID]
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
	CreatedAt        option.Option[time.Time]
}

type CreateSubscriptionCommandHandler struct {
	subscriptionRepository ports.SubscriptionRepository
	authorization          ports.Authorization
	familyRepository       ports.FamilyRepository
	entitlement            ports.EntitlementResolver
}

func NewCreateSubscriptionCommandHandler(
	subscriptionRepository ports.SubscriptionRepository,
	authorization ports.Authorization,
	familyRepository ports.FamilyRepository,
	entitlement ports.EntitlementResolver) *CreateSubscriptionCommandHandler {
	return &CreateSubscriptionCommandHandler{
		subscriptionRepository: subscriptionRepository,
		familyRepository:       familyRepository,
		authorization:          authorization,
		entitlement:            entitlement,
	}
}

func (h CreateSubscriptionCommandHandler) Handle(
	ctx context.Context,
	cmd CreateSubscriptionCommand) result.Result[subscription.Subscription] {
	// Prepare subscription ID if provided
	var subscriptionID types.SubscriptionID
	if cmd.SubscriptionID != nil && cmd.SubscriptionID.IsSome() {
		subscriptionID = *cmd.SubscriptionID.Value()
	}

	// Determine creation timestamp
	var createdAt time.Time
	if cmd.CreatedAt != nil && cmd.CreatedAt.IsSome() {
		createdAt = *cmd.CreatedAt.Value()
	} else {
		createdAt = time.Now()
	}

	// Guard price creation: only create Price aggregate if amount provided (avoid nil amount panic)
	var price subscription.Price
	if cmd.Price != nil { // interface can be nil
		price = subscription.NewPrice(cmd.Price)
	}

	// Build the subscription aggregate early so authorization can run before repository lookups
	newSub := subscription.NewSubscription(
		subscriptionID,
		cmd.FriendlyName,
		cmd.FreeTrial,
		cmd.ProviderID,
		price,
		cmd.Owner,
		cmd.Payer,
		cmd.FamilyUsers,
		subscription.NewSubscriptionLabelRefs(cmd.Labels),
		cmd.StartDate,
		cmd.EndDate,
		cmd.Recurrency,
		cmd.CustomRecurrency,
		createdAt,
		createdAt,
	)

	// Authorization first to avoid leaking existence information
	if err := h.authorization.Can(ctx, authorization.PermissionWrite).For(newSub); err != nil {
		return result.Fail[subscription.Subscription](err)
	}

	// Duplication check only after permission granted
	if cmd.SubscriptionID != nil && cmd.SubscriptionID.IsSome() {
		sub, err := h.subscriptionRepository.GetById(ctx, subscriptionID)
		if err != nil {
			return result.Fail[subscription.Subscription](err)
		}
		if sub != nil {
			return result.Fail[subscription.Subscription](subscription.ErrSubscriptionAlreadyExists)
		}
	}

	// Quota check (ensure int64 literal for matcher compatibility)
	allowed, _, err := h.entitlement.CheckQuota(ctx, billing.FeatureIdActiveSubscriptionsCount, int64(1))
	if err != nil {
		return result.Fail[subscription.Subscription](err)
	}
	if !allowed {
		return result.Fail[subscription.Subscription](subscription.ErrActiveSubscriptionLimitReached)
	}

	return h.createSubscription(ctx, newSub)
}

func (h CreateSubscriptionCommandHandler) createSubscription(
	ctx context.Context,
	newSub subscription.Subscription) result.Result[subscription.Subscription] {
	if err := ensureRelatedEntityExists(ctx, h.familyRepository, newSub); err != nil {
		return result.Fail[subscription.Subscription](err)
	}

	if err := newSub.GetValidationErrors(); err != nil {
		return result.Fail[subscription.Subscription](err)
	}

	if err := h.subscriptionRepository.Save(ctx, newSub); err != nil {
		return result.Fail[subscription.Subscription](err)
	}

	return result.Success(newSub)
}
