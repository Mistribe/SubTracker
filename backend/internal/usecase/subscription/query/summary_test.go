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
	"github.com/mistribe/subtracker/internal/usecase/subscription/query"
)

// build subscription with given params
func buildSub(amount float64, cur currency.Unit, start time.Time, rec subscription.RecurrencyType, months int, labels []subscription.LabelRef, provider types.ProviderID) subscription.Subscription {
	name := "Sub"
	price := subscription.NewPrice(currency.NewAmount(amount, cur))
	var custom *int32
	if rec == subscription.CustomRecurrency {
		c := int32(months)
		custom = &c
	}
	return subscription.NewSubscription(
		types.NewSubscriptionID(), &name, nil, provider, price,
		types.NewPersonalOwner(types.UserID("user-1")), nil, nil, labels,
		start, nil, rec, custom, time.Now(), time.Now(),
	)
}

func TestSummaryQueryHandler_Handle(t *testing.T) {
	userID := types.UserID("user-1")
	newMockConnectedAccount := func(userID types.UserID) *account.MockConnectedAccount {
		acc := account.NewMockConnectedAccount(t)
		acc.EXPECT().UserID().Return(userID).Maybe()
		return acc
	}

	t.Run("returns fault when account repository errors", func(t *testing.T) {
		subRepo := ports.NewMockSubscriptionRepository(t)
		curRepo := ports.NewMockCurrencyRepository(t)
		acctRepo := ports.NewMockAccountRepository(t)
		auth := ports.NewMockAuthentication(t)
		exch := ports.NewMockExchange(t)

		acc := newMockConnectedAccount(userID)
		auth.EXPECT().MustGetConnectedAccount(t.Context()).Return(acc)
		acctRepo.EXPECT().GetById(t.Context(), userID).Return(nil, errors.New("db"))

		h := query.NewSummaryQueryHandler(subRepo, curRepo, acctRepo, auth, exch)
		res := h.Handle(t.Context(), query.SummaryQuery{})
		assert.True(t, res.IsFaulted())
	})

	t.Run("returns fault when currency repository errors", func(t *testing.T) {
		subRepo := ports.NewMockSubscriptionRepository(t)
		curRepo := ports.NewMockCurrencyRepository(t)
		acctRepo := ports.NewMockAccountRepository(t)
		auth := ports.NewMockAuthentication(t)
		exch := ports.NewMockExchange(t)

		acc := newMockConnectedAccount(userID)
		auth.EXPECT().MustGetConnectedAccount(t.Context()).Return(acc)
		acctRepo.EXPECT().GetById(t.Context(), userID).Return(nil, nil)
		curRepo.EXPECT().GetRatesByDate(t.Context(), mock.Anything).Return(currency.Rates{}, errors.New("db"))

		h := query.NewSummaryQueryHandler(subRepo, curRepo, acctRepo, auth, exch)
		res := h.Handle(t.Context(), query.SummaryQuery{})
		assert.True(t, res.IsFaulted())
	})

	t.Run("computes totals and rankings", func(t *testing.T) {
		subRepo := ports.NewMockSubscriptionRepository(t)
		curRepo := ports.NewMockCurrencyRepository(t)
		acctRepo := ports.NewMockAccountRepository(t)
		auth := ports.NewMockAuthentication(t)
		exch := ports.NewMockExchange(t)

		acc := newMockConnectedAccount(userID)
		auth.EXPECT().MustGetConnectedAccount(t.Context()).Return(acc)
		acctRepo.EXPECT().GetById(t.Context(), userID).Return(nil, nil)
		// rates
		curRepo.EXPECT().GetRatesByDate(t.Context(), mock.Anything).Return(currency.Rates{}, nil)

		// Build subscriptions
		start := time.Now().AddDate(0, -2, 0)
		prov1 := types.ProviderID(uuid.Must(uuid.NewV7()))
		prov2 := types.ProviderID(uuid.Must(uuid.NewV7()))
		lab1 := subscription.LabelRef{LabelId: types.LabelID(uuid.Must(uuid.NewV7()))}
		lab2 := subscription.LabelRef{LabelId: types.LabelID(uuid.Must(uuid.NewV7()))}
		s1 := buildSub(30, currency.USD, start, subscription.MonthlyRecurrency, 0, []subscription.LabelRef{lab1}, prov1)
		s2 := buildSub(120, currency.USD, start, subscription.YearlyRecurrency, 0, []subscription.LabelRef{lab1, lab2}, prov2)
		// iterator mock via channel
		ch := make(chan subscription.Subscription, 2)
		ch <- s1
		ch <- s2
		close(ch)
		subRepo.EXPECT().GetAllIt(t.Context(), userID, "").Return(func(yield func(subscription.Subscription) bool) {
			for sub := range ch {
				yield(sub)
			}
		})
		// conversions - we call convert for recurrency amounts; simplify by returning same amount
		exch.EXPECT().ToCurrencyAt(t.Context(), mock.Anything, mock.Anything, mock.Anything).Return(currency.NewAmount(30, currency.USD), nil).Maybe()
		exch.EXPECT().ToCurrencyAt(t.Context(), mock.Anything, mock.Anything, mock.Anything).Return(currency.NewAmount(120, currency.USD), nil).Maybe()

		h := query.NewSummaryQueryHandler(subRepo, curRepo, acctRepo, auth, exch)
		res := h.Handle(t.Context(), query.SummaryQuery{TopProviders: 2, TopLabels: 2, UpcomingRenewals: 2, TotalMonthly: true, TotalYearly: true})
		assert.True(t, res.IsSuccess())
		var out query.SummaryQueryResponse
		res.IfSuccess(func(r query.SummaryQueryResponse) { out = r })
		assert.GreaterOrEqual(t, out.Active, uint16(1))
		assert.NotEmpty(t, out.TopProviders)
		assert.NotEmpty(t, out.TopLabels)
	})
}
