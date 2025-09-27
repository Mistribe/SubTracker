package command_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/mistribe/subtracker/internal/domain/label"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/label/command"
)

func TestCreateLabel(t *testing.T) {
	userId := "userID-Test"

	t.Run("all working", func(t *testing.T) {
		labelRepository := ports.NewMockLabelRepository(t)
		familyRepository := ports.NewMockFamilyRepository(t)
		authorization := ports.NewMockAuthorization(t)
		permissionRequest := ports.NewMockPermissionRequest(t)
		authorization.EXPECT().Can(t.Context(), user.PermissionWrite).Return(permissionRequest)
		permissionRequest.EXPECT().For(mock.Anything).Return(nil)
		labelRepository.EXPECT().GetById(t.Context(), mock.Anything).Return(nil, nil)
		authorization.EXPECT().EnsureWithinLimit(t.Context(), user.FeatureCustomLabels, 1).Return(nil)
		labelRepository.EXPECT().Save(t.Context(), mock.Anything).Return(nil)
		handler := command.NewCreateLabelCommandHandler(labelRepository, familyRepository, authorization)

		owner := types.NewPersonalOwner(userId)
		cmd := command.CreateLabelCommand{
			Label: label.NewLabel(
				uuid.Must(uuid.NewV7()),
				owner,
				"label test",
				nil,
				"#FFFFFF",
				time.Now(),
				time.Now(),
			),
		}

		result := handler.Handle(t.Context(), cmd)

		assert.True(t, result.IsSuccess())
	})

	t.Run("without permission", func(t *testing.T) {
		labelRepository := ports.NewMockLabelRepository(t)
		familyRepository := ports.NewMockFamilyRepository(t)
		authorization := ports.NewMockAuthorization(t)
		permissionRequest := ports.NewMockPermissionRequest(t)
		authorization.EXPECT().Can(t.Context(), user.PermissionWrite).Return(permissionRequest)
		permissionRequest.EXPECT().For(mock.Anything).Return(user.ErrUnauthorized)
		handler := command.NewCreateLabelCommandHandler(labelRepository, familyRepository, authorization)

		owner := types.NewPersonalOwner(userId)
		cmd := command.CreateLabelCommand{
			Label: label.NewLabel(
				uuid.Must(uuid.NewV7()),
				owner,
				"family label test",
				nil,
				"#000000",
				time.Now(),
				time.Now(),
			),
		}

		result := handler.Handle(t.Context(), cmd)

		assert.True(t, result.IsFaulted())
	})

	t.Run("with label already exists but different", func(t *testing.T) {
		labelRepository := ports.NewMockLabelRepository(t)
		familyRepository := ports.NewMockFamilyRepository(t)
		authorization := ports.NewMockAuthorization(t)
		permissionRequest := ports.NewMockPermissionRequest(t)
		authorization.EXPECT().Can(t.Context(), user.PermissionWrite).Return(permissionRequest)
		permissionRequest.EXPECT().For(mock.Anything).Return(nil)
		lbl := label.NewLabel(
			uuid.Must(uuid.NewV7()),
			types.NewPersonalOwner(userId),
			"label test",
			nil,
			"#FFFFFF",
			time.Now(),
			time.Now(),
		)
		labelRepository.EXPECT().GetById(t.Context(), mock.Anything).Return(lbl, nil)
		authorization.EXPECT().EnsureWithinLimit(t.Context(), user.FeatureCustomLabels, 1).Return(nil)
		handler := command.NewCreateLabelCommandHandler(labelRepository, familyRepository, authorization)

		cmd := command.CreateLabelCommand{
			Label: label.NewLabel(
				uuid.Must(uuid.NewV7()),
				types.NewPersonalOwner(userId),
				"label test 2",
				nil,
				"#FFFFFF",
				time.Now(),
				time.Now(),
			),
		}

		result := handler.Handle(t.Context(), cmd)
		assert.True(t, result.IsFaulted())
	})

	t.Run("with label already exists but identical", func(t *testing.T) {
		labelRepository := ports.NewMockLabelRepository(t)
		familyRepository := ports.NewMockFamilyRepository(t)
		authorization := ports.NewMockAuthorization(t)
		permissionRequest := ports.NewMockPermissionRequest(t)
		authorization.EXPECT().Can(t.Context(), user.PermissionWrite).Return(permissionRequest)
		permissionRequest.EXPECT().For(mock.Anything).Return(nil)
		lbl := label.NewLabel(
			uuid.Must(uuid.NewV7()),
			types.NewPersonalOwner(userId),
			"label test",
			nil,
			"#FFFFFF",
			time.Now(),
			time.Now(),
		)
		labelRepository.EXPECT().GetById(t.Context(), mock.Anything).Return(lbl, nil)
		authorization.EXPECT().EnsureWithinLimit(t.Context(), user.FeatureCustomLabels, 1).Return(nil)
		handler := command.NewCreateLabelCommandHandler(labelRepository, familyRepository, authorization)

		cmd := command.CreateLabelCommand{
			Label: lbl,
		}

		result := handler.Handle(t.Context(), cmd)
		assert.True(t, result.IsSuccess())
	})

}
