package user

import (
	"github.com/google/uuid"
	"golang.org/x/text/currency"

	"github.com/mistribe/subtracker/pkg/x"
)

type User interface {
	Id() string
	Currency() currency.Unit
	SetCurrency(newCurrency currency.Unit)
	Plan() Plan
	Stats() Stats
}

type profile struct {
	id        string
	currency  currency.Unit
	plan      *Plan
	familyIds []uuid.UUID
	stats     Stats
}

func New(id string,
	currency currency.Unit,
	plan *Plan,
	familyIds []uuid.UUID,
	stats Stats) User {
	return &profile{
		id:        id,
		currency:  currency,
		plan:      plan,
		familyIds: familyIds,
		stats:     stats,
	}
}

func NewDefault(id string) User {
	return New(id, currency.USD, x.P(PlanFree), nil, NewEmptyStats())
}

func (u *profile) Stats() Stats {
	return u.stats
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

func (u *profile) Plan() Plan {
	if u.plan == nil {
		return PlanFree
	}

	return *u.plan
}
