package persistence

import (
	"github.com/oleexo/subtracker/internal/domain/currency"
	"github.com/oleexo/subtracker/internal/infrastructure/persistence/jet/app/public/model"
	"github.com/oleexo/subtracker/internal/infrastructure/persistence/sql"
)

func createCurrencyRateFromJet(jetCurrencyRate model.CurrencyRates) (currency.Rate, error) {
	fromCurrency, err := currency.ParseISO(jetCurrencyRate.FromCurrency)
	if err != nil {
		return nil, err
	}
	toCurrency, err := currency.ParseISO(jetCurrencyRate.ToCurrency)
	if err != nil {
		return nil, err
	}

	return currency.NewRate(
		jetCurrencyRate.ID,
		fromCurrency,
		toCurrency,
		jetCurrencyRate.RateDate,
		jetCurrencyRate.ExchangeRate,
		jetCurrencyRate.CreatedAt,
		jetCurrencyRate.UpdatedAt,
	), nil
}

func createCurrencyRateFromSqlc(sqlcCurrencyRate sql.CurrencyRate) currency.Rate {
	fromCurrency, err := currency.ParseISO(sqlcCurrencyRate.FromCurrency)
	if err != nil {
		panic(err)
	}
	toCurrency, err := currency.ParseISO(sqlcCurrencyRate.ToCurrency)
	if err != nil {
		panic(err)
	}
	return currency.NewRate(
		sqlcCurrencyRate.ID,
		fromCurrency,
		toCurrency,
		sqlcCurrencyRate.RateDate.Time,
		sqlcCurrencyRate.ExchangeRate,
		sqlcCurrencyRate.CreatedAt,
		sqlcCurrencyRate.UpdatedAt,
	)
}
