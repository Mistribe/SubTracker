package currency

import (
	"golang.org/x/text/currency"
)

type Amount interface {
	IsZero() bool
	IsNegative() bool
	IsPositive() bool
	IsEqual(other Amount) bool
	IsGreaterThan(other Amount) bool
	IsLessThan(other Amount) bool
	Value() float64
	Source() Amount
	Currency() Unit
	IsValid() bool
	ToCurrency(to currency.Unit, rates Rates) Amount
	Add(otherAmount Amount) Amount
}
type amount struct {
	value    float64
	currency Unit
	source   Amount
	isValid  bool
}

func (a amount) Add(otherAmount Amount) Amount {
	if !a.isValid || !otherAmount.IsValid() {
		return NewInvalidAmount()
	}

	if a.currency != otherAmount.Currency() {
		return NewInvalidAmount()
	}
	return amount{
		value:    a.value + otherAmount.Value(),
		currency: a.currency,
		source:   nil,
		isValid:  true,
	}
}

func (a amount) Source() Amount {
	return a.source
}

func (a amount) IsZero() bool {
	return a.value == 0
}

func (a amount) IsNegative() bool {
	return a.value < 0
}

func (a amount) IsPositive() bool {
	return a.value > 0
}

func (a amount) IsEqual(other Amount) bool {
	return a.value == other.Value() &&
		a.currency == other.Currency()
}

func (a amount) IsGreaterThan(other Amount) bool {
	return a.value > other.Value()
}

func (a amount) IsLessThan(other Amount) bool {
	return a.value < other.Value()
}

func (a amount) Value() float64 {
	return a.value
}

func (a amount) Currency() Unit {
	return a.currency
}

func (a amount) IsValid() bool {
	return a.isValid
}

func (a amount) ToCurrency(
	to currency.Unit,
	rates Rates) Amount {
	if !a.isValid {
		return a
	}
	if a.currency == to {
		return a
	}

	var usdAmount Amount
	if a.currency == USD {
		usdAmount = a
	} else {
		usdAmount = a.ToUSD(rates)
	}
	exchangeRate, ok := rates.FindExchangeRate(USD, to)
	if !ok {
		return NewInvalidAmount()
	}

	currencyValue := usdAmount.Value() * exchangeRate
	return amount{
		value:    currencyValue,
		currency: to,
		source:   &a,
		isValid:  true,
	}
}

func (a amount) ToUSD(rates Rates) Amount {
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

	return amount{
		value:    a.value * exchangeRate,
		currency: USD,
		source:   &a,
		isValid:  true,
	}
}

func NewAmount(value float64, unit Unit) Amount {
	return amount{
		value:    value,
		currency: unit,
		isValid:  true,
	}
}

func NewInvalidAmount() Amount {
	return amount{
		isValid: false,
	}
}
