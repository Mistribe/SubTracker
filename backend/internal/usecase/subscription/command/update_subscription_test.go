package command_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/mistribe/subtracker/internal/domain/authorization"
	"github.com/mistribe/subtracker/internal/domain/subscription"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/subscription/command"
)

func TestUpdateSubscriptionCommandHandler_Handle(t *testing.T) {
	// user identifier used in ownership tests
	_ = "userID-Test"

	t.Run("returns fault when repository GetById returns error", func(t *testing.T) {
		subRepo := ports.NewMockSubscriptionRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)

		existing := newPersonalSubscription()
		subRepo.EXPECT().GetById(t.Context(), existing.Id()).Return(nil, errors.New("db error"))

		h := command.NewUpdateSubscriptionCommandHandler(subRepo, familyRepo, authz)
		cmd := command.UpdateSubscriptionCommand{SubscriptionID: existing.Id()}
		res := h.Handle(t.Context(), cmd)

		assert.True(t, res.IsFaulted())
	})

	t.Run("returns fault when subscription not found", func(t *testing.T) {
		subRepo := ports.NewMockSubscriptionRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)

		existing := newPersonalSubscription()
		subRepo.EXPECT().GetById(t.Context(), existing.Id()).Return(nil, nil)

		h := command.NewUpdateSubscriptionCommandHandler(subRepo, familyRepo, authz)
		cmd := command.UpdateSubscriptionCommand{SubscriptionID: existing.Id()}
		res := h.Handle(t.Context(), cmd)

		assert.True(t, res.IsFaulted())
	})

	t.Run("returns fault when permission denied", func(t *testing.T) {
		subRepo := ports.NewMockSubscriptionRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)

		existing := newPersonalSubscription()
		subRepo.EXPECT().GetById(t.Context(), existing.Id()).Return(existing, nil)
		authz.EXPECT().Can(t.Context(), authorization.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(authorization.ErrUnauthorized)

		h := command.NewUpdateSubscriptionCommandHandler(subRepo, familyRepo, authz)
		cmd := command.UpdateSubscriptionCommand{SubscriptionID: existing.Id()}
		res := h.Handle(t.Context(), cmd)

		assert.True(t, res.IsFaulted())
	})

	t.Run("returns fault when related family member does not exist", func(t *testing.T) {
		subRepo := ports.NewMockSubscriptionRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)

		existing := newPersonalSubscription()
		famID := uuid.Must(uuid.NewV7())
		member := uuid.Must(uuid.NewV7())
		cmdSub := newFamilySubscriptionWithMembers(famID, []uuid.UUID{member})
		subRepo.EXPECT().GetById(t.Context(), cmdSub.Id()).Return(existing, nil)
		authz.EXPECT().Can(t.Context(), authorization.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)
		familyRepo.EXPECT().MemberExists(t.Context(), types.FamilyID(famID), []types.FamilyMemberID{types.FamilyMemberID(member)}).Return(false, nil)

		h := command.NewUpdateSubscriptionCommandHandler(subRepo, familyRepo, authz)
		cmd := command.UpdateSubscriptionCommand{SubscriptionID: cmdSub.Id(), Owner: cmdSub.Owner(), FamilyUsers: cmdSub.FamilyUsers().Values()}
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
		authz.EXPECT().Can(t.Context(), authorization.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)
		familyRepo.EXPECT().MemberExists(t.Context(), types.FamilyID(famID), mock.Anything).Return(false, errors.New("repo error"))

		h := command.NewUpdateSubscriptionCommandHandler(subRepo, familyRepo, authz)
		cmd := command.UpdateSubscriptionCommand{SubscriptionID: cmdSub.Id(), Owner: cmdSub.Owner(), FamilyUsers: cmdSub.FamilyUsers().Values()}
		res := h.Handle(t.Context(), cmd)

		assert.True(t, res.IsFaulted())
	})

	t.Run("returns fault when validation fails", func(t *testing.T) {
		subRepo := ports.NewMockSubscriptionRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)

		existing := newPersonalSubscription()
		name := "Invalid"
		subRepo.EXPECT().GetById(t.Context(), existing.Id()).Return(existing, nil)
		authz.EXPECT().Can(t.Context(), authorization.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)

		h := command.NewUpdateSubscriptionCommandHandler(subRepo, familyRepo, authz)
		cmd := command.UpdateSubscriptionCommand{SubscriptionID: existing.Id(), FriendlyName: &name, ProviderID: existing.ProviderId(), Owner: existing.Owner(), StartDate: time.Now().Add(-time.Hour), Recurrency: subscription.CustomRecurrency}
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
		subRepo.EXPECT().GetById(t.Context(), existing.Id()).Return(existing, nil)
		authz.EXPECT().Can(t.Context(), authorization.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)
		subRepo.EXPECT().Save(t.Context(), mock.MatchedBy(func(entities []subscription.Subscription) bool {
			return len(entities) == 1 && entities[0] == existing
		})).Return(errors.New("save failed"))

		h := command.NewUpdateSubscriptionCommandHandler(subRepo, familyRepo, authz)
		cmd := command.UpdateSubscriptionCommand{SubscriptionID: existing.Id(), FriendlyName: &name, ProviderID: existing.ProviderId(), Owner: existing.Owner(), StartDate: time.Now().Add(-2 * time.Hour), Recurrency: subscription.MonthlyRecurrency}
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
		subRepo.EXPECT().GetById(t.Context(), existing.Id()).Return(existing, nil)
		authz.EXPECT().Can(t.Context(), authorization.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)
		subRepo.EXPECT().Save(t.Context(), mock.MatchedBy(func(entities []subscription.Subscription) bool {
			return len(entities) == 1 && entities[0] == existing
		})).Return(nil)

		h := command.NewUpdateSubscriptionCommandHandler(subRepo, familyRepo, authz)
		cmd := command.UpdateSubscriptionCommand{SubscriptionID: existing.Id(), FriendlyName: &name, ProviderID: existing.ProviderId(), Owner: existing.Owner(), StartDate: time.Now().Add(-2 * time.Hour), Recurrency: subscription.MonthlyRecurrency}
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
		subRepo.EXPECT().GetById(t.Context(), existing.Id()).Return(existing, nil)
		authz.EXPECT().Can(t.Context(), authorization.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)
		subRepo.EXPECT().Save(t.Context(), mock.Anything).Return(nil)

		h := command.NewUpdateSubscriptionCommandHandler(subRepo, familyRepo, authz)
		cmd := command.UpdateSubscriptionCommand{SubscriptionID: existing.Id(), FriendlyName: &name, ProviderID: existing.ProviderId(), Owner: existing.Owner(), StartDate: start, EndDate: end, Recurrency: subscription.CustomRecurrency, CustomRecurrency: &customRec}
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
