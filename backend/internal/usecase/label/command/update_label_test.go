package command_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/mistribe/subtracker/internal/domain/auth"
	"github.com/mistribe/subtracker/internal/domain/label"
	"github.com/mistribe/subtracker/internal/domain/user"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/label/command"
	"github.com/mistribe/subtracker/pkg/langext/option"
)

func TestUpdateLabelCommandHandler_Handle(t *testing.T) {
	id := uuid.Must(uuid.NewV7())

	// Common valid existing label used in multiple tests
	newPersonalLabel := func() label.Label {
		return label.NewLabel(
			uuid.Must(uuid.NewV7()),
			auth.NewPersonalOwner("userID-Test"),
			"label test",
			nil,
			"#FFFFFF",
			time.Now().Add(-time.Hour),
			time.Now().Add(-time.Hour),
		)
	}

	// 1. repository GetById returns error
	// Single assertion: result is faulted
	//
	t.Run("returns fault when repository GetById returns error", func(t *testing.T) {
		labelRepo := ports.NewMockLabelRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)

		labelRepo.EXPECT().GetById(t.Context(), id).Return(nil, errors.New("db error"))

		h := command.NewUpdateLabelCommandHandler(labelRepo, familyRepo, authz)
		cmd := command.UpdateLabelCommand{Id: id, Name: "n", Color: "#000000", UpdatedAt: option.None[time.Time]()}
		res := h.Handle(t.Context(), cmd)

		assert.True(t, res.IsFaulted())
	})

	// 2. repository GetById returns nil (not found)
	// Single assertion: result is faulted
	//
	t.Run("returns fault when label not found", func(t *testing.T) {
		labelRepo := ports.NewMockLabelRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)

		labelRepo.EXPECT().GetById(t.Context(), id).Return(nil, nil)

		h := command.NewUpdateLabelCommandHandler(labelRepo, familyRepo, authz)
		cmd := command.UpdateLabelCommand{Id: id, Name: "n", Color: "#000000", UpdatedAt: option.None[time.Time]()}
		res := h.Handle(t.Context(), cmd)

		assert.True(t, res.IsFaulted())
	})

	// 3. permission denied
	// Single assertion: result is faulted
	//
	t.Run("returns fault when permission denied", func(t *testing.T) {
		labelRepo := ports.NewMockLabelRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)

		existing := newPersonalLabel()
		labelRepo.EXPECT().GetById(t.Context(), id).Return(existing, nil)
		authz.EXPECT().Can(t.Context(), user.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(user.ErrUnauthorized)

		h := command.NewUpdateLabelCommandHandler(labelRepo, familyRepo, authz)
		cmd := command.UpdateLabelCommand{Id: id, Name: "n", Color: "#000000", UpdatedAt: option.None[time.Time]()}
		res := h.Handle(t.Context(), cmd)

		assert.True(t, res.IsFaulted())
	})

	// 4. validation fails on updated data
	// Single assertion: result is faulted
	//
	t.Run("returns fault when validation fails on updated data", func(t *testing.T) {
		labelRepo := ports.NewMockLabelRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)

		existing := newPersonalLabel()
		labelRepo.EXPECT().GetById(t.Context(), id).Return(existing, nil)
		authz.EXPECT().Can(t.Context(), user.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)

		h := command.NewUpdateLabelCommandHandler(labelRepo, familyRepo, authz)
		cmd := command.UpdateLabelCommand{Id: id, Name: "", Color: "invalid-color", UpdatedAt: option.None[time.Time]()}
		res := h.Handle(t.Context(), cmd)

		assert.True(t, res.IsFaulted())
	})

	// 5. repository Save returns error
	// Single assertion: result is faulted
	//
	t.Run("returns fault when repository Save returns error", func(t *testing.T) {
		labelRepo := ports.NewMockLabelRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)

		existing := newPersonalLabel()
		labelRepo.EXPECT().GetById(t.Context(), id).Return(existing, nil)
		authz.EXPECT().Can(t.Context(), user.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)
		labelRepo.EXPECT().Save(t.Context(), mock.Anything).Return(errors.New("save failed"))

		h := command.NewUpdateLabelCommandHandler(labelRepo, familyRepo, authz)
		cmd := command.UpdateLabelCommand{Id: id, Name: "new name", Color: "#123456", UpdatedAt: option.None[time.Time]()}
		res := h.Handle(t.Context(), cmd)

		assert.True(t, res.IsFaulted())
	})

	// 6. success with UpdatedAt provided
	// Assertions: success, returned label has updated name/color/updatedAt
	//
	t.Run("returns success and updates fields when UpdatedAt provided", func(t *testing.T) {
		labelRepo := ports.NewMockLabelRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)

		existing := newPersonalLabel()
		labelRepo.EXPECT().GetById(t.Context(), id).Return(existing, nil)
		authz.EXPECT().Can(t.Context(), user.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)

		captured := make([]label.Label, 0, 1)
		labelRepo.EXPECT().Save(t.Context(), mock.Anything).Run(func(_ context.Context, entities ...label.Label) {
			captured = append(captured, entities...)
		}).Return(nil)

		h := command.NewUpdateLabelCommandHandler(labelRepo, familyRepo, authz)

		newTime := time.Now().Add(time.Minute)
		cmd := command.UpdateLabelCommand{Id: id, Name: "updated name", Color: "#000000", UpdatedAt: option.Some(newTime)}
		res := h.Handle(t.Context(), cmd)

		assert.True(t, res.IsSuccess())
		// Assert on the saved entity
		if assert.Len(t, captured, 1) {
			updated := captured[0]
			assert.Equal(t, "updated name", updated.Name())
			assert.Equal(t, "#000000", updated.Color())
			assert.WithinDuration(t, newTime, updated.UpdatedAt(), time.Millisecond*50)
		}
	})

	// 7. success with UpdatedAt not provided -> UpdatedAt set to now
	// Assertions: success, updatedAt close to now
	//
	t.Run("returns success and sets UpdatedAt to now when not provided", func(t *testing.T) {
		labelRepo := ports.NewMockLabelRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)

		existing := newPersonalLabel()
		before := time.Now()
		labelRepo.EXPECT().GetById(t.Context(), id).Return(existing, nil)
		authz.EXPECT().Can(t.Context(), user.PermissionWrite).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)

		labelRepo.EXPECT().Save(t.Context(), mock.Anything).Return(nil)

		h := command.NewUpdateLabelCommandHandler(labelRepo, familyRepo, authz)
		cmd := command.UpdateLabelCommand{Id: id, Name: "updated name 2", Color: "#ABCDEF", UpdatedAt: option.None[time.Time]()}
		res := h.Handle(t.Context(), cmd)

		assert.True(t, res.IsSuccess())
		// The existing entity should have been mutated
		assert.Equal(t, "updated name 2", existing.Name())
		assert.Equal(t, "#ABCDEF", existing.Color())
		after := time.Now()
		assert.True(t, !existing.UpdatedAt().Before(before) && !existing.UpdatedAt().After(after))
	})
}
