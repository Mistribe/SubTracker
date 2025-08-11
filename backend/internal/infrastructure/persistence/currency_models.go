package persistence

import (
	"github.com/oleexo/subtracker/internal/domain/currency"
	"github.com/oleexo/subtracker/internal/infrastructure/persistence/sql"
)

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
