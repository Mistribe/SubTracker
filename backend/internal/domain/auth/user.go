package auth

import (
	"golang.org/x/text/currency"
)

type UserProfile interface {
	Id() string
	Currency() currency.Unit
	SetCurrency(newCurrency currency.Unit)
}

type userProfile struct {
	id       string
	currency currency.Unit
}

func NewUserProfile(id string, currency currency.Unit) UserProfile {
	return &userProfile{
		id:       id,
		currency: currency,
	}
}

func (u *userProfile) Id() string {
	return u.id
}

func (u *userProfile) Currency() currency.Unit {
	return u.currency
}
func (u *userProfile) SetCurrency(newCurrency currency.Unit) {
	u.currency = newCurrency
}
