package command_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/mistribe/subtracker/internal/domain/authorization"
	"github.com/mistribe/subtracker/internal/domain/billing"
	"github.com/mistribe/subtracker/internal/domain/subscription"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/subscription/command"
	"github.com/mistribe/subtracker/pkg/langext/option"
)

func newPersonalSubscription() subscription.Subscription {
	name := "Netflix"
	return subscription.NewSubscription(
		types.NewSubscriptionID(),
		&name,
		nil,
		types.ProviderID(uuid.Must(uuid.NewV7())),
		nil,
		types.NewPersonalOwner(types.UserID("userID-Test")),
		nil,
		nil,
		nil,
		time.Now().Add(-time.Hour),
		nil,
		subscription.MonthlyRecurrency,
		nil,
		time.Now(),
		time.Now(),
	)
}

func newFamilySubscriptionWithMembers(familyId uuid.UUID, members []uuid.UUID) subscription.Subscription {
	name := "Disney+"
	typedMembers := make([]types.FamilyMemberID, len(members))
	for i, m := range members {
		typedMembers[i] = types.FamilyMemberID(m)
	}
	return subscription.NewSubscription(
		types.NewSubscriptionID(),
		&name,
		nil,
		types.ProviderID(uuid.Must(uuid.NewV7())),
		nil,
		types.NewFamilyOwner(types.FamilyID(familyId)),
		nil,
		typedMembers,
		nil,
		time.Now().Add(-time.Hour),
		nil,
		subscription.MonthlyRecurrency,
		nil,
		time.Now(),
		time.Now(),
	)
}

func TestCreateSubscriptionCommandHandler_Handle(t *testing.T) {
	userId := "userID-Test"

	t.Run("returns fault when permission denied", func(t *testing.T) {
		subRepo := ports.NewMockSubscriptionRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		entitlement := ports.NewMockEntitlementResolver(t)
		perm := ports.NewMockPermissionRequest(t)
		authz.EXPECT().Can(t.Context(), authorization.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(authorization.ErrUnauthorized)

		h := command.NewCreateSubscriptionCommandHandler(subRepo, authz, familyRepo, entitlement)
		sub := newPersonalSubscription()
		cmd := command.CreateSubscriptionCommand{SubscriptionID: option.Some(sub.Id()), FriendlyName: sub.FriendlyName(), ProviderID: sub.ProviderId(), Owner: sub.Owner(), StartDate: sub.StartDate(), Recurrency: sub.Recurrency(), CreatedAt: option.Some(time.Now())}
		res := h.Handle(t.Context(), cmd)

		assert.True(t, res.IsFaulted())
	})

	t.Run("returns fault when repository GetById returns error", func(t *testing.T) {
		subRepo := ports.NewMockSubscriptionRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		entitlement := ports.NewMockEntitlementResolver(t)
		perm := ports.NewMockPermissionRequest(t)
		authz.EXPECT().Can(t.Context(), authorization.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)
		// Expect duplication check after authorization
		subRepo.EXPECT().GetById(t.Context(), mock.Anything).Return(nil, errors.New("db error"))

		h := command.NewCreateSubscriptionCommandHandler(subRepo, authz, familyRepo, entitlement)
		sub := newPersonalSubscription()
		cmd := command.CreateSubscriptionCommand{SubscriptionID: option.Some(sub.Id()), FriendlyName: sub.FriendlyName(), ProviderID: sub.ProviderId(), Owner: sub.Owner(), StartDate: sub.StartDate(), Recurrency: sub.Recurrency(), CreatedAt: option.Some(time.Now())}
		res := h.Handle(t.Context(), cmd)

		assert.True(t, res.IsFaulted())
	})

	t.Run("returns fault when subscription already exists", func(t *testing.T) {
		subRepo := ports.NewMockSubscriptionRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		entitlement := ports.NewMockEntitlementResolver(t)
		perm := ports.NewMockPermissionRequest(t)
		authz.EXPECT().Can(t.Context(), authorization.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)

		existing := newPersonalSubscription()
		subRepo.EXPECT().GetById(t.Context(), existing.Id()).Return(existing, nil)

		h := command.NewCreateSubscriptionCommandHandler(subRepo, authz, familyRepo, entitlement)
		cmd := command.CreateSubscriptionCommand{SubscriptionID: option.Some(existing.Id()), FriendlyName: existing.FriendlyName(), ProviderID: existing.ProviderId(), Owner: existing.Owner(), StartDate: existing.StartDate(), Recurrency: existing.Recurrency(), CreatedAt: option.Some(time.Now())}
		res := h.Handle(t.Context(), cmd)

		assert.True(t, res.IsFaulted())
	})

	// removed differing existing test since behavior changed to simple already-exists fault

	t.Run("returns fault when entitlement quota check fails", func(t *testing.T) {
		subRepo := ports.NewMockSubscriptionRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		entitlement := ports.NewMockEntitlementResolver(t)
		perm := ports.NewMockPermissionRequest(t)
		authz.EXPECT().Can(t.Context(), authorization.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)
		// duplication check (ID provided)
		subRepo.EXPECT().GetById(t.Context(), mock.Anything).Return(nil, nil)
		entitlement.EXPECT().CheckQuota(t.Context(), billing.FeatureIdActiveSubscriptionsCount, int64(1)).Return(false, billing.EffectiveEntitlement{}, errors.New("limit"))

		h := command.NewCreateSubscriptionCommandHandler(subRepo, authz, familyRepo, entitlement)
		sub := newPersonalSubscription()
		cmd := command.CreateSubscriptionCommand{SubscriptionID: option.Some(sub.Id()), FriendlyName: sub.FriendlyName(), ProviderID: sub.ProviderId(), Owner: sub.Owner(), StartDate: sub.StartDate(), Recurrency: sub.Recurrency(), CreatedAt: option.Some(time.Now())}
		res := h.Handle(t.Context(), cmd)

		assert.True(t, res.IsFaulted())
	})

	t.Run("returns fault when related family member does not exist", func(t *testing.T) {
		subRepo := ports.NewMockSubscriptionRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		entitlement := ports.NewMockEntitlementResolver(t)
		perm := ports.NewMockPermissionRequest(t)
		authz.EXPECT().Can(t.Context(), authorization.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)
		subRepo.EXPECT().GetById(t.Context(), mock.Anything).Return(nil, nil)
		entitlement.EXPECT().CheckQuota(t.Context(), billing.FeatureIdActiveSubscriptionsCount, int64(1)).Return(true, billing.EffectiveEntitlement{}, nil)

		famID := uuid.Must(uuid.NewV7())
		member := uuid.Must(uuid.NewV7())
		// expect MemberExists check and return false
		familyRepo.EXPECT().MemberExists(t.Context(), types.FamilyID(famID), []types.FamilyMemberID{types.FamilyMemberID(member)}).Return(false, nil)

		h := command.NewCreateSubscriptionCommandHandler(subRepo, authz, familyRepo, entitlement)
		sub := newFamilySubscriptionWithMembers(famID, []uuid.UUID{member})
		cmd := command.CreateSubscriptionCommand{SubscriptionID: option.Some(sub.Id()), FriendlyName: sub.FriendlyName(), ProviderID: sub.ProviderId(), Owner: sub.Owner(), StartDate: sub.StartDate(), Recurrency: sub.Recurrency(), CreatedAt: option.Some(time.Now()), FamilyUsers: sub.FamilyUsers().Values()}
		res := h.Handle(t.Context(), cmd)

		assert.True(t, res.IsFaulted())
	})

	t.Run("returns fault when family repository errors during MemberExists", func(t *testing.T) {
		subRepo := ports.NewMockSubscriptionRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		entitlement := ports.NewMockEntitlementResolver(t)
		perm := ports.NewMockPermissionRequest(t)
		authz.EXPECT().Can(t.Context(), authorization.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)
		subRepo.EXPECT().GetById(t.Context(), mock.Anything).Return(nil, nil)
		entitlement.EXPECT().CheckQuota(t.Context(), billing.FeatureIdActiveSubscriptionsCount, int64(1)).Return(true, billing.EffectiveEntitlement{}, nil)

		famID := uuid.Must(uuid.NewV7())
		member := uuid.Must(uuid.NewV7())
		familyRepo.EXPECT().MemberExists(t.Context(), types.FamilyID(famID), mock.Anything).Return(false, errors.New("repo error"))

		h := command.NewCreateSubscriptionCommandHandler(subRepo, authz, familyRepo, entitlement)
		sub := newFamilySubscriptionWithMembers(famID, []uuid.UUID{member})
		cmd := command.CreateSubscriptionCommand{SubscriptionID: option.Some(sub.Id()), FriendlyName: sub.FriendlyName(), ProviderID: sub.ProviderId(), Owner: sub.Owner(), StartDate: sub.StartDate(), Recurrency: sub.Recurrency(), CreatedAt: option.Some(time.Now()), FamilyUsers: sub.FamilyUsers().Values()}
		res := h.Handle(t.Context(), cmd)

		assert.True(t, res.IsFaulted())
	})

	t.Run("returns fault when validation fails", func(t *testing.T) {
		subRepo := ports.NewMockSubscriptionRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		entitlement := ports.NewMockEntitlementResolver(t)
		perm := ports.NewMockPermissionRequest(t)
		authz.EXPECT().Can(t.Context(), authorization.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)
		// No SubscriptionID provided so no GetById expectation
		entitlement.EXPECT().CheckQuota(t.Context(), billing.FeatureIdActiveSubscriptionsCount, int64(1)).Return(true, billing.EffectiveEntitlement{}, nil)

		// invalid: Custom recurrency without amount
		name := "Invalid"
		h := command.NewCreateSubscriptionCommandHandler(subRepo, authz, familyRepo, entitlement)
		cmd := command.CreateSubscriptionCommand{FriendlyName: &name, ProviderID: types.ProviderID(uuid.Must(uuid.NewV7())), Owner: types.NewPersonalOwner(types.UserID(userId)), StartDate: time.Now().Add(-time.Hour), Recurrency: subscription.CustomRecurrency}
		res := h.Handle(t.Context(), cmd)

		assert.True(t, res.IsFaulted())
	})

	t.Run("returns fault when save fails", func(t *testing.T) {
		subRepo := ports.NewMockSubscriptionRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		entitlement := ports.NewMockEntitlementResolver(t)
		perm := ports.NewMockPermissionRequest(t)
		authz.EXPECT().Can(t.Context(), authorization.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)
		// duplication check (ID provided)
		subRepo.EXPECT().GetById(t.Context(), mock.Anything).Return(nil, nil)
		entitlement.EXPECT().CheckQuota(t.Context(), billing.FeatureIdActiveSubscriptionsCount, int64(1)).Return(true, billing.EffectiveEntitlement{}, nil)

		sub := newPersonalSubscription()
		subRepo.EXPECT().Save(t.Context(),
			mock.MatchedBy(func(entities []subscription.Subscription) bool { return len(entities) == 1 })).Return(errors.New("save failed"))

		h := command.NewCreateSubscriptionCommandHandler(subRepo, authz, familyRepo, entitlement)
		cmd := command.CreateSubscriptionCommand{SubscriptionID: option.Some(sub.Id()), FriendlyName: sub.FriendlyName(), ProviderID: sub.ProviderId(), Owner: sub.Owner(), StartDate: sub.StartDate(), Recurrency: sub.Recurrency(), CreatedAt: option.Some(time.Now())}
		res := h.Handle(t.Context(), cmd)

		assert.True(t, res.IsFaulted())
	})

	t.Run("returns success when created and saved", func(t *testing.T) {
		subRepo := ports.NewMockSubscriptionRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		entitlement := ports.NewMockEntitlementResolver(t)
		perm := ports.NewMockPermissionRequest(t)
		authz.EXPECT().Can(t.Context(), authorization.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)
		// duplication check (ID provided)
		subRepo.EXPECT().GetById(t.Context(), mock.Anything).Return(nil, nil)
		entitlement.EXPECT().CheckQuota(t.Context(), billing.FeatureIdActiveSubscriptionsCount, int64(1)).Return(true, billing.EffectiveEntitlement{}, nil)

		sub := newPersonalSubscription()
		subRepo.EXPECT().Save(t.Context(),
			mock.MatchedBy(func(entities []subscription.Subscription) bool { return len(entities) == 1 })).Return(nil)

		h := command.NewCreateSubscriptionCommandHandler(subRepo, authz, familyRepo, entitlement)
		cmd := command.CreateSubscriptionCommand{SubscriptionID: option.Some(sub.Id()), FriendlyName: sub.FriendlyName(), ProviderID: sub.ProviderId(), Owner: sub.Owner(), StartDate: sub.StartDate(), Recurrency: sub.Recurrency(), CreatedAt: option.Some(time.Now())}
		res := h.Handle(t.Context(), cmd)

		assert.True(t, res.IsSuccess())
	})
}
