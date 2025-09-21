package command_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/mistribe/subtracker/internal/domain/auth"
	"github.com/mistribe/subtracker/internal/domain/provider"
	"github.com/mistribe/subtracker/internal/domain/user"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/provider/command"
)

func makeProvider(id uuid.UUID) provider.Provider {
	owner := auth.NewPersonalOwner("user-1")
	return provider.NewProvider(
		id,
		"Test Provider",
		nil, // key
		nil, // description
		nil, // iconUrl
		nil, // url
		nil, // pricingPageUrl
		nil, // labels
		nil, // plans
		owner,
		time.Now(),
		time.Now(),
	)
}

func TestDeleteProviderCommandHandler_Handle(t *testing.T) {
	ctx := context.Background()

	t.Run("success returns ok=true", func(t *testing.T) {
		repo := ports.NewMockProviderRepository(t)
		authz := ports.NewMockAuthorization(t)

		id := uuid.New()
		prov := makeProvider(id)

		repo.EXPECT().GetById(mock.Anything, id).Return(prov, nil)
		// Authorization allows deletion
		permReq := ports.NewMockPermissionRequest(t)
		authz.EXPECT().Can(mock.Anything, user.PermissionDelete).Return(permReq)
		permReq.EXPECT().For(mock.Anything).Return(nil)
		repo.EXPECT().Delete(mock.Anything, id).Return(true, nil)

		h := command.NewDeleteProviderCommandHandler(repo, authz)
		res := h.Handle(ctx, command.NewDeleteProviderCommand(id))

		require.True(t, res.IsSuccess())
		var ok bool
		res.IfSuccess(func(v bool) { ok = v })
		require.True(t, ok)
	})

	t.Run("success returns ok=false when repository returns false", func(t *testing.T) {
		repo := ports.NewMockProviderRepository(t)
		authz := ports.NewMockAuthorization(t)

		id := uuid.New()
		prov := makeProvider(id)

		repo.EXPECT().GetById(mock.Anything, id).Return(prov, nil)
		// Authorization allows deletion
		permReq := ports.NewMockPermissionRequest(t)
		authz.EXPECT().Can(mock.Anything, user.PermissionDelete).Return(permReq)
		permReq.EXPECT().For(mock.Anything).Return(nil)
		repo.EXPECT().Delete(mock.Anything, id).Return(false, nil)

		h := command.NewDeleteProviderCommandHandler(repo, authz)
		res := h.Handle(ctx, command.NewDeleteProviderCommand(id))

		require.True(t, res.IsSuccess())
		var ok bool
		res.IfSuccess(func(v bool) { ok = v })
		require.False(t, ok)
	})

	t.Run("returns error when repository GetById fails", func(t *testing.T) {
		repo := ports.NewMockProviderRepository(t)
		authz := ports.NewMockAuthorization(t)

		id := uuid.New()
		expectedErr := errors.New("db error")

		repo.EXPECT().GetById(mock.Anything, id).Return(nil, expectedErr)

		h := command.NewDeleteProviderCommandHandler(repo, authz)
		res := h.Handle(ctx, command.NewDeleteProviderCommand(id))

		require.True(t, res.IsFaulted())
		var got error
		res.IfFailure(func(err error) { got = err })
		require.ErrorIs(t, got, expectedErr)
	})

	t.Run("returns provider.ErrProviderNotFound when provider does not exist", func(t *testing.T) {
		repo := ports.NewMockProviderRepository(t)
		authz := ports.NewMockAuthorization(t)

		id := uuid.New()

		repo.EXPECT().GetById(mock.Anything, id).Return(nil, nil)

		h := command.NewDeleteProviderCommandHandler(repo, authz)
		res := h.Handle(ctx, command.NewDeleteProviderCommand(id))

		require.True(t, res.IsFaulted())
		var got error
		res.IfFailure(func(err error) { got = err })
		require.ErrorIs(t, got, provider.ErrProviderNotFound)
	})

	t.Run("returns error when authorization denies deletion", func(t *testing.T) {
		repo := ports.NewMockProviderRepository(t)

		id := uuid.New()
		prov := makeProvider(id)
		expectedErr := errors.New("forbidden")
		authz := ports.NewMockAuthorization(t)

		repo.EXPECT().GetById(mock.Anything, id).Return(prov, nil)
		// Authorization denies deletion
		permReq := ports.NewMockPermissionRequest(t)
		authz.EXPECT().Can(mock.Anything, user.PermissionDelete).Return(permReq)
		permReq.EXPECT().For(mock.Anything).Return(expectedErr)

		h := command.NewDeleteProviderCommandHandler(repo, authz)
		res := h.Handle(ctx, command.NewDeleteProviderCommand(id))

		require.True(t, res.IsFaulted())
		var got error
		res.IfFailure(func(err error) { got = err })
		require.ErrorIs(t, got, expectedErr)
	})

	t.Run("returns error when repository Delete fails", func(t *testing.T) {
		repo := ports.NewMockProviderRepository(t)

		id := uuid.New()
		prov := makeProvider(id)
		expectedErr := errors.New("delete failed")
		authz := ports.NewMockAuthorization(t)

		repo.EXPECT().GetById(mock.Anything, id).Return(prov, nil)
		// Authorization allows deletion
		permReq := ports.NewMockPermissionRequest(t)
		authz.EXPECT().Can(mock.Anything, user.PermissionDelete).Return(permReq)
		permReq.EXPECT().For(mock.Anything).Return(nil)
		repo.EXPECT().Delete(mock.Anything, id).Return(false, expectedErr)

		h := command.NewDeleteProviderCommandHandler(repo, authz)
		res := h.Handle(ctx, command.NewDeleteProviderCommand(id))

		require.True(t, res.IsFaulted())
		var got error
		res.IfFailure(func(err error) { got = err })
		require.ErrorIs(t, got, expectedErr)
	})
}
