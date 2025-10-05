package query_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/mistribe/subtracker/internal/domain/account"
	"github.com/mistribe/subtracker/internal/domain/currency"
	"github.com/mistribe/subtracker/internal/domain/subscription"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/shared"
	"github.com/mistribe/subtracker/internal/usecase/subscription/query"
)

// helper to build a simple subscription with price
func newSubscriptionWithPrice(amount float64, cur currency.Unit) subscription.Subscription {
	name := "TestSub"
	price := subscription.NewPrice(currency.NewAmount(amount, cur))
	return subscription.NewSubscription(
		types.NewSubscriptionID(),
		&name,
		nil,
		types.ProviderID(uuid.Must(uuid.NewV7())),
		price,
		types.NewPersonalOwner(types.UserID("user-1")),
		nil,
		nil,
		nil,
		time.Now().Add(-24*time.Hour),
		nil,
		subscription.MonthlyRecurrency,
		nil,
		time.Now(),
		time.Now(),
	)
}

func TestFindAllQueryHandler_Handle(t *testing.T) {
	userID := types.UserID("user-1")

	newMockConnectedAccount := func(userID types.UserID) *account.MockConnectedAccount {
		acc := account.NewMockConnectedAccount(t)
		acc.EXPECT().UserID().Return(userID).Maybe()
		return acc
	}

	buildParams := func() query.FindAllQuery { return query.NewFindAllQuery("", nil, nil, nil, nil, nil, false, 10, 0) }

	t.Run("returns fault when repository returns error", func(t *testing.T) {
		subRepo := ports.NewMockSubscriptionRepository(t)
		acctSvc := ports.NewMockAccountService(t)
		exchange := ports.NewMockExchange(t)
		auth := ports.NewMockAuthentication(t)
		authz := ports.NewMockAuthorization(t)

		acc := newMockConnectedAccount(userID)
		auth.EXPECT().MustGetConnectedAccount(t.Context()).Return(acc)
		params := ports.NewSubscriptionQueryParameters("", nil, nil, nil, nil, nil, false, 10, 0)
		subRepo.EXPECT().GetAllForUser(t.Context(), userID, params).Return(nil, int64(0), errors.New("db"))

		h := query.NewFindAllQueryHandler(subRepo, acctSvc, exchange, auth, authz)
		res := h.Handle(t.Context(), buildParams())
		assert.True(t, res.IsFaulted())
	})

	t.Run("returns fault when exchange fails for a subscription", func(t *testing.T) {
		subRepo := ports.NewMockSubscriptionRepository(t)
		acctSvc := ports.NewMockAccountService(t)
		exchange := ports.NewMockExchange(t)
		auth := ports.NewMockAuthentication(t)
		authz := ports.NewMockAuthorization(t)

		acc := newMockConnectedAccount(userID)
		auth.EXPECT().MustGetConnectedAccount(t.Context()).Return(acc)

		s1 := newSubscriptionWithPrice(10, currency.USD)
		params := ports.NewSubscriptionQueryParameters("", nil, nil, nil, nil, nil, false, 10, 0)
		subRepo.EXPECT().GetAllForUser(t.Context(), userID, params).Return([]subscription.Subscription{s1}, int64(1), nil)
		acctSvc.EXPECT().GetPreferredCurrency(t.Context(), userID).Return(currency.USD)
		exchange.EXPECT().ToCurrencyAt(t.Context(), mock.Anything, currency.USD, mock.Anything).Return(currency.NewInvalidAmount(), errors.New("conv"))

		h := query.NewFindAllQueryHandler(subRepo, acctSvc, exchange, auth, authz)
		res := h.Handle(t.Context(), buildParams())
		assert.True(t, res.IsFaulted())
	})

	t.Run("returns success and converts prices", func(t *testing.T) {
		subRepo := ports.NewMockSubscriptionRepository(t)
		acctSvc := ports.NewMockAccountService(t)
		exchange := ports.NewMockExchange(t)
		auth := ports.NewMockAuthentication(t)
		authz := ports.NewMockAuthorization(t)

		acc := newMockConnectedAccount(userID)
		auth.EXPECT().MustGetConnectedAccount(t.Context()).Return(acc)

		s1 := newSubscriptionWithPrice(10, currency.USD)
		s2 := newSubscriptionWithPrice(20, currency.USD)
		params := ports.NewSubscriptionQueryParameters("", nil, nil, nil, nil, nil, false, 10, 0)
		subRepo.EXPECT().GetAllForUser(t.Context(), userID, params).Return([]subscription.Subscription{s1, s2}, int64(2), nil)
		acctSvc.EXPECT().GetPreferredCurrency(t.Context(), userID).Return(currency.EUR)
		exchange.EXPECT().ToCurrencyAt(t.Context(), mock.Anything, currency.EUR, mock.Anything).Return(currency.NewAmount(11, currency.EUR), nil).Twice()

		h := query.NewFindAllQueryHandler(subRepo, acctSvc, exchange, auth, authz)
		res := h.Handle(t.Context(), buildParams())
		assert.True(t, res.IsSuccess())
		var pg shared.PaginatedResponse[subscription.Subscription]
		res.IfSuccess(func(v shared.PaginatedResponse[subscription.Subscription]) { pg = v })
		assert.Equal(t, int64(2), pg.Total())
		assert.Len(t, pg.Data(), 2)
		assert.Equal(t, 11.0, pg.Data()[0].Price().Amount().Value())
		assert.Equal(t, currency.EUR, pg.Data()[0].Price().Amount().Currency())
	})
}
