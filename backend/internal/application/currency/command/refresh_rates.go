package command

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/currency"
	"github.com/oleexo/subtracker/internal/infrastructure/exch"
	"github.com/oleexo/subtracker/pkg/langext/result"
	"github.com/oleexo/subtracker/pkg/slicesx"
)

// RefreshCurrencyRatesCommand represents a command to update currency rates if they are older than one day
type RefreshCurrencyRatesCommand struct {
	// No parameters needed as this is a system operation
}

type RefreshCurrencyRatesResponse struct {
	Timestamp time.Time
	Rates     []currency.Rate
}

// RefreshCurrencyRatesCommandHandler handles the update currency rates command
type RefreshCurrencyRatesCommandHandler struct {
	repository currency.Repository
	exchClient exch.Client
	logger     *slog.Logger
}

// NewRefreshCurrencyRatesCommandHandler creates a new RefreshCurrencyRatesCommandHandler
func NewRefreshCurrencyRatesCommandHandler(
	repository currency.Repository,
	exchClient exch.Client,
	logger *slog.Logger,
) *RefreshCurrencyRatesCommandHandler {
	return &RefreshCurrencyRatesCommandHandler{
		repository: repository,
		exchClient: exchClient,
		logger:     logger,
	}
}

// Handle processes the RefreshCurrencyRatesCommand
func (h *RefreshCurrencyRatesCommandHandler) Handle(ctx context.Context,
	_ RefreshCurrencyRatesCommand) result.Result[RefreshCurrencyRatesResponse] {
	// Get the supported currencies
	supportedCurrencies := currency.GetSupportedCurrencies()

	lastUpdatedDate, err := h.repository.GetLatestUpdateDate(ctx)
	if err != nil {
		h.logger.Error("Failed to get latest currency rates update date", "error", err)
		return result.Fail[RefreshCurrencyRatesResponse](err)
	}

	if lastUpdatedDate.Add(currency.TimeBetweenSync).After(time.Now()) {
		h.logger.Info("Currency rates are up to date")
		rates, err := h.repository.GetRatesByDate(ctx, lastUpdatedDate)
		if err != nil {
			h.logger.Error("Failed to get currency rates", "error", err)
			return result.Fail[RefreshCurrencyRatesResponse](err)
		}
		return result.Success(RefreshCurrencyRatesResponse{
			Timestamp: lastUpdatedDate,
			Rates:     rates,
		})
	}

	latestRates, err := h.exchClient.GetLiveExchangeRates(supportedCurrencies)
	if err != nil {
		h.logger.Error("Failed to get live exchange rates", "error", err)
		return result.Fail[RefreshCurrencyRatesResponse](err)
	}

	rates := slicesx.MapToArr(latestRates.Rates, func(c currency.Unit, v float64) currency.Rate {
		return currency.NewRate(
			uuid.Must(uuid.NewV7()),
			currency.USD,
			c,
			latestRates.Timestamp,
			v,
			time.Now(),
			time.Now(),
		)
	})

	if err := h.repository.Save(ctx, rates...); err != nil {
		h.logger.Error("Failed to save currency rates", "error", err)
		return result.Fail[RefreshCurrencyRatesResponse](err)
	}

	h.logger.Info("Successfully updated currency rates", "count", len(rates))
	return result.Success(RefreshCurrencyRatesResponse{
		Timestamp: latestRates.Timestamp,
		Rates:     rates,
	})
}
