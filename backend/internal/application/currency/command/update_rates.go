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

// UpdateCurrencyRatesCommand represents a command to update currency rates if they are older than one day
type UpdateCurrencyRatesCommand struct {
	// No parameters needed as this is a system operation
}

// UpdateCurrencyRatesCommandHandler handles the update currency rates command
type UpdateCurrencyRatesCommandHandler struct {
	repository currency.Repository
	exchClient exch.Client
	logger     *slog.Logger
}

// NewUpdateCurrencyRatesCommandHandler creates a new UpdateCurrencyRatesCommandHandler
func NewUpdateCurrencyRatesCommandHandler(
	repository currency.Repository,
	exchClient exch.Client,
	logger *slog.Logger,
) *UpdateCurrencyRatesCommandHandler {
	return &UpdateCurrencyRatesCommandHandler{
		repository: repository,
		exchClient: exchClient,
		logger:     logger,
	}
}

// Handle processes the UpdateCurrencyRatesCommand
func (h *UpdateCurrencyRatesCommandHandler) Handle(ctx context.Context,
	_ UpdateCurrencyRatesCommand) result.Result[bool] {
	// Get the supported currencies
	supportedCurrencies := currency.GetSupportedCurrencies()

	lastUpdatedDate, err := h.repository.GetLatestUpdateDate(ctx)
	if err != nil {
		h.logger.Error("Failed to get latest currency rates update date", "error", err)
		return result.Fail[bool](err)
	}

	if lastUpdatedDate.Add(currency.TimeBetweenSync).After(time.Now()) {
		h.logger.Info("Currency rates are up to date")
		return result.Success(true)
	}

	latestRates, err := h.exchClient.GetLiveExchangeRates(supportedCurrencies)
	if err != nil {
		h.logger.Error("Failed to get live exchange rates", "error", err)
		return result.Fail[bool](err)
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
		return result.Fail[bool](err)
	}

	h.logger.Info("Successfully updated currency rates", "count", len(rates))
	return result.Success(true)
}
