package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	currencyunit "golang.org/x/text/currency"

	"github.com/oleexo/subtracker/internal/domain/currency"
	"github.com/oleexo/subtracker/internal/domain/entity"
	"github.com/oleexo/subtracker/internal/infrastructure/persistence/sql"
	"github.com/oleexo/subtracker/pkg/slicesx"
)

// CurrencyRateRepository implements the currency.Repository interface
type CurrencyRateRepository struct {
	dbContext *DatabaseContext
}

// NewCurrencyRateRepository creates a new CurrencyRateRepository
func NewCurrencyRateRepository(dbContext *DatabaseContext) currency.Repository {
	return &CurrencyRateRepository{
		dbContext: dbContext,
	}
}

// GetById retrieves a currency rate by its ID
func (r CurrencyRateRepository) GetById(ctx context.Context, id uuid.UUID) (*currency.Rate, error) {
	rate, err := r.dbContext.GetQueries(ctx).GetCurrencyRateById(ctx, id)
	if err != nil {
		return nil, err
	}

	fromCurrency, err := currencyunit.ParseISO(rate.FromCurrency)
	if err != nil {
		return nil, err
	}

	toCurrency, err := currencyunit.ParseISO(rate.ToCurrency)
	if err != nil {
		return nil, err
	}

	return currency.FromExisting(
		rate.ID,
		fromCurrency,
		toCurrency,
		rate.RateDate,
		rate.ExchangeRate,
		rate.CreatedAt,
		rate.UpdatedAt,
		rate.Etag,
	), nil
}

// GetRateByDate retrieves a currency rate for a specific currency pair and date
func (r CurrencyRateRepository) GetRateByDate(
	ctx context.Context,
	fromCurrency, toCurrency currency.Unit,
	date time.Time) (*currency.Rate, error) {
	rate, err := r.dbContext.GetQueries(ctx).GetCurrencyRateByDate(ctx, sql.GetCurrencyRateByDateParams{
		FromCurrency: fromCurrency.String(),
		ToCurrency:   toCurrency.String(),
		RateDate:     date,
	})
	if err != nil {
		return nil, err
	}

	fromCur, err := currency.ParseISO(rate.FromCurrency)
	if err != nil {
		return nil, err
	}

	toCur, err := currency.ParseISO(rate.ToCurrency)
	if err != nil {
		return nil, err
	}

	return currency.FromExisting(
		rate.ID,
		fromCur,
		toCur,
		rate.RateDate,
		rate.ExchangeRate,
		rate.CreatedAt,
		rate.UpdatedAt,
		rate.Etag,
	), nil
}

// GetLatestRate retrieves the most recent currency rate for a specific currency pair
func (r CurrencyRateRepository) GetLatestRate(
	ctx context.Context,
	fromCurrency, toCurrency currency.Unit) (*currency.Rate, error) {
	rate, err := r.dbContext.GetQueries(ctx).GetLatestCurrencyRate(ctx, sql.GetLatestCurrencyRateParams{
		FromCurrency: fromCurrency.String(),
		ToCurrency:   toCurrency.String(),
	})
	if err != nil {
		return nil, err
	}

	fromCur, err := currency.ParseISO(rate.FromCurrency)
	if err != nil {
		return nil, err
	}

	toCur, err := currency.ParseISO(rate.ToCurrency)
	if err != nil {
		return nil, err
	}

	return currency.FromExisting(
		rate.ID,
		fromCur,
		toCur,
		rate.RateDate,
		rate.ExchangeRate,
		rate.CreatedAt,
		rate.UpdatedAt,
		rate.Etag,
	), nil
}

// ConvertAmount converts an amount from one currency to another at a specific date
func (r CurrencyRateRepository) ConvertAmount(
	ctx context.Context,
	amount float64,
	fromCurrency, toCurrency currency.Unit,
	date time.Time) (float64, error) {
	// If currencies are the same, return the original amount
	if fromCurrency == toCurrency {
		return amount, nil
	}

	// Try to get direct conversion rate
	rate, err := r.GetRateByDate(ctx, fromCurrency, toCurrency, date)
	if err == nil {
		return amount * rate.ExchangeRate(), nil
	}

	// Try reverse conversion
	rate, err = r.GetRateByDate(ctx, toCurrency, fromCurrency, date)
	if err == nil {
		return amount / rate.ExchangeRate(), nil
	}

	// If no direct or reverse rate, try to convert through USD
	if fromCurrency != currency.USD && toCurrency != currency.USD {
		// Convert from source currency to USD
		usdAmount, err := r.ConvertAmount(ctx, amount, fromCurrency, currency.USD, date)
		if err != nil {
			return 0, err
		}

		// Convert from USD to target currency
		return r.ConvertAmount(ctx, usdAmount, currency.USD, toCurrency, date)
	}

	return 0, errors.New("no conversion rate available")
}

// Save saves one or more currency rates
func (r CurrencyRateRepository) Save(ctx context.Context, rates ...*currency.Rate) error {
	var newRates []*currency.Rate
	for _, rate := range rates {
		if !rate.IsExists() {
			newRates = append(newRates, rate)
		} else {
			if err := r.update(ctx, rate); err != nil {
				return err
			}
		}
	}

	if len(newRates) > 0 {
		if err := r.create(ctx, newRates); err != nil {
			return err
		}
	}

	for _, rate := range rates {
		rate.Clean()
	}

	return nil
}

// create creates new currency rates
func (r CurrencyRateRepository) create(ctx context.Context, rates []*currency.Rate) error {
	if len(rates) == 0 {
		return nil
	}

	rateParams := slicesx.Select(rates, func(rate *currency.Rate) sql.CreateCurrencyRatesParams {
		return sql.CreateCurrencyRatesParams{
			ID:           rate.Id(),
			FromCurrency: rate.FromCurrency().String(),
			ToCurrency:   rate.ToCurrency().String(),
			RateDate:     rate.RateDate(),
			ExchangeRate: rate.ExchangeRate(),
			CreatedAt:    rate.CreatedAt(),
			UpdatedAt:    rate.UpdatedAt(),
			Etag:         rate.ETag(),
		}
	})

	_, err := r.dbContext.GetQueries(ctx).CreateCurrencyRates(ctx, rateParams)
	if err != nil {
		return err
	}
	return nil
}

// update updates an existing currency rate
func (r CurrencyRateRepository) update(ctx context.Context, rate *currency.Rate) error {
	if !rate.IsDirty() {
		return nil
	}

	err := r.dbContext.GetQueries(ctx).UpdateCurrencyRate(ctx, sql.UpdateCurrencyRateParams{
		ID:           rate.Id(),
		FromCurrency: rate.FromCurrency().String(),
		ToCurrency:   rate.ToCurrency().String(),
		RateDate:     rate.RateDate(),
		ExchangeRate: rate.ExchangeRate(),
		UpdatedAt:    rate.UpdatedAt(),
		Etag:         rate.ETag(),
	})

	return err
}

// Delete deletes a currency rate by its ID
func (r CurrencyRateRepository) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	err := r.dbContext.GetQueries(ctx).DeleteCurrencyRate(ctx, id)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Exists checks if currency rates with the given IDs exist
func (r CurrencyRateRepository) Exists(ctx context.Context, ids ...uuid.UUID) (bool, error) {
	count, err := r.dbContext.GetQueries(ctx).IsCurrencyRateExists(ctx, ids)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetAll retrieves all currency rates with pagination
func (r CurrencyRateRepository) GetAll(ctx context.Context, parameters entity.QueryParameters) (
	[]*currency.Rate,
	int64,
	error) {
	rates, err := r.dbContext.GetQueries(ctx).GetCurrencyRates(ctx, sql.GetCurrencyRatesParams{
		Limit:  parameters.Limit,
		Offset: parameters.Offset,
	})
	if err != nil {
		return nil, 0, err
	}

	count, err := r.dbContext.GetQueries(ctx).GetCurrencyRatesCount(ctx)
	if err != nil {
		return nil, 0, err
	}

	result := make([]*currency.Rate, 0, len(rates))
	for _, rate := range rates {
		fromCur, err := currency.ParseISO(rate.FromCurrency)
		if err != nil {
			return nil, 0, err
		}

		toCur, err := currency.ParseISO(rate.ToCurrency)
		if err != nil {
			return nil, 0, err
		}

		result = append(result, currency.FromExisting(
			rate.ID,
			fromCur,
			toCur,
			rate.RateDate,
			rate.ExchangeRate,
			rate.CreatedAt,
			rate.UpdatedAt,
			rate.Etag,
		))
	}

	return result, count, nil
}
