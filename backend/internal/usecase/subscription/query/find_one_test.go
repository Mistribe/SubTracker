package query_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/mistribe/subtracker/internal/domain/account"
	"github.com/mistribe/subtracker/internal/domain/currency"
	"github.com/mistribe/subtracker/internal/domain/subscription"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/subscription/query"
)

func TestFindOneQueryHandler_Handle(t *testing.T) {
	userID := types.UserID("user-1")
	newMockConnectedAccount := func(userID types.UserID) *account.MockConnectedAccount {
		acc := account.NewMockConnectedAccount(t)
		acc.EXPECT().UserID().Return(userID).Maybe()
		return acc
	}

	buildSub := func() subscription.Subscription {
		name := "FindOne"
		price := subscription.NewPrice(currency.NewAmount(10, currency.USD))
		return subscription.NewSubscription(
			types.NewSubscriptionID(), &name, nil, types.ProviderID(uuid.Must(uuid.NewV7())), price,
			types.NewPersonalOwner(userID), nil, nil, nil,
			time.Now().Add(-24*time.Hour), nil, subscription.MonthlyRecurrency, nil, time.Now(), time.Now())
	}

	t.Run("returns fault when repository errors", func(t *testing.T) {
		subRepo := ports.NewMockSubscriptionRepository(t)
		auth := ports.NewMockAuthentication(t)
		acc := newMockConnectedAccount(userID)
		auth.EXPECT().MustGetConnectedAccount(t.Context()).Return(acc)
		id := types.SubscriptionID(uuid.Must(uuid.NewV7()))
		subRepo.EXPECT().GetByIdForUser(t.Context(), userID, id).Return(nil, errors.New("db"))

		h := query.NewFindOneQueryHandler(subRepo, auth)
		res := h.Handle(t.Context(), query.NewFindOneQuery(id))
		assert.True(t, res.IsFaulted())
	})

	t.Run("returns fault when subscription not found", func(t *testing.T) {
		subRepo := ports.NewMockSubscriptionRepository(t)
		auth := ports.NewMockAuthentication(t)
		acc := newMockConnectedAccount(userID)
		auth.EXPECT().MustGetConnectedAccount(t.Context()).Return(acc)
		id := types.SubscriptionID(uuid.Must(uuid.NewV7()))
		subRepo.EXPECT().GetByIdForUser(t.Context(), userID, id).Return(nil, nil)

		h := query.NewFindOneQueryHandler(subRepo, auth)
		res := h.Handle(t.Context(), query.NewFindOneQuery(id))
		assert.True(t, res.IsFaulted())
	})

	t.Run("returns success when subscription found", func(t *testing.T) {
		subRepo := ports.NewMockSubscriptionRepository(t)
		auth := ports.NewMockAuthentication(t)
		acc := newMockConnectedAccount(userID)
		auth.EXPECT().MustGetConnectedAccount(t.Context()).Return(acc)
		sub := buildSub()
		subRepo.EXPECT().GetByIdForUser(t.Context(), userID, sub.Id()).Return(sub, nil)

		h := query.NewFindOneQueryHandler(subRepo, auth)
		res := h.Handle(t.Context(), query.NewFindOneQuery(sub.Id()))
		assert.True(t, res.IsSuccess())
		var got subscription.Subscription
		res.IfSuccess(func(s subscription.Subscription) { got = s })
		assert.Equal(t, sub.Id(), got.Id())
	})
}
