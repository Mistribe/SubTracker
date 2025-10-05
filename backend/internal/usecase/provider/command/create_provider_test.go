package command_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/mistribe/subtracker/internal/domain/authorization"
	"github.com/mistribe/subtracker/internal/domain/billing"
	"github.com/mistribe/subtracker/internal/domain/provider"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/provider/command"
	"github.com/mistribe/subtracker/pkg/langext/option"
)

func TestCreateProviderCommandHandler_Handle(t *testing.T) {
	userId := types.UserID("userID-Test")
	owner := types.NewPersonalOwner(userId)

	t.Run("success with provided id", func(t *testing.T) {
		provRepo := ports.NewMockProviderRepository(t)
		labelRepo := ports.NewMockLabelRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)
		ent := ports.NewMockEntitlementResolver(t)

		id := types.ProviderID(uuid.Must(uuid.NewV7()))
		labels := []types.LabelID{types.LabelID(uuid.Must(uuid.NewV7()))}
		createdAt := time.Now()

		provRepo.EXPECT().Exists(mock.Anything, mock.Anything).Return(false, nil)
		labelRepo.EXPECT().Exists(mock.Anything, mock.Anything).Return(true, nil)
		authz.EXPECT().Can(mock.Anything, authorization.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)
		ent.EXPECT().CheckQuota(mock.Anything, billing.FeatureIdCustomProvidersCount, int64(1)).Return(true, billing.EffectiveEntitlement{}, nil)
		provRepo.EXPECT().Save(mock.Anything, mock.Anything).Return(nil)

		h := command.NewCreateProviderCommandHandler(provRepo, labelRepo, authz, ent)
		cmd := command.CreateProviderCommand{
			ProviderID:  option.Some(id),
			Name:        "Test Provider",
			Labels:      labels,
			Owner:       owner,
			CreatedAt:   option.Some(createdAt),
			Description: nil,
		}

		res := h.Handle(t.Context(), cmd)
		assert.True(t, res.IsSuccess())
	})

	t.Run("fail when provided id already exists", func(t *testing.T) {
		provRepo := ports.NewMockProviderRepository(t)
		labelRepo := ports.NewMockLabelRepository(t)
		authz := ports.NewMockAuthorization(t)
		ent := ports.NewMockEntitlementResolver(t)

		id := types.ProviderID(uuid.Must(uuid.NewV7()))

		provRepo.EXPECT().Exists(mock.Anything, mock.Anything).Return(true, nil)

		h := command.NewCreateProviderCommandHandler(provRepo, labelRepo, authz, ent)
		cmd := command.CreateProviderCommand{ProviderID: option.Some(id), Name: "X", Owner: owner}

		res := h.Handle(t.Context(), cmd)
		assert.True(t, res.IsFaulted())
		var err error
		res.IfFailure(func(e error) { err = e })
		assert.ErrorIs(t, err, provider.ErrProviderAlreadyExists)
	})

	t.Run("fail when repository.Exists returns error", func(t *testing.T) {
		provRepo := ports.NewMockProviderRepository(t)
		labelRepo := ports.NewMockLabelRepository(t)
		authz := ports.NewMockAuthorization(t)
		ent := ports.NewMockEntitlementResolver(t)

		id := types.ProviderID(uuid.Must(uuid.NewV7()))
		existsErr := errors.New("exists error")

		provRepo.EXPECT().Exists(mock.Anything, mock.Anything).Return(false, existsErr)

		h := command.NewCreateProviderCommandHandler(provRepo, labelRepo, authz, ent)
		cmd := command.CreateProviderCommand{ProviderID: option.Some(id), Name: "X", Owner: owner}

		res := h.Handle(t.Context(), cmd)
		assert.True(t, res.IsFaulted())
		var err error
		res.IfFailure(func(e error) { err = e })
		assert.ErrorIs(t, err, existsErr)
	})

	t.Run("success auto-generates id when none", func(t *testing.T) {
		provRepo := ports.NewMockProviderRepository(t)
		labelRepo := ports.NewMockLabelRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)
		ent := ports.NewMockEntitlementResolver(t)

		labels := []types.LabelID{types.LabelID(uuid.Must(uuid.NewV7()))}
		var savedId types.ProviderID
		labelRepo.EXPECT().Exists(mock.Anything, mock.Anything).Return(true, nil)
		authz.EXPECT().Can(mock.Anything, authorization.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)
		ent.EXPECT().CheckQuota(mock.Anything, billing.FeatureIdCustomProvidersCount, int64(1)).Return(true, billing.EffectiveEntitlement{}, nil)
		provRepo.EXPECT().Save(mock.Anything, mock.Anything).Run(func(_ context.Context, entities ...provider.Provider) {
			if len(entities) > 0 {
				savedId = entities[0].Id()
			}
		}).Return(nil)

		h := command.NewCreateProviderCommandHandler(provRepo, labelRepo, authz, ent)
		cmd := command.CreateProviderCommand{
			ProviderID: option.None[types.ProviderID](),
			Name:       "Auto LabelID",
			Labels:     labels,
			Owner:      owner,
		}

		res := h.Handle(t.Context(), cmd)
		assert.True(t, res.IsSuccess())
		assert.NotEqual(t, types.ProviderID(uuid.Nil), savedId)
	})

	t.Run("fail when labels do not exist", func(t *testing.T) {
		provRepo := ports.NewMockProviderRepository(t)
		labelRepo := ports.NewMockLabelRepository(t)
		authz := ports.NewMockAuthorization(t)
		ent := ports.NewMockEntitlementResolver(t)

		id := types.ProviderID(uuid.Must(uuid.NewV7()))
		labels := []types.LabelID{types.LabelID(uuid.Must(uuid.NewV7()))}

		provRepo.EXPECT().Exists(mock.Anything, mock.Anything).Return(false, nil)
		labelRepo.EXPECT().Exists(mock.Anything, mock.Anything).Return(false, nil)

		h := command.NewCreateProviderCommandHandler(provRepo, labelRepo, authz, ent)
		cmd := command.CreateProviderCommand{ProviderID: option.Some(id), Name: "X", Labels: labels, Owner: owner}

		res := h.Handle(t.Context(), cmd)
		assert.True(t, res.IsFaulted())
	})

	t.Run("fail when label repository Exists returns error", func(t *testing.T) {
		provRepo := ports.NewMockProviderRepository(t)
		labelRepo := ports.NewMockLabelRepository(t)
		authz := ports.NewMockAuthorization(t)
		ent := ports.NewMockEntitlementResolver(t)

		id := types.ProviderID(uuid.Must(uuid.NewV7()))
		labels := []types.LabelID{types.LabelID(uuid.Must(uuid.NewV7()))}
		existsErr := errors.New("label exists error")

		provRepo.EXPECT().Exists(mock.Anything, mock.Anything).Return(false, nil)
		labelRepo.EXPECT().Exists(mock.Anything, mock.Anything).Return(false, existsErr)

		h := command.NewCreateProviderCommandHandler(provRepo, labelRepo, authz, ent)
		cmd := command.CreateProviderCommand{ProviderID: option.Some(id), Name: "X", Labels: labels, Owner: owner}

		res := h.Handle(t.Context(), cmd)
		assert.True(t, res.IsFaulted())
	})

	t.Run("fail when unauthorized", func(t *testing.T) {
		provRepo := ports.NewMockProviderRepository(t)
		labelRepo := ports.NewMockLabelRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)
		ent := ports.NewMockEntitlementResolver(t)

		id := types.ProviderID(uuid.Must(uuid.NewV7()))
		labels := []types.LabelID{types.LabelID(uuid.Must(uuid.NewV7()))}

		provRepo.EXPECT().Exists(mock.Anything, mock.Anything).Return(false, nil)
		labelRepo.EXPECT().Exists(mock.Anything, mock.Anything).Return(true, nil)
		authz.EXPECT().Can(mock.Anything, authorization.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(authorization.ErrUnauthorized)

		h := command.NewCreateProviderCommandHandler(provRepo, labelRepo, authz, ent)
		cmd := command.CreateProviderCommand{ProviderID: option.Some(id), Name: "X", Labels: labels, Owner: owner}

		res := h.Handle(t.Context(), cmd)
		assert.True(t, res.IsFaulted())
	})

	t.Run("fail when quota check returns not allowed", func(t *testing.T) {
		provRepo := ports.NewMockProviderRepository(t)
		labelRepo := ports.NewMockLabelRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)
		ent := ports.NewMockEntitlementResolver(t)

		id := types.ProviderID(uuid.Must(uuid.NewV7()))
		labels := []types.LabelID{types.LabelID(uuid.Must(uuid.NewV7()))}

		provRepo.EXPECT().Exists(mock.Anything, mock.Anything).Return(false, nil)
		labelRepo.EXPECT().Exists(mock.Anything, mock.Anything).Return(true, nil)
		authz.EXPECT().Can(mock.Anything, authorization.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)
		ent.EXPECT().CheckQuota(mock.Anything, billing.FeatureIdCustomProvidersCount, int64(1)).Return(false, billing.EffectiveEntitlement{}, nil)

		h := command.NewCreateProviderCommandHandler(provRepo, labelRepo, authz, ent)
		cmd := command.CreateProviderCommand{ProviderID: option.Some(id), Name: "X", Labels: labels, Owner: owner}

		res := h.Handle(t.Context(), cmd)
		assert.True(t, res.IsFaulted())
	})

	t.Run("fail when validation error (empty name)", func(t *testing.T) {
		provRepo := ports.NewMockProviderRepository(t)
		labelRepo := ports.NewMockLabelRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)
		ent := ports.NewMockEntitlementResolver(t)

		id := types.ProviderID(uuid.Must(uuid.NewV7()))
		labels := []types.LabelID{types.LabelID(uuid.Must(uuid.NewV7()))}

		provRepo.EXPECT().Exists(mock.Anything, mock.Anything).Return(false, nil)
		labelRepo.EXPECT().Exists(mock.Anything, mock.Anything).Return(true, nil)
		authz.EXPECT().Can(mock.Anything, authorization.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)
		ent.EXPECT().CheckQuota(mock.Anything, billing.FeatureIdCustomProvidersCount, int64(1)).Return(true, billing.EffectiveEntitlement{}, nil)

		h := command.NewCreateProviderCommandHandler(provRepo, labelRepo, authz, ent)
		cmd := command.CreateProviderCommand{ProviderID: option.Some(id), Name: "   ", Labels: labels, Owner: owner}

		res := h.Handle(t.Context(), cmd)
		assert.True(t, res.IsFaulted())
	})

	t.Run("fail when repository.Save returns error", func(t *testing.T) {
		provRepo := ports.NewMockProviderRepository(t)
		labelRepo := ports.NewMockLabelRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)
		ent := ports.NewMockEntitlementResolver(t)

		id := types.ProviderID(uuid.Must(uuid.NewV7()))
		labels := []types.LabelID{types.LabelID(uuid.Must(uuid.NewV7()))}
		saveErr := errors.New("save failed")

		provRepo.EXPECT().Exists(mock.Anything, mock.Anything).Return(false, nil)
		labelRepo.EXPECT().Exists(mock.Anything, mock.Anything).Return(true, nil)
		authz.EXPECT().Can(mock.Anything, authorization.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)
		ent.EXPECT().CheckQuota(mock.Anything, billing.FeatureIdCustomProvidersCount, int64(1)).Return(true, billing.EffectiveEntitlement{}, nil)
		provRepo.EXPECT().Save(mock.Anything, mock.Anything).Return(saveErr)

		h := command.NewCreateProviderCommandHandler(provRepo, labelRepo, authz, ent)
		cmd := command.CreateProviderCommand{ProviderID: option.Some(id), Name: "X", Labels: labels, Owner: owner}

		res := h.Handle(t.Context(), cmd)
		assert.True(t, res.IsFaulted())
	})
}
