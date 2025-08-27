package models

import (
	"github.com/mistribe/subtracker/internal/adapters/persistence/db/jet/app/public/model"
	"github.com/mistribe/subtracker/internal/domain/currency"
)

func CreateCurrencyRateFromModel(source model.CurrencyRates) (currency.Rate, error) {
	fromCurrency, err := currency.ParseISO(source.FromCurrency)
	if err != nil {
		return nil, err
	}
	toCurrency, err := currency.ParseISO(source.ToCurrency)
	if err != nil {
		return nil, err
	}

	return currency.NewRate(
		source.ID,
		fromCurrency,
		toCurrency,
		source.RateDate,
		source.ExchangeRate,
		source.CreatedAt,
		source.UpdatedAt,
	), nil
}
