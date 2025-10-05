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
	"github.com/mistribe/subtracker/internal/domain/label"
	"github.com/mistribe/subtracker/internal/domain/provider"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/provider/command"
	"github.com/mistribe/subtracker/pkg/langext/option"
)

func TestUpdateProviderCommandHandler_Handle(t *testing.T) {
	userId := types.UserID("user-id-update-test")
	owner := types.NewPersonalOwner(userId)
	baseProv := func() provider.Provider {
		id := types.ProviderID(uuid.Must(uuid.NewV7()))
		createdAt := time.Now().Add(-24 * time.Hour)
		updatedAt := createdAt
		return provider.NewProvider(
			id,
			"Original Name",
			nil,
			nil,
			nil,
			nil,
			nil,
			[]types.LabelID{},
			owner,
			createdAt,
			updatedAt,
		)
	}

	t.Run("fail when repository.GetById returns error", func(t *testing.T) {
		provRepo := ports.NewMockProviderRepository(t)
		labelRepo := ports.NewMockLabelRepository(t)
		authz := ports.NewMockAuthorization(t)

		id := types.ProviderID(uuid.Must(uuid.NewV7()))
		getErr := errors.New("get error")
		provRepo.EXPECT().GetById(mock.Anything, mock.Anything).Return(nil, getErr)

		h := command.NewUpdateProviderCommandHandler(provRepo, labelRepo, authz)
		cmd := command.UpdateProviderCommand{ProviderID: id, Name: "X"}

		res := h.Handle(t.Context(), cmd)
		assert.True(t, res.IsFaulted())
		var err error
		res.IfFailure(func(e error) { err = e })
		assert.ErrorIs(t, err, getErr)
	})

	t.Run("fail when provider not found", func(t *testing.T) {
		provRepo := ports.NewMockProviderRepository(t)
		labelRepo := ports.NewMockLabelRepository(t)
		authz := ports.NewMockAuthorization(t)

		id := types.ProviderID(uuid.Must(uuid.NewV7()))
		provRepo.EXPECT().GetById(mock.Anything, mock.Anything).Return(nil, nil)

		h := command.NewUpdateProviderCommandHandler(provRepo, labelRepo, authz)
		cmd := command.UpdateProviderCommand{ProviderID: id, Name: "X"}

		res := h.Handle(t.Context(), cmd)
		assert.True(t, res.IsFaulted())
		res.IfFailure(func(e error) {
			assert.ErrorIs(t, e, provider.ErrProviderNotFound)
		})
	})

	t.Run("fail when unauthorized", func(t *testing.T) {
		provRepo := ports.NewMockProviderRepository(t)
		labelRepo := ports.NewMockLabelRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)

		prov := baseProv()
		provRepo.EXPECT().GetById(mock.Anything, mock.Anything).Return(prov, nil)
		authz.EXPECT().Can(mock.Anything, authorization.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(authorization.ErrUnauthorized)

		h := command.NewUpdateProviderCommandHandler(provRepo, labelRepo, authz)
		cmd := command.UpdateProviderCommand{ProviderID: prov.Id(), Name: "X"}

		res := h.Handle(t.Context(), cmd)
		assert.True(t, res.IsFaulted())
		res.IfFailure(func(e error) {
			assert.ErrorIs(t, e, authorization.ErrUnauthorized)
		})
	})

	t.Run("fail when labels do not exist", func(t *testing.T) {
		provRepo := ports.NewMockProviderRepository(t)
		labelRepo := ports.NewMockLabelRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)

		prov := baseProv()
		labels := []types.LabelID{types.LabelID(uuid.Must(uuid.NewV7()))}

		provRepo.EXPECT().GetById(mock.Anything, mock.Anything).Return(prov, nil)
		authz.EXPECT().Can(mock.Anything, authorization.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)
		labelRepo.EXPECT().Exists(mock.Anything, mock.Anything).Return(false, nil)

		h := command.NewUpdateProviderCommandHandler(provRepo, labelRepo, authz)
		cmd := command.UpdateProviderCommand{ProviderID: prov.Id(), Name: "New Name", Labels: labels}

		res := h.Handle(t.Context(), cmd)
		assert.True(t, res.IsFaulted())
		res.IfFailure(func(e error) {
			assert.ErrorIs(t, e, label.ErrLabelNotFound)
		})
	})

	t.Run("fail when label repository Exists returns error", func(t *testing.T) {
		provRepo := ports.NewMockProviderRepository(t)
		labelRepo := ports.NewMockLabelRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)

		prov := baseProv()
		labels := []types.LabelID{types.LabelID(uuid.Must(uuid.NewV7()))}
		existsErr := errors.New("label exists error")

		provRepo.EXPECT().GetById(mock.Anything, mock.Anything).Return(prov, nil)
		authz.EXPECT().Can(mock.Anything, authorization.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)
		labelRepo.EXPECT().Exists(mock.Anything, mock.Anything).Return(false, existsErr)

		h := command.NewUpdateProviderCommandHandler(provRepo, labelRepo, authz)
		cmd := command.UpdateProviderCommand{ProviderID: prov.Id(), Name: "New Name", Labels: labels}

		res := h.Handle(t.Context(), cmd)
		assert.True(t, res.IsFaulted())
		res.IfFailure(func(e error) {
			assert.ErrorIs(t, e, existsErr)
		})
	})

	t.Run("fail when validation error (empty name)", func(t *testing.T) {
		provRepo := ports.NewMockProviderRepository(t)
		labelRepo := ports.NewMockLabelRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)

		prov := baseProv()
		labels := []types.LabelID{types.LabelID(uuid.Must(uuid.NewV7()))}

		provRepo.EXPECT().GetById(mock.Anything, mock.Anything).Return(prov, nil)
		authz.EXPECT().Can(mock.Anything, authorization.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)
		labelRepo.EXPECT().Exists(mock.Anything, mock.Anything).Return(true, nil)

		h := command.NewUpdateProviderCommandHandler(provRepo, labelRepo, authz)
		cmd := command.UpdateProviderCommand{ProviderID: prov.Id(), Name: "   ", Labels: labels}

		res := h.Handle(t.Context(), cmd)
		assert.True(t, res.IsFaulted())
		res.IfFailure(func(e error) {
			assert.Error(t, e)
		})
	})

	t.Run("fail when repository.Save returns error", func(t *testing.T) {
		provRepo := ports.NewMockProviderRepository(t)
		labelRepo := ports.NewMockLabelRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)

		prov := baseProv()
		labels := []types.LabelID{types.LabelID(uuid.Must(uuid.NewV7()))}
		saveErr := errors.New("save failed")

		provRepo.EXPECT().GetById(mock.Anything, mock.Anything).Return(prov, nil)
		authz.EXPECT().Can(mock.Anything, authorization.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)
		labelRepo.EXPECT().Exists(mock.Anything, mock.Anything).Return(true, nil)
		provRepo.EXPECT().Save(mock.Anything, mock.Anything).Return(saveErr)

		h := command.NewUpdateProviderCommandHandler(provRepo, labelRepo, authz)
		cmd := command.UpdateProviderCommand{ProviderID: prov.Id(), Name: "New Name", Labels: labels}

		res := h.Handle(t.Context(), cmd)
		assert.True(t, res.IsFaulted())
		res.IfFailure(func(e error) {
			assert.ErrorIs(t, e, saveErr)
		})
	})

	t.Run("success updates fields and uses provided UpdatedAt", func(t *testing.T) {
		provRepo := ports.NewMockProviderRepository(t)
		labelRepo := ports.NewMockLabelRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)

		prov := baseProv()
		labels := []types.LabelID{types.LabelID(uuid.Must(uuid.NewV7()))}
		desc := "New Desc"
		icon := "https://icon"
		url := "https://site"
		pricing := "https://site/pricing"
		customUpdatedAt := time.Now().Add(10 * time.Minute)

		var saved provider.Provider
		provRepo.EXPECT().GetById(mock.Anything, mock.Anything).Return(prov, nil)
		authz.EXPECT().Can(mock.Anything, authorization.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)
		labelRepo.EXPECT().Exists(mock.Anything, mock.Anything).Return(true, nil)
		provRepo.EXPECT().Save(mock.Anything, mock.Anything).Run(func(_ context.Context, entities ...provider.Provider) {
			if len(entities) > 0 {
				saved = entities[0]
			}
		}).Return(nil)

		h := command.NewUpdateProviderCommandHandler(provRepo, labelRepo, authz)
		cmd := command.UpdateProviderCommand{
			ProviderID:     prov.Id(),
			Name:           "Updated Name",
			Description:    &desc,
			IconUrl:        &icon,
			Url:            &url,
			PricingPageUrl: &pricing,
			Labels:         labels,
			UpdatedAt:      option.Some(customUpdatedAt),
		}

		res := h.Handle(t.Context(), cmd)
		assert.True(t, res.IsSuccess())
		// Validate saved entity fields
		assert.Equal(t, "Updated Name", saved.Name())
		assert.Equal(t, &desc, saved.Description())
		assert.Equal(t, &icon, saved.IconUrl())
		assert.Equal(t, &url, saved.Url())
		assert.Equal(t, &pricing, saved.PricingPageUrl())
		assert.Equal(t, labels, saved.Labels().Values())
		assert.Equal(t, customUpdatedAt, saved.UpdatedAt())
	})

	t.Run("success updates updatedAt when none (uses now)", func(t *testing.T) {
		provRepo := ports.NewMockProviderRepository(t)
		labelRepo := ports.NewMockLabelRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)

		prov := baseProv()
		labels := []types.LabelID{types.LabelID(uuid.Must(uuid.NewV7()))}
		before := time.Now()
		var saved provider.Provider

		provRepo.EXPECT().GetById(mock.Anything, mock.Anything).Return(prov, nil)
		authz.EXPECT().Can(mock.Anything, authorization.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)
		labelRepo.EXPECT().Exists(mock.Anything, mock.Anything).Return(true, nil)
		provRepo.EXPECT().Save(mock.Anything, mock.Anything).Run(func(_ context.Context, entities ...provider.Provider) {
			if len(entities) > 0 {
				saved = entities[0]
			}
		}).Return(nil)

		h := command.NewUpdateProviderCommandHandler(provRepo, labelRepo, authz)
		cmd := command.UpdateProviderCommand{ProviderID: prov.Id(), Name: "N", Labels: labels, UpdatedAt: option.None[time.Time]()}

		res := h.Handle(t.Context(), cmd)
		assert.True(t, res.IsSuccess())
		after := time.Now()
		assert.True(t, saved.UpdatedAt().After(before) || saved.UpdatedAt().Equal(before))
		assert.True(t, saved.UpdatedAt().Before(after) || saved.UpdatedAt().Equal(after))
	})
}
