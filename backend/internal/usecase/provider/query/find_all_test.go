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
	"github.com/mistribe/subtracker/internal/shared"
	"github.com/mistribe/subtracker/internal/usecase/provider/query"
)

func newProvider(name string) provider.Provider {
	id := types.ProviderID(uuid.Must(uuid.NewV7()))
	created := time.Now()
	updated := created
	return provider.NewProvider(id, name, nil, nil, nil, nil, nil, types.NewPersonalOwner(types.UserID("user-1")), created, updated)
}

func TestProviderFindAllQueryHandler_Handle(t *testing.T) {
	userID := types.UserID("user-1")
	newMockConnectedAccount := func(userID types.UserID) *account.MockConnectedAccount {
		acc := account.NewMockConnectedAccount(t)
		acc.EXPECT().UserID().Return(userID).Maybe()
		return acc
	}

	buildQuery := func() query.FindAllQuery { return query.NewFindAllQuery("net", 10, 0) }

	t.Run("returns fault when repository errors", func(t *testing.T) {
		provRepo := ports.NewMockProviderRepository(t)
		auth := ports.NewMockAuthentication(t)
		authz := ports.NewMockAuthorization(t)
		acc := newMockConnectedAccount(userID)
		auth.EXPECT().MustGetConnectedAccount(t.Context()).Return(acc)
		provRepo.EXPECT().GetAllForUser(t.Context(), userID, ports.NewProviderQueryParameters("net", 10, 0)).Return(nil, int64(0), errors.New("db"))

		h := query.NewFindAllQueryHandler(provRepo, auth, authz)
		res := h.Handle(t.Context(), buildQuery())
		assert.True(t, res.IsFaulted())
	})

	t.Run("returns success with pagination", func(t *testing.T) {
		provRepo := ports.NewMockProviderRepository(t)
		auth := ports.NewMockAuthentication(t)
		authz := ports.NewMockAuthorization(t)
		acc := newMockConnectedAccount(userID)
		auth.EXPECT().MustGetConnectedAccount(t.Context()).Return(acc)

		list := []provider.Provider{newProvider("Netflix"), newProvider("Hulu")}
		provRepo.EXPECT().GetAllForUser(t.Context(), userID, ports.NewProviderQueryParameters("net", 10, 0)).Return(list, int64(2), nil)

		h := query.NewFindAllQueryHandler(provRepo, auth, authz)
		res := h.Handle(t.Context(), buildQuery())
		assert.True(t, res.IsSuccess())
		var pg shared.PaginatedResponse[provider.Provider]
		res.IfSuccess(func(v shared.PaginatedResponse[provider.Provider]) { pg = v })
		assert.Equal(t, int64(2), pg.Total())
		assert.Len(t, pg.Data(), 2)
	})
}
