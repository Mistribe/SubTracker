package exchange

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/domain/currency"
	"github.com/mistribe/subtracker/internal/ports"
)

const (
	externalSourceUrl = "https://raw.githubusercontent.com/Mistribe/SubTracker-Data/refs/heads/main/exchange/"
)

type externalSourceModel struct {
	From string             `json:"from"`
	To   map[string]float64 `json:"to"`
}
type exchange struct {
	cache      ports.Cache
	repository ports.CurrencyRepository
	logger     *slog.Logger
}

func getCacheKey(from, to currency.Unit, at time.Time) string {
	return fmt.Sprintf("%s-%s-%s", from, to, at.Format("2006-01-02"))
}

func (e exchange) getRateFromCache(ctx context.Context, key string) float64 {
	rate, ok := e.cache.From(ctx, ports.CacheLevelServer).Get(key).(float64)
	if rate == 0 || !ok {
		return 0
	}

	return rate
}

func (e exchange) getRateFromDatabase(ctx context.Context, from, to currency.Unit, at time.Time) (float64, error) {
	rate, err := e.repository.GetRateAt(ctx, from, to, at)
	if err != nil {
		return 0, err
	}

	if rate == nil {
		return 0, nil
	}
	return rate.ExchangeRate(), nil
}

func (e exchange) getRateFromExternalSource(from, to currency.Unit, at time.Time) (
	float64,
	error) {
	year := at.Year()
	month := at.Month()
	day := at.Day()

	url := fmt.Sprintf("%s/%d/%02d/%02d.json", externalSourceUrl, year, month, day)
	response, err := http.Get(url)
	if err != nil {
		return 0, err
	}

	var source externalSourceModel
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusNotFound {
			return 0, nil
		}
		return 0, fmt.Errorf("external source returned status code %d", response.StatusCode)
	}

	if err = json.NewDecoder(response.Body).Decode(&source); err != nil {
		return 0, err
	}

	// The external source contains only USD to currencies
	usdTo := source.To
	fromCode := from.String()
	toCode := to.String()

	// Identity conversion
	if fromCode == toCode {
		return 1, nil
	}

	// USD -> target
	if fromCode == currency.USD.String() {
		r, ok := usdTo[toCode]
		if !ok || r == 0 {
			return 0, fmt.Errorf("missing USD->%s rate from external source", toCode)
		}
		return r, nil
	}

	// source -> USD
	if toCode == currency.USD.String() {
		rFrom, ok := usdTo[fromCode]
		if !ok || rFrom == 0 {
			return 0, fmt.Errorf("missing USD->%s rate from external source (for inverse to USD)", fromCode)
		}
		return 1 / rFrom, nil
	}

	// source -> target via USD
	rFrom, okFrom := usdTo[fromCode]
	rTo, okTo := usdTo[toCode]
	if !okFrom || rFrom == 0 {
		return 0, fmt.Errorf("missing USD->%s rate from external source (for cross conversion)", fromCode)
	}
	if !okTo || rTo == 0 {
		return 0, fmt.Errorf("missing USD->%s rate from external source (for cross conversion)", toCode)
	}

	// from -> USD -> to : (1 / (USD->from)) * (USD->to)
	return (1 / rFrom) * rTo, nil
}

func (e exchange) getRateAt(ctx context.Context, from, to currency.Unit, at time.Time) (float64, error) {
	key := getCacheKey(from, to, at)
	rate := e.getRateFromCache(ctx, key)
	if rate == 0 {
		var err error
		rate, err = e.getRateFromDatabase(ctx, from, to, at)
		if err != nil {
			return 0, err
		}
	}
	if rate == 0 {
		var err error
		rate, err = e.getRateFromExternalSource(from, to, at)
		if err != nil {
			return 0, err
		}
		if rate == 0 {
			e.logger.Warn("missing rate",
				slog.String("from", from.String()),
				slog.String("to", to.String()),
				slog.Time("at", at))
			rate, err = e.getRateFromExternalSource(from, to, time.Now())
			if err != nil {
				return 0, err
			}
		} else {
			if err = e.saveRateToDatabase(ctx, from, to, rate, at); err != nil {
				return 0, err
			}
		}
	}

	if rate > 0 {
		e.cache.From(ctx, ports.CacheLevelServer).Set(key, rate, ports.WithDuration(time.Hour*24))
		return rate, nil
	}

	return 0, errors.New("missing rate")
}

func (e exchange) ToCurrencyAt(
	ctx context.Context,
	initial currency.Amount,
	target currency.Unit,
	at time.Time) (currency.Amount, error) {
	if !initial.IsValid() {
		return initial, errors.New("invalid amount")
	}
	if initial.Currency() == target {
		return initial, nil
	}
	rate, err := e.getRateAt(ctx, initial.Currency(), target, at)
	if err != nil {
		return initial, err
	}
	newValue := initial.Value() * rate
	return currency.NewAmountWithSource(newValue, target, initial), nil
}

func (e exchange) ToCurrency(
	ctx context.Context,
	initial currency.Amount,
	target currency.Unit) (amount currency.Amount, err error) {
	return e.ToCurrencyAt(ctx, initial, target, time.Now())
}

func (e exchange) saveRateToDatabase(
	ctx context.Context,
	from currency.Unit,
	to currency.Unit,
	rate float64,
	at time.Time) error {

	rateId, err := uuid.NewV7()
	if err != nil {
		return err
	}
	dateOnly := time.Date(at.Year(), at.Month(), at.Day(), 0, 0, 0, 0, at.Location())
	now := time.Now()
	r := currency.NewRate(
		rateId,
		from,
		to,
		dateOnly,
		rate,
		now,
		now,
	)

	return e.repository.Save(ctx, r)
}

func New(
	cache ports.Cache,
	repository ports.CurrencyRepository,
	logger *slog.Logger) ports.Exchange {
	return exchange{
		cache:      cache,
		repository: repository,
		logger:     logger,
	}
}
