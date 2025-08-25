package currency

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/domain/currency"
	"github.com/mistribe/subtracker/internal/infrastructure/cache"
	"github.com/mistribe/subtracker/internal/infrastructure/exch"
	"github.com/mistribe/subtracker/pkg/slicesx"
)

const (
	RatesCacheName     = "currency_rates"
	RatesCacheDuration = 24 * time.Hour
)

type service struct {
	currencyRepository currency.Repository
	localCache         cache.Local
	exchangeClient     exch.Client
}

func (s service) getRatesFromApi(at time.Time) (currency.Rates, error) {
	if time.Now().Sub(at) < time.Hour {
		liveRates, err := s.exchangeClient.GetLiveExchangeRates(currency.GetSupportedCurrencies())
		if err != nil {
			return nil, err
		}
		rates := slicesx.MapToArr(liveRates.Rates, func(c currency.Unit, v float64) currency.Rate {
			return currency.NewRate(
				uuid.Must(uuid.NewV7()),
				currency.USD,
				c,
				liveRates.Timestamp,
				v,
				time.Now(),
				time.Now(),
			)
		})

		return rates, nil
	} else {
		panic("not implemented")
	}
}

func (s service) getRatesFromCache(at time.Time) (currency.Rates, string, error) {
	// Normalize to date (UTC) for daily rates
	dateKey := at.In(time.UTC).Format("2006-01-02")
	cacheKey := RatesCacheName + ":" + dateKey

	if cacheRates, ok := s.localCache.Get(cacheKey).(currency.Rates); ok {
		return cacheRates, cacheKey, nil
	}

	return nil, cacheKey, nil
}

func (s service) getRatesFromDatabase(ctx context.Context, at time.Time) (currency.Rates, error) {
	// Query repository using normalized date (UTC midnight)
	normalizedAt := time.Date(at.Year(), at.Month(), at.Day(), 0, 0, 0, 0, time.UTC)
	rates, err := s.currencyRepository.GetRatesByDate(ctx, normalizedAt)
	if err != nil {
		return nil, err
	}

	return rates, nil
}

func (s service) getRates(ctx context.Context, at time.Time) (currency.Rates, error) {
	rates, cacheKey, err := s.getRatesFromCache(at)
	if err != nil {
		return nil, err
	}
	if rates != nil {
		return rates, nil
	}

	rates, err = s.getRatesFromDatabase(ctx, at)
	if err != nil {
		return nil, err
	}

	if rates != nil {
		rates = rates.WithReverse()
		if len(rates) > 0 {
			s.localCache.Set(cacheKey, rates, cache.WithDuration(RatesCacheDuration))
		}
		return rates, nil
	}

	rates, err = s.getRatesFromApi(at)
	if err != nil {
		return nil, err
	}
	err = s.currencyRepository.Save(ctx, rates...)
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
	localCache cache.Local,
	exchangeClient exch.Client) currency.Service {
	return service{
		currencyRepository: currencyRepository,
		localCache:         localCache,
		exchangeClient:     exchangeClient,
	}
}
