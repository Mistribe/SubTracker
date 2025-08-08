package command

import (
	"context"
	"log/slog"
	"time"

	"golang.org/x/text/currency"

	"github.com/oleexo/subtracker/internal/application/core"
	dcurrency "github.com/oleexo/subtracker/internal/domain/currency"
	"github.com/oleexo/subtracker/internal/infrastructure/exch"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

// UpdateCurrencyRatesCommand represents a command to update currency rates if they are older than one day
type UpdateCurrencyRatesCommand struct {
	// No parameters needed as this is a system operation
}

// UpdateCurrencyRatesCommandHandler handles the update currency rates command
type UpdateCurrencyRatesCommandHandler struct {
	repository dcurrency.Repository
	exchClient exch.Client
	logger     *slog.Logger
}

// NewUpdateCurrencyRatesCommandHandler creates a new UpdateCurrencyRatesCommandHandler
func NewUpdateCurrencyRatesCommandHandler(
	repository dcurrency.Repository,
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
func (h *UpdateCurrencyRatesCommandHandler) Handle(ctx context.Context, _ UpdateCurrencyRatesCommand) result.Result[bool] {
	// Get the supported currencies
	supportedCurrencies := dcurrency.GetSupportedCurrencies()

	// Check if we need to update rates
	needsUpdate, err := h.needsRateUpdate(ctx, supportedCurrencies)
	if err != nil {
		h.logger.Error("Failed to check if rates need update", "error", err)
		return result.Fail[bool](err)
	}

	if !needsUpdate {
		h.logger.Info("Currency rates are up to date")
		return result.Success(true)
	}

	// Fetch new rates
	h.logger.Info("Fetching new currency rates")
	liveRates, err := h.exchClient.GetLiveExchangeRates()
	if err != nil {
		h.logger.Error("Failed to fetch live exchange rates", "error", err)
		return result.Fail[bool](err)
	}

	// Convert to domain rates and save
	rates := h.convertToRates(liveRates)
	if err := h.repository.Save(ctx, rates...); err != nil {
		h.logger.Error("Failed to save currency rates", "error", err)
		return result.Fail[bool](err)
	}

	h.logger.Info("Successfully updated currency rates", "count", len(rates))
	return result.Success(true)
}

// needsRateUpdate checks if any of the currency pairs need to be updated
func (h *UpdateCurrencyRatesCommandHandler) needsRateUpdate(ctx context.Context, currencies []currency.Unit) (bool, error) {
	// We only need to check USD to other currency rates as per the requirement
	oneDayAgo := time.Now().AddDate(0, 0, -1)

	for _, cur := range currencies {
		if cur == currency.USD {
			continue // Skip USD to USD rate
		}

		// Check if we have a recent rate from USD to this currency
		rate, err := h.repository.GetLatestRate(ctx, currency.USD, cur)
		if err != nil {
			// If we don't have a rate at all, we need to update
			return true, nil
		}

		// Check if the rate is older than one day
		if rate.RateDate().Before(oneDayAgo) {
			return true, nil
		}
	}

	return false, nil
}

// convertToRates converts live exchange rates to domain rate entities
func (h *UpdateCurrencyRatesCommandHandler) convertToRates(liveRates exch.LiveExchangeRates) []*dcurrency.Rate {
	var rates []*dcurrency.Rate

	// Create a rate for each currency pair
	for targetCur, rate := range liveRates.Rates {
		// Skip if the target currency is USD (USD to USD rate is always 1)
		if targetCur == currency.USD {
			continue
		}

		// Create a new rate from USD to the target currency
		newRate := dcurrency.NewRate(
			currency.USD,
			targetCur,
			liveRates.Timestamp,
			rate,
		)

		rates = append(rates, newRate)
	}

	return rates
}

// Ensure UpdateCurrencyRatesCommandHandler implements the CommandHandler interface
var _ core.CommandHandler[UpdateCurrencyRatesCommand, bool] = (*UpdateCurrencyRatesCommandHandler)(nil)
