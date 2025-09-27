package command_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	auth "github.com/mistribe/subtracker/internal/domain/authorization"
	"github.com/mistribe/subtracker/internal/domain/billing"
	"github.com/mistribe/subtracker/internal/domain/label"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/label/command"
	"github.com/mistribe/subtracker/pkg/langext/option"
)

func TestCreateLabel(t *testing.T) {
	userId := types.UserID("userID-Test")

	t.Run("all working", func(t *testing.T) {
		labelRepository := ports.NewMockLabelRepository(t)
		familyRepository := ports.NewMockFamilyRepository(t)
		authorization := ports.NewMockAuthorization(t)
		permissionRequest := ports.NewMockPermissionRequest(t)
		entitlementResolver := ports.NewMockEntitlementResolver(t)

		labelID := types.NewLabelID()

		authorization.EXPECT().Can(t.Context(), auth.PermissionWrite).Return(permissionRequest)
		permissionRequest.EXPECT().For(mock.Anything).Return(nil)
		labelRepository.EXPECT().GetById(t.Context(), labelID).Return(nil, nil)
		entitlementResolver.EXPECT().CheckQuota(t.Context(), billing.FeatureIdCustomLabels, 1).Return(true,
			billing.EffectiveEntitlement{}, nil)
		labelRepository.EXPECT().Save(t.Context(), mock.Anything).Return(nil)

		handler := command.NewCreateLabelCommandHandler(labelRepository, familyRepository, authorization, entitlementResolver)

		owner := types.NewPersonalOwner(userId)
		cmd := command.CreateLabelCommand{
			LabelID:   option.Some(labelID),
			Name:      "label test",
			Color:     "#FFFFFF",
			Owner:     owner,
			CreatedAt: option.Some(time.Now()),
		}

		result := handler.Handle(t.Context(), cmd)
		assert.True(t, result.IsSuccess())
	})

	t.Run("without permission", func(t *testing.T) {
		labelRepository := ports.NewMockLabelRepository(t)
		familyRepository := ports.NewMockFamilyRepository(t)
		authorization := ports.NewMockAuthorization(t)
		permissionRequest := ports.NewMockPermissionRequest(t)
		entitlementResolver := ports.NewMockEntitlementResolver(t)

		authorization.EXPECT().Can(t.Context(), auth.PermissionWrite).Return(permissionRequest)
		permissionRequest.EXPECT().For(mock.Anything).Return(auth.ErrUnauthorized)

		handler := command.NewCreateLabelCommandHandler(labelRepository, familyRepository, authorization, entitlementResolver)

		owner := types.NewPersonalOwner(userId)
		cmd := command.CreateLabelCommand{
			Name:      "family label test",
			Color:     "#000000",
			Owner:     owner,
			CreatedAt: option.Some(time.Now()),
		}

		result := handler.Handle(t.Context(), cmd)
		assert.True(t, result.IsFaulted())
	})

	t.Run("with label already exists but different", func(t *testing.T) {
		labelRepository := ports.NewMockLabelRepository(t)
		familyRepository := ports.NewMockFamilyRepository(t)
		authorization := ports.NewMockAuthorization(t)
		entitlementResolver := ports.NewMockEntitlementResolver(t)

		owner := types.NewPersonalOwner(userId)
		labelID := types.NewLabelID()
		existing := label.NewLabel(
			labelID,
			owner,
			"label test",
			nil,
			"#FFFFFF",
			time.Now(),
			time.Now(),
		)

		labelRepository.EXPECT().GetById(t.Context(), labelID).Return(existing, nil)

		handler := command.NewCreateLabelCommandHandler(labelRepository, familyRepository, authorization, entitlementResolver)

		cmd := command.CreateLabelCommand{
			LabelID:   option.Some(labelID),
			Name:      "label test 2",
			Color:     "#FFFFFF",
			Owner:     owner,
			CreatedAt: option.Some(time.Now()),
		}

		result := handler.Handle(t.Context(), cmd)
		assert.True(t, result.IsFaulted())
	})

	t.Run("with label already exists but identical", func(t *testing.T) {
		labelRepository := ports.NewMockLabelRepository(t)
		familyRepository := ports.NewMockFamilyRepository(t)
		authorization := ports.NewMockAuthorization(t)
		entitlementResolver := ports.NewMockEntitlementResolver(t)

		owner := types.NewPersonalOwner(userId)
		labelID := types.NewLabelID()
		existing := label.NewLabel(
			labelID,
			owner,
			"label test",
			nil,
			"#FFFFFF",
			time.Now(),
			time.Now(),
		)

		labelRepository.EXPECT().GetById(t.Context(), labelID).Return(existing, nil)

		handler := command.NewCreateLabelCommandHandler(labelRepository, familyRepository, authorization, entitlementResolver)

		cmd := command.CreateLabelCommand{
			LabelID:   option.Some(labelID),
			Name:      existing.Name(),
			Color:     existing.Color(),
			Owner:     owner,
			CreatedAt: option.Some(time.Now()),
		}

		result := handler.Handle(t.Context(), cmd)
		assert.True(t, result.IsFaulted())
	})
}
