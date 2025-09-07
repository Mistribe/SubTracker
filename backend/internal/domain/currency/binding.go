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

func MustParseISO(code string) Unit {
	return currency.MustParseISO(code)
}
