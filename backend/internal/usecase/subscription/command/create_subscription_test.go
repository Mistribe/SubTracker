package command_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/mistribe/subtracker/internal/domain/auth"
	"github.com/mistribe/subtracker/internal/domain/subscription"
	"github.com/mistribe/subtracker/internal/domain/user"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/subscription/command"
)

func newPersonalSubscription() subscription.Subscription {
	name := "Netflix"
	return subscription.NewSubscription(
		uuid.Must(uuid.NewV7()),
		&name,
		nil, // freeTrial
		uuid.Must(uuid.NewV7()),
		nil, // planId
		nil, // priceId
		nil, // customPrice
		auth.NewPersonalOwner("userID-Test"),
		nil, // payer
		nil, // serviceUsers
		nil, // labels
		time.Now().Add(-time.Hour),
		nil, // endDate
		subscription.MonthlyRecurrency,
		nil, // customRecurrency
		time.Now(),
		time.Now(),
	)
}

func newFamilySubscriptionWithMembers(familyId uuid.UUID, members []uuid.UUID) subscription.Subscription {
	name := "Disney+"
	return subscription.NewSubscription(
		uuid.Must(uuid.NewV7()),
		&name,
		nil,
		uuid.Must(uuid.NewV7()),
		nil,
		nil,
		nil,
		auth.NewFamilyOwner(familyId),
		nil,
		members,
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
		perm := ports.NewMockPermissionRequest(t)
		authz.EXPECT().Can(t.Context(), user.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(user.ErrUnauthorized)

		h := command.NewCreateSubscriptionCommandHandler(subRepo, authz, familyRepo)
		cmd := command.CreateSubscriptionCommand{Subscription: newPersonalSubscription()}
		res := h.Handle(t.Context(), cmd)

		assert.True(t, res.IsFaulted())
	})

	t.Run("returns fault when repository GetById returns error", func(t *testing.T) {
		subRepo := ports.NewMockSubscriptionRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)
		authz.EXPECT().Can(t.Context(), user.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)
		subRepo.EXPECT().GetById(t.Context(), mock.Anything).Return(nil, errors.New("db error"))

		h := command.NewCreateSubscriptionCommandHandler(subRepo, authz, familyRepo)
		cmd := command.CreateSubscriptionCommand{Subscription: newPersonalSubscription()}
		res := h.Handle(t.Context(), cmd)

		assert.True(t, res.IsFaulted())
	})

	t.Run("returns success when subscription already exists and is identical", func(t *testing.T) {
		subRepo := ports.NewMockSubscriptionRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)
		authz.EXPECT().Can(t.Context(), user.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)

		existing := newPersonalSubscription()
		subRepo.EXPECT().GetById(t.Context(), existing.Id()).Return(existing, nil)

		h := command.NewCreateSubscriptionCommandHandler(subRepo, authz, familyRepo)
		cmd := command.CreateSubscriptionCommand{Subscription: existing}
		res := h.Handle(t.Context(), cmd)

		assert.True(t, res.IsSuccess())
	})

	t.Run("returns fault when subscription already exists but differs", func(t *testing.T) {
		subRepo := ports.NewMockSubscriptionRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)
		authz.EXPECT().Can(t.Context(), user.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)

		existing := newPersonalSubscription()
		subRepo.EXPECT().GetById(t.Context(), mock.Anything).Return(existing, nil)

		h := command.NewCreateSubscriptionCommandHandler(subRepo, authz, familyRepo)
		// make different subscription (different FriendlyName)
		name := "Different"
		newSub := subscription.NewSubscription(
			uuid.Must(uuid.NewV7()),
			&name,
			nil,
			uuid.Must(uuid.NewV7()),
			nil,
			nil,
			nil,
			auth.NewPersonalOwner(userId),
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
		cmd := command.CreateSubscriptionCommand{Subscription: newSub}
		res := h.Handle(t.Context(), cmd)

		assert.True(t, res.IsFaulted())
	})

	t.Run("returns fault when EnsureWithinLimit fails", func(t *testing.T) {
		subRepo := ports.NewMockSubscriptionRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)
		authz.EXPECT().Can(t.Context(), user.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)
		subRepo.EXPECT().GetById(t.Context(), mock.Anything).Return(nil, nil)
		authz.EXPECT().EnsureWithinLimit(t.Context(), user.FeatureActiveSubscriptions, 1).Return(errors.New("limit"))

		h := command.NewCreateSubscriptionCommandHandler(subRepo, authz, familyRepo)
		cmd := command.CreateSubscriptionCommand{Subscription: newPersonalSubscription()}
		res := h.Handle(t.Context(), cmd)

		assert.True(t, res.IsFaulted())
	})

	t.Run("returns fault when related family member does not exist", func(t *testing.T) {
		subRepo := ports.NewMockSubscriptionRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)
		authz.EXPECT().Can(t.Context(), user.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)
		subRepo.EXPECT().GetById(t.Context(), mock.Anything).Return(nil, nil)
		authz.EXPECT().EnsureWithinLimit(t.Context(), user.FeatureActiveSubscriptions, 1).Return(nil)

		famID := uuid.Must(uuid.NewV7())
		member := uuid.Must(uuid.NewV7())
		// expect MemberExists check and return false
		familyRepo.EXPECT().MemberExists(t.Context(), famID,
			mock.MatchedBy(func(members []uuid.UUID) bool { return len(members) == 1 && members[0] == member })).Return(false,
			nil)

		h := command.NewCreateSubscriptionCommandHandler(subRepo, authz, familyRepo)
		sub := newFamilySubscriptionWithMembers(famID, []uuid.UUID{member})
		cmd := command.CreateSubscriptionCommand{Subscription: sub}
		res := h.Handle(t.Context(), cmd)

		assert.True(t, res.IsFaulted())
	})

	t.Run("returns fault when family repository errors during MemberExists", func(t *testing.T) {
		subRepo := ports.NewMockSubscriptionRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)
		authz.EXPECT().Can(t.Context(), user.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)
		subRepo.EXPECT().GetById(t.Context(), mock.Anything).Return(nil, nil)
		authz.EXPECT().EnsureWithinLimit(t.Context(), user.FeatureActiveSubscriptions, 1).Return(nil)

		famID := uuid.Must(uuid.NewV7())
		member := uuid.Must(uuid.NewV7())
		familyRepo.EXPECT().MemberExists(t.Context(), famID, mock.Anything).Return(false, errors.New("repo error"))

		h := command.NewCreateSubscriptionCommandHandler(subRepo, authz, familyRepo)
		sub := newFamilySubscriptionWithMembers(famID, []uuid.UUID{member})
		cmd := command.CreateSubscriptionCommand{Subscription: sub}
		res := h.Handle(t.Context(), cmd)

		assert.True(t, res.IsFaulted())
	})

	t.Run("returns fault when validation fails", func(t *testing.T) {
		subRepo := ports.NewMockSubscriptionRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)
		authz.EXPECT().Can(t.Context(), user.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)
		subRepo.EXPECT().GetById(t.Context(), mock.Anything).Return(nil, nil)
		authz.EXPECT().EnsureWithinLimit(t.Context(), user.FeatureActiveSubscriptions, 1).Return(nil)

		// invalid: Custom recurrency without amount
		name := "Invalid"
		invalid := subscription.NewSubscription(
			uuid.Must(uuid.NewV7()),
			&name,
			nil,
			uuid.Must(uuid.NewV7()),
			nil,
			nil,
			nil,
			auth.NewPersonalOwner(userId),
			nil,
			nil,
			nil,
			time.Now().Add(-time.Hour),
			nil,
			subscription.CustomRecurrency,
			nil, // missing customRecurrency
			time.Now(),
			time.Now(),
		)

		h := command.NewCreateSubscriptionCommandHandler(subRepo, authz, familyRepo)
		cmd := command.CreateSubscriptionCommand{Subscription: invalid}
		res := h.Handle(t.Context(), cmd)

		assert.True(t, res.IsFaulted())
	})

	t.Run("returns fault when save fails", func(t *testing.T) {
		subRepo := ports.NewMockSubscriptionRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)
		authz.EXPECT().Can(t.Context(), user.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)
		subRepo.EXPECT().GetById(t.Context(), mock.Anything).Return(nil, nil)
		authz.EXPECT().EnsureWithinLimit(t.Context(), user.FeatureActiveSubscriptions, 1).Return(nil)

		sub := newPersonalSubscription()
		subRepo.EXPECT().Save(t.Context(),
			mock.MatchedBy(func(entities []subscription.Subscription) bool { return len(entities) == 1 && entities[0] == sub })).Return(errors.New("save failed"))

		h := command.NewCreateSubscriptionCommandHandler(subRepo, authz, familyRepo)
		cmd := command.CreateSubscriptionCommand{Subscription: sub}
		res := h.Handle(t.Context(), cmd)

		assert.True(t, res.IsFaulted())
	})

	t.Run("returns success when created and saved", func(t *testing.T) {
		subRepo := ports.NewMockSubscriptionRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)
		authz.EXPECT().Can(t.Context(), user.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)
		subRepo.EXPECT().GetById(t.Context(), mock.Anything).Return(nil, nil)
		authz.EXPECT().EnsureWithinLimit(t.Context(), user.FeatureActiveSubscriptions, 1).Return(nil)

		sub := newPersonalSubscription()
		subRepo.EXPECT().Save(t.Context(),
			mock.MatchedBy(func(entities []subscription.Subscription) bool { return len(entities) == 1 && entities[0] == sub })).Return(nil)

		h := command.NewCreateSubscriptionCommandHandler(subRepo, authz, familyRepo)
		cmd := command.CreateSubscriptionCommand{Subscription: sub}
		res := h.Handle(t.Context(), cmd)

		assert.True(t, res.IsSuccess())
	})
}
