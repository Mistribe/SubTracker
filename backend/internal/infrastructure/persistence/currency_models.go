package persistence

import (
	"github.com/oleexo/subtracker/internal/domain/currency"
	"github.com/oleexo/subtracker/internal/infrastructure/persistence/jet/app/public/model"
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
