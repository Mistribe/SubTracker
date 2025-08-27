package exchange

import (
	"time"

	"golang.org/x/text/currency"
)

func getLiveFakeData() LiveExchangeRates {
	return LiveExchangeRates{
		Timestamp: time.Now(),
		Rates: map[currency.Unit]float64{
			currency.USD: 1.0,
			currency.EUR: 0.85731,
			currency.JPY: 147.731501,
			currency.GBP: 0.74336,
			currency.CHF: 0.807125,
			currency.AUD: 1.531851,
			currency.NZD: 1.67905,
			currency.CAD: 1.374865,
			currency.SEK: 9.56854,
			currency.NOK: 10.27247,
		},
	}
}
