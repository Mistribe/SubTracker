package exchange_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/mistribe/subtracker/internal/adapters/exchange"
	dcurrency "github.com/mistribe/subtracker/internal/domain/currency"
	"github.com/mistribe/subtracker/internal/ports"
)

func TestExchange_ToCurrencyAt_InvalidAmount(t *testing.T) {
	localCacheMock := ports.NewMockLocalCache(t)
	currencyRepositoryMock := ports.NewMockCurrencyRepository(t)
	service := exchange.New(localCacheMock, currencyRepositoryMock)

	invalid := dcurrency.NewInvalidAmount()
	out, err := service.ToCurrencyAt(context.Background(), invalid, dcurrency.USD, time.Now())
	require.Error(t, err)
	require.Equal(t, "invalid amount", err.Error())
	require.True(t, out.IsEqual(invalid))
}

func TestExchange_ToCurrencyAt_SameCurrency(t *testing.T) {
	localCacheMock := ports.NewMockLocalCache(t)
	currencyRepositoryMock := ports.NewMockCurrencyRepository(t)
	service := exchange.New(localCacheMock, currencyRepositoryMock)

	initial := dcurrency.NewAmount(42.0, dcurrency.USD)
	out, err := service.ToCurrencyAt(context.Background(), initial, dcurrency.USD, time.Now())
	require.NoError(t, err)
	require.True(t, out.IsEqual(initial))
	require.Nil(t, out.Source())
}

func TestExchange_ToCurrencyAt_ErrorFromRepository(t *testing.T) {
	localCacheMock := ports.NewMockLocalCache(t)
	currencyRepositoryMock := ports.NewMockCurrencyRepository(t)
	service := exchange.New(localCacheMock, currencyRepositoryMock)

	at := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)

	// Cache miss
	localCacheMock.On("Get", mock.Anything).Return(0.0).Once()
	// Repository error
	currencyRepositoryMock.On(
		"GetRateAt",
		mock.Anything, // ctx
		mock.Anything, // from
		mock.Anything, // to
		mock.Anything, // at
	).Return(nil, errors.New("db error")).Once()

	initial := dcurrency.NewAmount(10, dcurrency.USD)
	out, err := service.ToCurrencyAt(context.Background(), initial, dcurrency.EUR, at)

	require.Error(t, err)
	require.Equal(t, "db error", err.Error())
	require.True(t, out.IsEqual(initial))
}

func TestExchange_ToCurrencyAt_SuccessFromCache(t *testing.T) {
	localCacheMock := ports.NewMockLocalCache(t)
	currencyRepositoryMock := ports.NewMockCurrencyRepository(t)
	service := exchange.New(localCacheMock, currencyRepositoryMock)

	at := time.Date(2024, 4, 2, 0, 0, 0, 0, time.UTC)
	rate := 1.10

	// Cache hit
	localCacheMock.On("Get", mock.Anything).Return(rate).Once()
	// Service re-sets the cache with the same rate
	localCacheMock.On("Set", mock.Anything, rate, mock.Anything).Return().Once()

	initial := dcurrency.NewAmount(100, dcurrency.USD)
	out, err := service.ToCurrencyAt(context.Background(), initial, dcurrency.EUR, at)

	require.NoError(t, err)
	require.Equal(t, dcurrency.EUR, out.Currency())
	require.InDelta(t, 110.0, out.Value(), 1e-9)
	require.NotNil(t, out.Source())
	require.True(t, out.Source().IsEqual(initial))

	// Ensure repository was not called on cache hit
	currencyRepositoryMock.AssertNotCalled(t, "GetRateAt", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestExchange_ToCurrencyAt_SuccessFromRepository(t *testing.T) {
	localCacheMock := ports.NewMockLocalCache(t)
	currencyRepositoryMock := ports.NewMockCurrencyRepository(t)
	service := exchange.New(localCacheMock, currencyRepositoryMock)

	at := time.Date(2024, 5, 3, 0, 0, 0, 0, time.UTC)
	rateVal := 1.25
	repoRate := dcurrency.NewRate(
		uuid.New(),
		dcurrency.USD,
		dcurrency.EUR,
		at,
		rateVal,
		time.Now(),
		time.Now(),
	)

	// Cache miss
	localCacheMock.On("Get", mock.Anything).Return(0.0).Once()
	// Repository provides the rate
	currencyRepositoryMock.On(
		"GetRateAt",
		mock.Anything, // ctx
		mock.Anything, // from
		mock.Anything, // to
		mock.Anything, // at
	).Return(repoRate, nil).Once()
	// Cache set after obtaining rate
	localCacheMock.On("Set", mock.Anything, rateVal, mock.Anything).Return().Once()

	initial := dcurrency.NewAmount(80, dcurrency.USD)
	out, err := service.ToCurrencyAt(context.Background(), initial, dcurrency.EUR, at)

	require.NoError(t, err)
	require.Equal(t, dcurrency.EUR, out.Currency())
	require.InDelta(t, 100.0, out.Value(), 1e-9) // 80 * 1.25
	require.NotNil(t, out.Source())
	require.True(t, out.Source().IsEqual(initial))
}
