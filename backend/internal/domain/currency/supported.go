package currency

import "golang.org/x/text/currency"

func GetSupportedCurrencies() []currency.Unit {
	return []currency.Unit{
		currency.USD,
		currency.EUR,
		currency.JPY,
		currency.GBP,
		currency.CHF,
		currency.AUD,
		currency.NZD,
		currency.CAD,
		currency.SEK,
		currency.NOK,
	}
}
