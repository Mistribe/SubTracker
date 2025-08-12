package currency

import (
	"golang.org/x/text/currency"
)

type Amount struct {
	value    float64
	currency Unit
	isValid  bool
}

func (a Amount) Value() float64 {
	return a.value
}

func (a Amount) Currency() Unit {
	return a.currency
}

func (a Amount) IsValid() bool {
	return a.isValid
}

func (a Amount) ToCurrency(
	to currency.Unit,
	rates Rates) Amount {
	if !a.isValid {
		return a
	}
	if a.currency == to {
		return a
	}

	var usdValue float64
	if a.currency == USD {
		usdValue = a.value
	} else {
		usdValue = a.ToUSD(rates).Value()
	}
	exchangeRate, ok := rates.FindExchangeRate(USD, to)
	if !ok {
		return NewInvalidAmount()
	}

	currencyValue := usdValue * exchangeRate
	return NewAmount(currencyValue, to)
}

func (a Amount) ToUSD(rates Rates) Amount {
	if !a.isValid {
		return a
	}

	if a.currency == USD {
		return a
	}

	exchangeRate, ok := rates.FindExchangeRate(a.currency, USD)
	if !ok {
		return NewInvalidAmount()
	}

	return NewAmount(a.value*exchangeRate, USD)
}

func NewAmount(value float64, unit Unit) Amount {
	return Amount{
		value:    value,
		currency: unit,
		isValid:  true,
	}
}

func NewInvalidAmount() Amount {
	return Amount{
		isValid: false,
	}
}
