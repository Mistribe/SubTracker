package command_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/mistribe/subtracker/internal/domain/user"

	"github.com/mistribe/subtracker/internal/domain/provider"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/provider/command"
)

func TestCreateProviderCommandHandler_Handle(t *testing.T) {
	userId := "userID-Test"
	owner := types.NewPersonalOwner(userId)

	t.Run("success with provided id", func(t *testing.T) {
		provRepo := ports.NewMockProviderRepository(t)
		labelRepo := ports.NewMockLabelRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)

		id := uuid.Must(uuid.NewV7())
		labels := []uuid.UUID{uuid.Must(uuid.NewV7())}
		createdAt := time.Now()

		provRepo.EXPECT().Exists(t.Context(), mock.Anything).Return(false, nil)
		labelRepo.EXPECT().Exists(t.Context(), mock.Anything).Return(true, nil)
		authz.EXPECT().Can(t.Context(), user.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)
		authz.EXPECT().EnsureWithinLimit(t.Context(), user.FeatureCustomProviders, 1).Return(nil)
		provRepo.EXPECT().Save(t.Context(), mock.Anything).Return(nil)

		h := command.NewCreateProviderCommandHandler(provRepo, labelRepo, authz)
		cmd := command.CreateProviderCommand{
			Id:          &id,
			Name:        "Test Provider",
			Labels:      labels,
			Owner:       owner,
			CreatedAt:   &createdAt,
			Description: nil,
		}

		res := h.Handle(t.Context(), cmd)
		assert.True(t, res.IsSuccess())
	})

	t.Run("fail when provided id already exists", func(t *testing.T) {
		provRepo := ports.NewMockProviderRepository(t)
		labelRepo := ports.NewMockLabelRepository(t)
		authz := ports.NewMockAuthorization(t)

		id := uuid.Must(uuid.NewV7())

		provRepo.EXPECT().Exists(t.Context(), mock.Anything).Return(true, nil)

		h := command.NewCreateProviderCommandHandler(provRepo, labelRepo, authz)
		cmd := command.CreateProviderCommand{Id: &id, Name: "X", Owner: owner}

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

		id := uuid.Must(uuid.NewV7())
		existsErr := errors.New("exists error")

		provRepo.EXPECT().Exists(t.Context(), mock.Anything).Return(false, existsErr)

		h := command.NewCreateProviderCommandHandler(provRepo, labelRepo, authz)
		cmd := command.CreateProviderCommand{Id: &id, Name: "X", Owner: owner}

		res := h.Handle(t.Context(), cmd)
		assert.True(t, res.IsFaulted())
		var err error
		res.IfFailure(func(e error) { err = e })
		assert.ErrorIs(t, err, existsErr)
	})

	t.Run("success auto-generates id when nil", func(t *testing.T) {
		provRepo := ports.NewMockProviderRepository(t)
		labelRepo := ports.NewMockLabelRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)

		labels := []uuid.UUID{uuid.Must(uuid.NewV7())}
		var savedId uuid.UUID
		labelRepo.EXPECT().Exists(t.Context(), mock.Anything).Return(true, nil)
		authz.EXPECT().Can(t.Context(), user.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)
		authz.EXPECT().EnsureWithinLimit(t.Context(), user.FeatureCustomProviders, 1).Return(nil)
		provRepo.EXPECT().Save(t.Context(), mock.Anything).Run(func(_ context.Context, entities ...provider.Provider) {
			if len(entities) > 0 {
				savedId = entities[0].Id()
			}
		}).Return(nil)

		h := command.NewCreateProviderCommandHandler(provRepo, labelRepo, authz)
		cmd := command.CreateProviderCommand{
			Id:     nil,
			Name:   "Auto LabelID",
			Labels: labels,
			Owner:  owner,
		}

		res := h.Handle(t.Context(), cmd)
		assert.True(t, res.IsSuccess())
		assert.NotEqual(t, uuid.Nil, savedId)
	})

	t.Run("fail when labels do not exist", func(t *testing.T) {
		provRepo := ports.NewMockProviderRepository(t)
		labelRepo := ports.NewMockLabelRepository(t)
		authz := ports.NewMockAuthorization(t)

		id := uuid.Must(uuid.NewV7())
		labels := []uuid.UUID{uuid.Must(uuid.NewV7())}

		provRepo.EXPECT().Exists(t.Context(), mock.Anything).Return(false, nil)
		labelRepo.EXPECT().Exists(t.Context(), mock.Anything).Return(false, nil)

		h := command.NewCreateProviderCommandHandler(provRepo, labelRepo, authz)
		cmd := command.CreateProviderCommand{Id: &id, Name: "X", Labels: labels, Owner: owner}

		res := h.Handle(t.Context(), cmd)
		assert.True(t, res.IsFaulted())
	})

	t.Run("fail when label repository Exists returns error", func(t *testing.T) {
		provRepo := ports.NewMockProviderRepository(t)
		labelRepo := ports.NewMockLabelRepository(t)
		authz := ports.NewMockAuthorization(t)

		id := uuid.Must(uuid.NewV7())
		labels := []uuid.UUID{uuid.Must(uuid.NewV7())}
		existsErr := errors.New("label exists error")

		provRepo.EXPECT().Exists(t.Context(), mock.Anything).Return(false, nil)
		labelRepo.EXPECT().Exists(t.Context(), mock.Anything).Return(false, existsErr)

		h := command.NewCreateProviderCommandHandler(provRepo, labelRepo, authz)
		cmd := command.CreateProviderCommand{Id: &id, Name: "X", Labels: labels, Owner: owner}

		res := h.Handle(t.Context(), cmd)
		assert.True(t, res.IsFaulted())
	})

	t.Run("fail when unauthorized", func(t *testing.T) {
		provRepo := ports.NewMockProviderRepository(t)
		labelRepo := ports.NewMockLabelRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)

		id := uuid.Must(uuid.NewV7())
		labels := []uuid.UUID{uuid.Must(uuid.NewV7())}

		provRepo.EXPECT().Exists(t.Context(), mock.Anything).Return(false, nil)
		labelRepo.EXPECT().Exists(t.Context(), mock.Anything).Return(true, nil)
		authz.EXPECT().Can(t.Context(), user.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(user.ErrUnauthorized)

		h := command.NewCreateProviderCommandHandler(provRepo, labelRepo, authz)
		cmd := command.CreateProviderCommand{Id: &id, Name: "X", Labels: labels, Owner: owner}

		res := h.Handle(t.Context(), cmd)
		assert.True(t, res.IsFaulted())
	})

	t.Run("fail when ensure within limit returns error", func(t *testing.T) {
		provRepo := ports.NewMockProviderRepository(t)
		labelRepo := ports.NewMockLabelRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)

		id := uuid.Must(uuid.NewV7())
		labels := []uuid.UUID{uuid.Must(uuid.NewV7())}
		limitErr := user.NewLimitExceededErr(user.FeatureCustomProviders, 1, 1)

		provRepo.EXPECT().Exists(t.Context(), mock.Anything).Return(false, nil)
		labelRepo.EXPECT().Exists(t.Context(), mock.Anything).Return(true, nil)
		authz.EXPECT().Can(t.Context(), user.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)
		authz.EXPECT().EnsureWithinLimit(t.Context(), user.FeatureCustomProviders, 1).Return(limitErr)

		h := command.NewCreateProviderCommandHandler(provRepo, labelRepo, authz)
		cmd := command.CreateProviderCommand{Id: &id, Name: "X", Labels: labels, Owner: owner}

		res := h.Handle(t.Context(), cmd)
		assert.True(t, res.IsFaulted())
	})

	t.Run("fail when validation error (empty name)", func(t *testing.T) {
		provRepo := ports.NewMockProviderRepository(t)
		labelRepo := ports.NewMockLabelRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)

		id := uuid.Must(uuid.NewV7())
		labels := []uuid.UUID{uuid.Must(uuid.NewV7())}

		provRepo.EXPECT().Exists(t.Context(), mock.Anything).Return(false, nil)
		labelRepo.EXPECT().Exists(t.Context(), mock.Anything).Return(true, nil)
		authz.EXPECT().Can(t.Context(), user.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)
		authz.EXPECT().EnsureWithinLimit(t.Context(), user.FeatureCustomProviders, 1).Return(nil)

		h := command.NewCreateProviderCommandHandler(provRepo, labelRepo, authz)
		cmd := command.CreateProviderCommand{Id: &id, Name: "   ", Labels: labels, Owner: owner}

		res := h.Handle(t.Context(), cmd)
		assert.True(t, res.IsFaulted())
	})

	t.Run("fail when repository.Save returns error", func(t *testing.T) {
		provRepo := ports.NewMockProviderRepository(t)
		labelRepo := ports.NewMockLabelRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)

		id := uuid.Must(uuid.NewV7())
		labels := []uuid.UUID{uuid.Must(uuid.NewV7())}
		saveErr := errors.New("save failed")

		provRepo.EXPECT().Exists(t.Context(), mock.Anything).Return(false, nil)
		labelRepo.EXPECT().Exists(t.Context(), mock.Anything).Return(true, nil)
		authz.EXPECT().Can(t.Context(), user.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)
		authz.EXPECT().EnsureWithinLimit(t.Context(), user.FeatureCustomProviders, 1).Return(nil)
		provRepo.EXPECT().Save(t.Context(), mock.Anything).Return(saveErr)

		h := command.NewCreateProviderCommandHandler(provRepo, labelRepo, authz)
		cmd := command.CreateProviderCommand{Id: &id, Name: "X", Labels: labels, Owner: owner}

		res := h.Handle(t.Context(), cmd)
		assert.True(t, res.IsFaulted())
	})
}
