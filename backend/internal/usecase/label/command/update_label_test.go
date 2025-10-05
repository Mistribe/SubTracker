package command_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	auth "github.com/mistribe/subtracker/internal/domain/authorization"

	"github.com/mistribe/subtracker/internal/domain/label"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/label/command"
	"github.com/mistribe/subtracker/pkg/langext/option"
)

func TestUpdateLabelCommandHandler_Handle(t *testing.T) {
	id := types.NewLabelID()

	newPersonalLabel := func() label.Label {
		return label.NewLabel(
			types.NewLabelID(),
			types.NewPersonalOwner("userID-Test"),
			"label test",
			nil,
			"#FFFFFF",
			time.Now().Add(-time.Hour),
			time.Now().Add(-time.Hour),
		)
	}

	t.Run("returns fault when repository GetById returns error", func(t *testing.T) {
		ctx := t.Context()
		labelRepo := ports.NewMockLabelRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)

		labelRepo.EXPECT().GetById(mock.Anything, id).Return(nil, errors.New("db error"))

		h := command.NewUpdateLabelCommandHandler(labelRepo, familyRepo, authz)
		cmd := command.UpdateLabelCommand{LabelID: id, Name: "n", Color: "#000000", UpdatedAt: option.None[time.Time]()}
		res := h.Handle(ctx, cmd)

		assert.True(t, res.IsFaulted())
	})

	t.Run("returns fault when label not found", func(t *testing.T) {
		ctx := t.Context()
		labelRepo := ports.NewMockLabelRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)

		labelRepo.EXPECT().GetById(mock.Anything, id).Return(nil, nil)

		h := command.NewUpdateLabelCommandHandler(labelRepo, familyRepo, authz)
		cmd := command.UpdateLabelCommand{LabelID: id, Name: "n", Color: "#000000", UpdatedAt: option.None[time.Time]()}
		res := h.Handle(ctx, cmd)

		assert.True(t, res.IsFaulted())
	})

	t.Run("returns fault when permission denied", func(t *testing.T) {
		ctx := t.Context()
		labelRepo := ports.NewMockLabelRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)

		existing := newPersonalLabel()
		labelRepo.EXPECT().GetById(mock.Anything, id).Return(existing, nil)
		authz.EXPECT().Can(mock.Anything, auth.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(auth.ErrUnauthorized)

		h := command.NewUpdateLabelCommandHandler(labelRepo, familyRepo, authz)
		cmd := command.UpdateLabelCommand{LabelID: id, Name: "n", Color: "#000000", UpdatedAt: option.None[time.Time]()}
		res := h.Handle(ctx, cmd)

		assert.True(t, res.IsFaulted())
	})

	t.Run("returns fault when validation fails on updated data", func(t *testing.T) {
		ctx := t.Context()
		labelRepo := ports.NewMockLabelRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)

		existing := newPersonalLabel()
		labelRepo.EXPECT().GetById(mock.Anything, id).Return(existing, nil)
		authz.EXPECT().Can(mock.Anything, auth.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)

		h := command.NewUpdateLabelCommandHandler(labelRepo, familyRepo, authz)
		cmd := command.UpdateLabelCommand{
			LabelID: id, Name: "", Color: "invalid-color", UpdatedAt: option.None[time.Time](),
		}
		res := h.Handle(ctx, cmd)

		assert.True(t, res.IsFaulted())
	})

	t.Run("returns fault when repository Save returns error", func(t *testing.T) {
		ctx := t.Context()
		labelRepo := ports.NewMockLabelRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)

		existing := newPersonalLabel()
		labelRepo.EXPECT().GetById(mock.Anything, id).Return(existing, nil)
		authz.EXPECT().Can(mock.Anything, auth.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)
		labelRepo.EXPECT().Save(mock.Anything, mock.Anything).Return(errors.New("save failed"))

		h := command.NewUpdateLabelCommandHandler(labelRepo, familyRepo, authz)
		cmd := command.UpdateLabelCommand{
			LabelID: id, Name: "new name", Color: "#123456", UpdatedAt: option.None[time.Time](),
		}
		res := h.Handle(ctx, cmd)

		assert.True(t, res.IsFaulted())
	})

	t.Run("returns success and updates fields when UpdatedAt provided", func(t *testing.T) {
		ctx := t.Context()
		labelRepo := ports.NewMockLabelRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)

		existing := newPersonalLabel()
		labelRepo.EXPECT().GetById(mock.Anything, id).Return(existing, nil)
		authz.EXPECT().Can(mock.Anything, auth.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)

		captured := make([]label.Label, 0, 1)
		labelRepo.EXPECT().Save(mock.Anything, mock.Anything).Run(func(_ context.Context, entities ...label.Label) {
			captured = append(captured, entities...)
		}).Return(nil)

		h := command.NewUpdateLabelCommandHandler(labelRepo, familyRepo, authz)

		newTime := time.Now().Add(time.Minute)
		cmd := command.UpdateLabelCommand{
			LabelID: id, Name: "updated name", Color: "#000000", UpdatedAt: option.Some(newTime),
		}
		res := h.Handle(ctx, cmd)

		assert.True(t, res.IsSuccess())
		if assert.Len(t, captured, 1) {
			updated := captured[0]
			assert.Equal(t, "updated name", updated.Name())
			assert.Equal(t, "#000000", updated.Color())
			assert.WithinDuration(t, newTime, updated.UpdatedAt(), time.Millisecond*50)
		}
	})

	t.Run("returns success and sets UpdatedAt to now when not provided", func(t *testing.T) {
		ctx := t.Context()
		labelRepo := ports.NewMockLabelRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)

		existing := newPersonalLabel()
		before := time.Now()
		labelRepo.EXPECT().GetById(mock.Anything, id).Return(existing, nil)
		authz.EXPECT().Can(mock.Anything, auth.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)

		labelRepo.EXPECT().Save(mock.Anything, mock.Anything).Return(nil)

		h := command.NewUpdateLabelCommandHandler(labelRepo, familyRepo, authz)
		cmd := command.UpdateLabelCommand{
			LabelID: id, Name: "updated name 2", Color: "#ABCDEF", UpdatedAt: option.None[time.Time](),
		}
		res := h.Handle(ctx, cmd)

		assert.True(t, res.IsSuccess())
		assert.Equal(t, "updated name 2", existing.Name())
		assert.Equal(t, "#ABCDEF", existing.Color())
		after := time.Now()
		assert.True(t, !existing.UpdatedAt().Before(before) && !existing.UpdatedAt().After(after))
	})
}
