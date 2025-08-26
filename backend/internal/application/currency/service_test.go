package currency_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	appcurrency "github.com/mistribe/subtracker/internal/application/currency"
	"github.com/mistribe/subtracker/internal/domain/currency"
	"github.com/mistribe/subtracker/internal/infrastructure/cache"
	"github.com/mistribe/subtracker/internal/infrastructure/exchange"
)

func TestConvertTo(t *testing.T) {
	ctx := context.Background()
	mockRepo := currency.NewMockRepository(t)
	mockCache := cache.NewMockLocal(t)
	mockExchange := exchange.NewMockClient(t)
	service := appcurrency.NewService(mockRepo, mockCache, mockExchange)

	t.Run("successful conversion", func(t *testing.T) {
		amount := currency.NewAmount(100, currency.USD)
		rates := currency.Rates{
			currency.NewRate(uuid.Must(uuid.NewV7()), currency.USD, currency.EUR, time.Now(), 0.85, time.Now(),
				time.Now()),
		}
		// Cache miss, then cache set after fetching rates from DB
		mockCache.On("Get", mock.Anything).Return(nil)
		mockCache.On("Set", mock.Anything, mock.Anything, mock.Anything).Return()
		mockRepo.On("GetRatesByDate", mock.Anything, mock.Anything).Return(rates, nil)

		result, err := service.ConvertTo(ctx, amount, currency.EUR)

		assert.NoError(t, err)
		assert.Equal(t, currency.EUR, result.Currency())
		assert.Equal(t, 85.0, result.Value())
	})

	t.Run("error getting rates", func(t *testing.T) {
		amount := currency.NewAmount(100, currency.USD)
		// Cache miss triggers repository call which returns an error
		mockCache.On("Get", mock.Anything).Return(nil)
		mockRepo.On("GetRatesByDate", mock.Anything, mock.Anything).Return(nil, errors.New("db error"))

		result, err := service.ConvertTo(ctx, amount, currency.EUR)

		assert.Error(t, err)
		assert.False(t, result.IsValid())
	})

	t.Run("invalid conversion", func(t *testing.T) {
		amount := currency.NewAmount(100, currency.USD)
		rates := currency.Rates{} // Empty rates
		// Cache miss; repository returns empty rates (no Set call since len == 0)
		mockCache.On("Get", mock.Anything).Return(nil)
		mockRepo.On("GetRatesByDate", mock.Anything, mock.Anything).Return(rates, nil)

		result, err := service.ConvertTo(ctx, amount, currency.EUR)

		assert.NoError(t, err)
		assert.False(t, result.IsValid())
	})

}
