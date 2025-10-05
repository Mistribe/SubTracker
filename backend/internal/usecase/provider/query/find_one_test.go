package query_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/mistribe/subtracker/internal/domain/account"
	"github.com/mistribe/subtracker/internal/domain/provider"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/provider/query"
)

func buildProvider(name string, owner types.Owner) provider.Provider {
	id := types.ProviderID(uuid.Must(uuid.NewV7()))
	created := time.Now()
	return provider.NewProvider(id, name, nil, nil, nil, nil, nil, nil, owner, created, created)
}

func TestProviderFindOneQueryHandler_Handle(t *testing.T) {
	userID := types.UserID("user-1")
	newMockConnectedAccount := func(userID types.UserID) *account.MockConnectedAccount {
		acc := account.NewMockConnectedAccount(t)
		acc.EXPECT().UserID().Return(userID).Maybe()
		return acc
	}

	t.Run("returns fault when repository errors", func(t *testing.T) {
		provRepo := ports.NewMockProviderRepository(t)
		auth := ports.NewMockAuthentication(t)
		acc := newMockConnectedAccount(userID)
		auth.EXPECT().MustGetConnectedAccount(t.Context()).Return(acc)
		id := types.ProviderID(uuid.Must(uuid.NewV7()))
		provRepo.EXPECT().GetByIdForUser(t.Context(), userID, id).Return(nil, errors.New("db"))

		h := query.NewFindOneQueryHandler(provRepo, auth)
		res := h.Handle(t.Context(), query.NewFindOneQuery(id))
		assert.True(t, res.IsFaulted())
	})

	t.Run("returns fault when provider not found", func(t *testing.T) {
		provRepo := ports.NewMockProviderRepository(t)
		auth := ports.NewMockAuthentication(t)
		acc := newMockConnectedAccount(userID)
		auth.EXPECT().MustGetConnectedAccount(t.Context()).Return(acc)
		id := types.ProviderID(uuid.Must(uuid.NewV7()))
		provRepo.EXPECT().GetByIdForUser(t.Context(), userID, id).Return(nil, nil)

		h := query.NewFindOneQueryHandler(provRepo, auth)
		res := h.Handle(t.Context(), query.NewFindOneQuery(id))
		assert.True(t, res.IsFaulted())
	})

	t.Run("returns success when provider found", func(t *testing.T) {
		provRepo := ports.NewMockProviderRepository(t)
		auth := ports.NewMockAuthentication(t)
		acc := newMockConnectedAccount(userID)
		auth.EXPECT().MustGetConnectedAccount(t.Context()).Return(acc)
		prov := buildProvider("Netflix", types.NewPersonalOwner(userID))
		provRepo.EXPECT().GetByIdForUser(t.Context(), userID, prov.Id()).Return(prov, nil)

		h := query.NewFindOneQueryHandler(provRepo, auth)
		res := h.Handle(t.Context(), query.NewFindOneQuery(prov.Id()))
		assert.True(t, res.IsSuccess())
		var got provider.Provider
		res.IfSuccess(func(p provider.Provider) { got = p })
		assert.Equal(t, prov.Id(), got.Id())
	})
}
