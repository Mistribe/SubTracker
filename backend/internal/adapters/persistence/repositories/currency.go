package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/adapters/persistence/db"
	"github.com/mistribe/subtracker/internal/adapters/persistence/db/jet/app/public/model"
	. "github.com/mistribe/subtracker/internal/adapters/persistence/db/jet/app/public/table"
	"github.com/mistribe/subtracker/internal/adapters/persistence/db/models"
	"github.com/mistribe/subtracker/internal/domain/currency"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/slicesx"

	. "github.com/go-jet/jet/v2/postgres"
)

// CurrencyRateRepository implements the currency.UserRepository interface
type CurrencyRateRepository struct {
	dbContext *db.Context
}

// NewCurrencyRateRepository creates a new CurrencyRateRepository
func NewCurrencyRateRepository(dbContext *db.Context) ports.CurrencyRepository {
	return &CurrencyRateRepository{
		dbContext: dbContext,
	}
}

func (r CurrencyRateRepository) GetLatestUpdateDate(ctx context.Context) (time.Time, error) {
	stmt := SELECT(
		MAX(CurrencyRates.UpdatedAt).AS("latest_update_date"),
	).FROM(CurrencyRates)
	var row struct {
		LatestUpdateDate *time.Time `json:"latest_update_date"`
	}
	if err := r.dbContext.Query(ctx, stmt, &row); err != nil {
		return time.Time{}, err
	}
	if row.LatestUpdateDate == nil {
		return time.Time{}, nil
	}
	return *row.LatestUpdateDate, nil
}

func (r CurrencyRateRepository) GetById(ctx context.Context, entityId uuid.UUID) (currency.Rate, error) {
	stmt := SELECT(CurrencyRates.AllColumns).
		FROM(CurrencyRates).
		WHERE(CurrencyRates.ID.EQ(UUID(entityId)))
	var row model.CurrencyRates
	if err := r.dbContext.Query(ctx, stmt, &row); err != nil {
		return nil, err
	}
	return models.CreateCurrencyRateFromModel(row)
}

func (r CurrencyRateRepository) GetRatesByDate(ctx context.Context, date time.Time) (currency.Rates, error) {
	stmt := SELECT(CurrencyRates.AllColumns).
		FROM(CurrencyRates).
		WHERE(CurrencyRates.RateDate.EQ(DateT(date))).
		ORDER_BY(CurrencyRates.FromCurrency.ASC(), CurrencyRates.ToCurrency.ASC())

	var rows []model.CurrencyRates
	if err := r.dbContext.Query(ctx, stmt, &rows); err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, nil
	}

	return slicesx.SelectErr(rows, func(in model.CurrencyRates) (currency.Rate, error) {
		return models.CreateCurrencyRateFromModel(in)
	})
}

func (r CurrencyRateRepository) GetRateAt(ctx context.Context, from, to currency.Unit, at time.Time) (
	currency.Rate,
	error) {
	stmt := SELECT(CurrencyRates.AllColumns).
		FROM(CurrencyRates).
		WHERE(CurrencyRates.RateDate.EQ(DateT(at)).
			AND(CurrencyRates.FromCurrency.EQ(String(from.String()))).
			AND(CurrencyRates.ToCurrency.EQ(String(to.String())))).
		LIMIT(1)

	var rows []model.CurrencyRates
	if err := r.dbContext.Query(ctx, stmt, &rows); err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, nil
	}

	return models.CreateCurrencyRateFromModel(rows[0])
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

	stmt := CurrencyRates.
		INSERT(
			CurrencyRates.ID,
			CurrencyRates.FromCurrency,
			CurrencyRates.ToCurrency,
			CurrencyRates.RateDate,
			CurrencyRates.ExchangeRate,
			CurrencyRates.CreatedAt,
			CurrencyRates.UpdatedAt,
			CurrencyRates.Etag,
		)

	for _, rate := range rates {
		stmt = stmt.VALUES(UUID(rate.Id()),
			String(rate.FromCurrency().String()),
			String(rate.ToCurrency().String()),
			DateT(rate.RateDate()),
			Float(rate.ExchangeRate()),
			TimestampzT(rate.CreatedAt()),
			TimestampzT(rate.UpdatedAt()),
			String(rate.ETag()),
		)
	}

	count, err := r.dbContext.Execute(ctx, stmt)
	if err != nil {
		return err
	}
	if count != int64(len(rates)) {
		return db.ErrMissMatchAffectRow
	}
	return nil
}

// update updates an existing currency rate
func (r CurrencyRateRepository) update(ctx context.Context, rate currency.Rate) error {
	if !rate.IsDirty() {
		return nil
	}

	stmt := CurrencyRates.
		UPDATE().
		SET(
			CurrencyRates.FromCurrency.SET(String(rate.FromCurrency().String())),
			CurrencyRates.ToCurrency.SET(String(rate.ToCurrency().String())),
			CurrencyRates.RateDate.SET(DateT(rate.RateDate())),
			CurrencyRates.ExchangeRate.SET(Float(rate.ExchangeRate())),
			CurrencyRates.UpdatedAt.SET(TimestampT(rate.UpdatedAt())),
			CurrencyRates.Etag.SET(String(rate.ETag())),
		).
		WHERE(CurrencyRates.ID.EQ(UUID(rate.Id())))

	count, err := r.dbContext.Execute(ctx, stmt)
	if err != nil {
		return err
	}
	if count == 0 {
		return db.ErrMissMatchAffectRow
	}
	return nil
}

// Delete deletes a currency rate by its ID
func (r CurrencyRateRepository) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	stmt := CurrencyRates.
		DELETE().
		WHERE(CurrencyRates.ID.EQ(UUID(id)))

	count, err := r.dbContext.Execute(ctx, stmt)
	if err != nil {
		return false, err
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
}

// Exists checks if currency rates with the given IDs exist
func (r CurrencyRateRepository) Exists(ctx context.Context, ids ...uuid.UUID) (bool, error) {
	if len(ids) == 0 {
		return true, nil
	}

	vals := make([]Expression, len(ids))
	for i, id := range ids {
		vals[i] = UUID(id)
	}

	stmt := SELECT(COUNT(CurrencyRates.ID).AS("count")).
		FROM(CurrencyRates).
		WHERE(CurrencyRates.ID.IN(vals...))

	var row struct {
		Count int
	}

	if err := r.dbContext.Query(ctx, stmt, &row); err != nil {
		return false, err
	}

	return row.Count == len(ids), nil
}
