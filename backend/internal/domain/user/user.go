package user

import (
	"golang.org/x/text/currency"
)

type Profile interface {
	Id() string
	Currency() currency.Unit
	SetCurrency(newCurrency currency.Unit)
}

type profile struct {
	id       string
	currency currency.Unit
}

func NewProfile(id string, currency currency.Unit) Profile {
	return &profile{
		id:       id,
		currency: currency,
	}
}

func (u *profile) Id() string {
	return u.id
}

func (u *profile) Currency() currency.Unit {
	return u.currency
}
func (u *profile) SetCurrency(newCurrency currency.Unit) {
	u.currency = newCurrency
}
