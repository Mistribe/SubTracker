package persistence

import (
	"context"
	dsql "database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/oleexo/subtracker/internal/domain/currency"
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

func (r CurrencyRateRepository) GetLatestUpdateDate(ctx context.Context) (time.Time, error) {
	row, err := r.dbContext.GetQueries(ctx).GetLatestUpdateDate(ctx)
	if err != nil {
		return time.Time{}, err
	}
	return row, nil
}

func (r CurrencyRateRepository) GetById(ctx context.Context, entityId uuid.UUID) (currency.Rate, error) {
	row, err := r.dbContext.GetQueries(ctx).GetCurrencyRateById(ctx, entityId)
	if err != nil {
		if errors.Is(err, dsql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return createCurrencyRateFromSqlc(row.CurrencyRate), nil
}

func (r CurrencyRateRepository) GetRatesByDate(ctx context.Context, date time.Time) (currency.Rates, error) {
	rows, err := r.dbContext.GetQueries(ctx).GetCurrencyRatesByDate(ctx, pgtype.Date{
		Time:  date,
		Valid: true,
	})
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, nil
	}

	return slicesx.Select(rows, func(in sql.GetCurrencyRatesByDateRow) currency.Rate {
		return createCurrencyRateFromSqlc(in.CurrencyRate)
	}), nil
}

// Save saves one or more currency rates
func (r CurrencyRateRepository) Save(ctx context.Context, rates ...currency.Rate) error {
	var newRates []currency.Rate
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
func (r CurrencyRateRepository) create(ctx context.Context, rates []currency.Rate) error {
	if len(rates) == 0 {
		return nil
	}

	rateParams := slicesx.Select(rates, func(rate currency.Rate) sql.CreateCurrencyRatesParams {
		return sql.CreateCurrencyRatesParams{
			ID:           rate.Id(),
			FromCurrency: rate.FromCurrency().String(),
			ToCurrency:   rate.ToCurrency().String(),
			RateDate:     pgtype.Date{Time: rate.RateDate(), Valid: true},
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
func (r CurrencyRateRepository) update(ctx context.Context, rate currency.Rate) error {
	if !rate.IsDirty() {
		return nil
	}

	err := r.dbContext.GetQueries(ctx).UpdateCurrencyRate(ctx, sql.UpdateCurrencyRateParams{
		ID:           rate.Id(),
		FromCurrency: rate.FromCurrency().String(),
		ToCurrency:   rate.ToCurrency().String(),
		RateDate:     pgtype.Date{Time: rate.RateDate(), Valid: true},
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
