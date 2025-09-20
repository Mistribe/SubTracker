package command_test

import (
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
)

func TestDeleteLabelCommandHandler_Handle(t *testing.T) {
	id := uuid.Must(uuid.NewV7())

	t.Run("returns fault when repository GetById returns error", func(t *testing.T) {
		labelRepo := ports.NewMockLabelRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)

		labelRepo.EXPECT().GetById(t.Context(), id).Return(nil, errors.New("db error"))

		h := command.NewDeleteLabelCommandHandler(labelRepo, familyRepo, authz)
		cmd := command.DeleteLabelCommand{Id: id}
		res := h.Handle(t.Context(), cmd)

		assert.True(t, res.IsFaulted())
	})

	t.Run("returns fault when label not found", func(t *testing.T) {
		labelRepo := ports.NewMockLabelRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)

		labelRepo.EXPECT().GetById(t.Context(), id).Return(nil, nil)

		h := command.NewDeleteLabelCommandHandler(labelRepo, familyRepo, authz)
		cmd := command.DeleteLabelCommand{Id: id}
		res := h.Handle(t.Context(), cmd)

		assert.True(t, res.IsFaulted())
	})

	t.Run("returns fault when permission denied", func(t *testing.T) {
		labelRepo := ports.NewMockLabelRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)

		existing := newPersonalLabel()
		labelRepo.EXPECT().GetById(t.Context(), id).Return(existing, nil)
		authz.EXPECT().Can(t.Context(), user.PermissionDelete).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(user.ErrUnauthorized)

		h := command.NewDeleteLabelCommandHandler(labelRepo, familyRepo, authz)
		cmd := command.DeleteLabelCommand{Id: id}
		res := h.Handle(t.Context(), cmd)

		assert.True(t, res.IsFaulted())
	})

	t.Run("returns fault when repository Delete returns error", func(t *testing.T) {
		labelRepo := ports.NewMockLabelRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)

		existing := newPersonalLabel()
		labelRepo.EXPECT().GetById(t.Context(), id).Return(existing, nil)
		authz.EXPECT().Can(t.Context(), user.PermissionDelete).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)
		labelRepo.EXPECT().Delete(t.Context(), id).Return(false, errors.New("delete failed"))

		h := command.NewDeleteLabelCommandHandler(labelRepo, familyRepo, authz)
		cmd := command.DeleteLabelCommand{Id: id}
		res := h.Handle(t.Context(), cmd)

		assert.True(t, res.IsFaulted())
	})

	t.Run("returns fault when repository Delete returns false (not found)", func(t *testing.T) {
		labelRepo := ports.NewMockLabelRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)

		existing := newPersonalLabel()
		labelRepo.EXPECT().GetById(t.Context(), id).Return(existing, nil)
		authz.EXPECT().Can(t.Context(), user.PermissionDelete).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)
		labelRepo.EXPECT().Delete(t.Context(), id).Return(false, nil)

		h := command.NewDeleteLabelCommandHandler(labelRepo, familyRepo, authz)
		cmd := command.DeleteLabelCommand{Id: id}
		res := h.Handle(t.Context(), cmd)

		assert.True(t, res.IsFaulted())
	})

	t.Run("returns success when repository Delete returns true", func(t *testing.T) {
		labelRepo := ports.NewMockLabelRepository(t)
		familyRepo := ports.NewMockFamilyRepository(t)
		authz := ports.NewMockAuthorization(t)
		perm := ports.NewMockPermissionRequest(t)

		existing := newPersonalLabel()
		labelRepo.EXPECT().GetById(t.Context(), id).Return(existing, nil)
		authz.EXPECT().Can(t.Context(), user.PermissionDelete).Return(perm)
		perm.EXPECT().For(mock.Anything).Return(nil)
		labelRepo.EXPECT().Delete(t.Context(), id).Return(true, nil)

		h := command.NewDeleteLabelCommandHandler(labelRepo, familyRepo, authz)
		cmd := command.DeleteLabelCommand{Id: id}
		res := h.Handle(t.Context(), cmd)

		assert.True(t, res.IsSuccess())
	})
}

func newPersonalLabel() label.Label {
	return label.NewLabel(
		uuid.Must(uuid.NewV7()),
		auth.NewPersonalOwner("userID-Test"),
		"label test",
		nil,
		"#FFFFFF",
		time.Now(),
		time.Now(),
	)
}
