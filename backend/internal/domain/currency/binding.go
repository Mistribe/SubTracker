package currency

import (
	"golang.org/x/text/currency"
)

type Unit = currency.Unit

var (
	USD = currency.USD
	EUR = currency.EUR
)

func ParseISO(code string) (Unit, error) {
	return currency.ParseISO(code)
}
