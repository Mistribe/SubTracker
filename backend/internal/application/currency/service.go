package currency

import (
	"context"
	"time"

	"github.com/oleexo/subtracker/internal/domain/currency"
	"github.com/oleexo/subtracker/internal/infrastructure/cache"
)

const (
	RatesCacheName     = "currency_rates"
	RatesCacheDuration = 24 * time.Hour
)

type service struct {
	currencyRepository currency.Repository
	localCache         cache.Local
}

func (s service) getRates(ctx context.Context, at time.Time) (currency.Rates, error) {
	// Normalize to date (UTC) for daily rates
	dateKey := at.In(time.UTC).Format("2006-01-02")
	cacheKey := RatesCacheName + ":" + dateKey

	if cacheRates, ok := s.localCache.Get(cacheKey).(currency.Rates); ok {
		return cacheRates, nil
	}

	// Query repository using normalized date (UTC midnight)
	normalizedAt := time.Date(at.Year(), at.Month(), at.Day(), 0, 0, 0, 0, time.UTC)
	rates, err := s.currencyRepository.GetRatesByDate(ctx, normalizedAt)
	if err != nil {
		return nil, err
	}
	rates = rates.WithReverse()

	if len(rates) > 0 {
		s.localCache.Set(cacheKey, rates, cache.WithDuration(RatesCacheDuration))
	}

	return rates, nil
}

func (s service) ConvertTo(ctx context.Context, from currency.Amount, to currency.Unit) (currency.Amount, error) {
	rates, err := s.getRates(ctx, time.Now())
	if err != nil {
		return currency.NewInvalidAmount(), err
	}
	return from.ToCurrency(to, rates), nil
}

func NewService(currencyRepository currency.Repository,
	localCache cache.Local) currency.Service {
	return service{
		currencyRepository: currencyRepository,
		localCache:         localCache,
	}
}
