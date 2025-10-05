package command_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/mistribe/subtracker/internal/domain/authorization"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/subscription/command"
)

func TestDeleteSubscriptionCommandHandler_Handle(t *testing.T) {
	id := uuid.Must(uuid.NewV7())
	subID := types.SubscriptionID(id)

	t.Run("returns fault when repository GetById returns error", func(t *testing.T) {
		subRepo := ports.NewMockSubscriptionRepository(t)
		authz := ports.NewMockAuthorization(t)
		subRepo.EXPECT().GetById(t.Context(), subID).Return(nil, errors.New("db error"))

		h := command.NewDeleteSubscriptionCommandHandler(subRepo, authz)
		cmd := command.DeleteSubscriptionCommand{SubscriptionID: subID}
		res := h.Handle(t.Context(), cmd)

		assert.True(t, res.IsFaulted())
	})

	t.Run("returns fault when subscription not found", func(t *testing.T) {
		subRepo := ports.NewMockSubscriptionRepository(t)
		authz := ports.NewMockAuthorization(t)
		subRepo.EXPECT().GetById(t.Context(), subID).Return(nil, nil)

		h := command.NewDeleteSubscriptionCommandHandler(subRepo, authz)
		cmd := command.DeleteSubscriptionCommand{SubscriptionID: subID}
		res := h.Handle(t.Context(), cmd)

		assert.True(t, res.IsFaulted())
	})

	t.Run("returns fault when permission denied", func(t *testing.T) {
		subRepo := ports.NewMockSubscriptionRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)

		existing := newPersonalSubscription()
		subRepo.EXPECT().GetById(t.Context(), subID).Return(existing, nil)
		authz.EXPECT().Can(t.Context(), authorization.PermissionDelete).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(authorization.ErrUnauthorized)

		h := command.NewDeleteSubscriptionCommandHandler(subRepo, authz)
		cmd := command.DeleteSubscriptionCommand{SubscriptionID: subID}
		res := h.Handle(t.Context(), cmd)

		assert.True(t, res.IsFaulted())
	})

	t.Run("returns fault when repository Delete returns error", func(t *testing.T) {
		subRepo := ports.NewMockSubscriptionRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)

		existing := newPersonalSubscription()
		subRepo.EXPECT().GetById(t.Context(), subID).Return(existing, nil)
		authz.EXPECT().Can(t.Context(), authorization.PermissionDelete).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)
		subRepo.EXPECT().Delete(t.Context(), subID).Return(false, errors.New("delete failed"))

		h := command.NewDeleteSubscriptionCommandHandler(subRepo, authz)
		cmd := command.DeleteSubscriptionCommand{SubscriptionID: subID}
		res := h.Handle(t.Context(), cmd)

		assert.True(t, res.IsFaulted())
	})

	t.Run("returns success when repository Delete succeeds (true)", func(t *testing.T) {
		subRepo := ports.NewMockSubscriptionRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)

		existing := newPersonalSubscription()
		subRepo.EXPECT().GetById(t.Context(), subID).Return(existing, nil)
		authz.EXPECT().Can(t.Context(), authorization.PermissionDelete).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)
		subRepo.EXPECT().Delete(t.Context(), subID).Return(true, nil)

		h := command.NewDeleteSubscriptionCommandHandler(subRepo, authz)
		cmd := command.DeleteSubscriptionCommand{SubscriptionID: subID}
		res := h.Handle(t.Context(), cmd)

		assert.True(t, res.IsSuccess())
	})

	t.Run("returns success when repository Delete succeeds (false)", func(t *testing.T) {
		subRepo := ports.NewMockSubscriptionRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)

		existing := newPersonalSubscription()
		subRepo.EXPECT().GetById(t.Context(), subID).Return(existing, nil)
		authz.EXPECT().Can(t.Context(), authorization.PermissionDelete).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)
		subRepo.EXPECT().Delete(t.Context(), subID).Return(false, nil)

		h := command.NewDeleteSubscriptionCommandHandler(subRepo, authz)
		cmd := command.DeleteSubscriptionCommand{SubscriptionID: subID}
		res := h.Handle(t.Context(), cmd)

		assert.True(t, res.IsSuccess())
	})
}
