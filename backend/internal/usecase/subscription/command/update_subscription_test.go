package command_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/mistribe/subtracker/internal/domain/user"

	"github.com/mistribe/subtracker/internal/domain/subscription"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/subscription/command"
)

func TestUpdateSubscriptionCommandHandler_Handle(t *testing.T) {
	userId := "userID-Test"

	t.Run("returns fault when repository GetById returns error", func(t *testing.T) {
		subRepo := ports.NewMockSubscriptionRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)

		// command with any id
		cmdSub := newPersonalSubscription()
		subRepo.EXPECT().GetById(t.Context(), cmdSub.Id()).Return(nil, errors.New("db error"))

		h := command.NewUpdateSubscriptionCommandHandler(subRepo, familyRepo, authz)
		cmd := command.UpdateSubscriptionCommand{Subscription: cmdSub}
		res := h.Handle(t.Context(), cmd)

		assert.True(t, res.IsFaulted())
	})

	t.Run("returns fault when subscription not found", func(t *testing.T) {
		subRepo := ports.NewMockSubscriptionRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)

		cmdSub := newPersonalSubscription()
		subRepo.EXPECT().GetById(t.Context(), cmdSub.Id()).Return(nil, nil)

		h := command.NewUpdateSubscriptionCommandHandler(subRepo, familyRepo, authz)
		cmd := command.UpdateSubscriptionCommand{Subscription: cmdSub}
		res := h.Handle(t.Context(), cmd)

		assert.True(t, res.IsFaulted())
	})

	t.Run("returns fault when permission denied", func(t *testing.T) {
		subRepo := ports.NewMockSubscriptionRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)

		// existing subscription (fetched)
		existing := newPersonalSubscription()
		cmdSub := subscription.NewSubscription(
			existing.Id(),
			existing.FriendlyName(),
			existing.FreeTrial(),
			existing.ProviderId(),
			existing.PlanId(),
			existing.PriceId(),
			existing.CustomPrice(),
			existing.Owner(),
			existing.Payer(),
			existing.ServiceUsers().Values(),
			nil,
			existing.StartDate(),
			existing.EndDate(),
			existing.Recurrency(),
			existing.CustomRecurrency(),
			time.Now(),
			time.Now(),
		)

		subRepo.EXPECT().GetById(t.Context(), cmdSub.Id()).Return(existing, nil)
		authz.EXPECT().Can(t.Context(), user.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(user.ErrUnauthorized)

		h := command.NewUpdateSubscriptionCommandHandler(subRepo, familyRepo, authz)
		cmd := command.UpdateSubscriptionCommand{Subscription: cmdSub}
		res := h.Handle(t.Context(), cmd)

		assert.True(t, res.IsFaulted())
	})

	t.Run("returns fault when related family member does not exist", func(t *testing.T) {
		subRepo := ports.NewMockSubscriptionRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)

		// existing subscription (any non-nil)
		existing := newPersonalSubscription()
		// cmd subscription owned by family with members
		famID := uuid.Must(uuid.NewV7())
		member := uuid.Must(uuid.NewV7())
		cmdSub := newFamilySubscriptionWithMembers(famID, []uuid.UUID{member})

		subRepo.EXPECT().GetById(t.Context(), cmdSub.Id()).Return(existing, nil)
		authz.EXPECT().Can(t.Context(), user.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)
		familyRepo.EXPECT().MemberExists(t.Context(), famID,
			mock.MatchedBy(func(members []uuid.UUID) bool { return len(members) == 1 && members[0] == member }),
		).Return(false, nil)

		h := command.NewUpdateSubscriptionCommandHandler(subRepo, familyRepo, authz)
		cmd := command.UpdateSubscriptionCommand{Subscription: cmdSub}
		res := h.Handle(t.Context(), cmd)

		assert.True(t, res.IsFaulted())
	})

	t.Run("returns fault when family repository errors during MemberExists", func(t *testing.T) {
		subRepo := ports.NewMockSubscriptionRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)

		existing := newPersonalSubscription()
		famID := uuid.Must(uuid.NewV7())
		member := uuid.Must(uuid.NewV7())
		cmdSub := newFamilySubscriptionWithMembers(famID, []uuid.UUID{member})

		subRepo.EXPECT().GetById(t.Context(), cmdSub.Id()).Return(existing, nil)
		authz.EXPECT().Can(t.Context(), user.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)
		familyRepo.EXPECT().MemberExists(t.Context(), famID, mock.Anything).Return(false, errors.New("repo error"))

		h := command.NewUpdateSubscriptionCommandHandler(subRepo, familyRepo, authz)
		cmd := command.UpdateSubscriptionCommand{Subscription: cmdSub}
		res := h.Handle(t.Context(), cmd)

		assert.True(t, res.IsFaulted())
	})

	t.Run("returns fault when validation fails", func(t *testing.T) {
		subRepo := ports.NewMockSubscriptionRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)

		// existing subscription (personal owner)
		existing := newPersonalSubscription()
		// invalid cmd: Custom recurrency without amount
		name := "Invalid"
		invalid := subscription.NewSubscription(
			existing.Id(),
			&name,
			nil,
			existing.ProviderId(),
			nil,
			nil,
			nil,
			types.NewPersonalOwner(userId),
			nil,
			nil,
			nil,
			time.Now().Add(-time.Hour),
			nil,
			subscription.CustomRecurrency,
			nil,
			time.Now(),
			time.Now(),
		)

		subRepo.EXPECT().GetById(t.Context(), invalid.Id()).Return(existing, nil)
		authz.EXPECT().Can(t.Context(), user.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)

		h := command.NewUpdateSubscriptionCommandHandler(subRepo, familyRepo, authz)
		cmd := command.UpdateSubscriptionCommand{Subscription: invalid}
		res := h.Handle(t.Context(), cmd)

		assert.True(t, res.IsFaulted())
	})

	t.Run("returns fault when save fails", func(t *testing.T) {
		subRepo := ports.NewMockSubscriptionRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)

		existing := newPersonalSubscription()
		name := "Updated"
		updated := subscription.NewSubscription(
			existing.Id(),
			&name,
			nil,
			existing.ProviderId(),
			nil,
			nil,
			nil,
			types.NewPersonalOwner(userId),
			nil,
			nil,
			nil,
			time.Now().Add(-2*time.Hour),
			nil,
			subscription.MonthlyRecurrency,
			nil,
			time.Now(),
			time.Now(),
		)

		subRepo.EXPECT().GetById(t.Context(), updated.Id()).Return(existing, nil)
		authz.EXPECT().Can(t.Context(), user.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)
		subRepo.EXPECT().Save(t.Context(), mock.MatchedBy(func(entities []subscription.Subscription) bool {
			return len(entities) == 1 && entities[0] == existing
		})).Return(errors.New("save failed"))

		h := command.NewUpdateSubscriptionCommandHandler(subRepo, familyRepo, authz)
		cmd := command.UpdateSubscriptionCommand{Subscription: updated}
		res := h.Handle(t.Context(), cmd)

		assert.True(t, res.IsFaulted())
	})

	t.Run("returns success when updated and saved", func(t *testing.T) {
		subRepo := ports.NewMockSubscriptionRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)

		existing := newPersonalSubscription()
		name := "Updated"
		updated := subscription.NewSubscription(
			existing.Id(),
			&name,
			subscription.NewFreeTrial(time.Now().Add(-48*time.Hour), time.Now().Add(-24*time.Hour)),
			existing.ProviderId(),
			nil,
			nil,
			nil,
			types.NewPersonalOwner(userId),
			nil,
			nil,
			nil,
			time.Now().Add(-2*time.Hour),
			nil,
			subscription.MonthlyRecurrency,
			nil,
			time.Now(),
			time.Now(),
		)

		subRepo.EXPECT().GetById(t.Context(), updated.Id()).Return(existing, nil)
		authz.EXPECT().Can(t.Context(), user.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)
		subRepo.EXPECT().Save(t.Context(), mock.MatchedBy(func(entities []subscription.Subscription) bool {
			return len(entities) == 1 && entities[0] == existing
		})).Return(nil)

		h := command.NewUpdateSubscriptionCommandHandler(subRepo, familyRepo, authz)
		cmd := command.UpdateSubscriptionCommand{Subscription: updated}
		res := h.Handle(t.Context(), cmd)

		assert.True(t, res.IsSuccess())
	})

	t.Run("updates fields accordingly", func(t *testing.T) {
		subRepo := ports.NewMockSubscriptionRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)

		existing := newPersonalSubscription()
		// new values to be set on existing
		name := "New Name"
		start := time.Now().Add(-72 * time.Hour)
		end := timePtr(time.Now().Add(72 * time.Hour))
		customRec := int32(5)
		updated := subscription.NewSubscription(
			existing.Id(),
			&name,
			subscription.NewFreeTrial(time.Now().Add(-10*time.Hour), time.Now().Add(-9*time.Hour)),
			existing.ProviderId(),
			nil,
			nil,
			nil,
			types.NewPersonalOwner(userId),
			nil,
			nil,
			nil,
			start,
			end,
			subscription.CustomRecurrency,
			&customRec,
			time.Now(),
			time.Now(),
		)

		subRepo.EXPECT().GetById(t.Context(), updated.Id()).Return(existing, nil)
		authz.EXPECT().Can(t.Context(), user.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)
		subRepo.EXPECT().Save(t.Context(), mock.Anything).Return(nil)

		h := command.NewUpdateSubscriptionCommandHandler(subRepo, familyRepo, authz)
		cmd := command.UpdateSubscriptionCommand{Subscription: updated}
		_ = h.Handle(t.Context(), cmd)

		// Assert fields were updated on the existing entity
		assert.Equal(t, name, *existing.FriendlyName())
		assert.Equal(t, start, existing.StartDate())
		assert.Equal(t, end, existing.EndDate())
		assert.Equal(t, subscription.CustomRecurrency, existing.Recurrency())
		assert.Equal(t, customRec, *existing.CustomRecurrency())
	})
}

func timePtr(ti time.Time) *time.Time { return &ti }
