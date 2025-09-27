package account

import (
	"github.com/mistribe/subtracker/internal/domain/billing"
	"github.com/mistribe/subtracker/internal/domain/currency"
	"github.com/mistribe/subtracker/internal/domain/types"
)

type ConnectedAccount interface {
	UserID() types.UserID
	PlanID() billing.PlanID
	Role() Role
}

type Account interface {
	ConnectedAccount

	Id() string
	Currency() currency.Unit
	SetCurrency(newCurrency currency.Unit)
	FamilyID() *types.FamilyID
}

type account struct {
	userID   types.UserID
	currency currency.Unit
	planID   billing.PlanID
	familyID *types.FamilyID
	role     Role
}

func (a *account) Id() string {
	return a.userID.String()
}

func (a *account) UserID() types.UserID {
	return a.userID
}

func (a *account) Currency() currency.Unit {
	return a.currency
}

func (a *account) SetCurrency(newCurrency currency.Unit) {
	a.currency = newCurrency
}

func (a *account) PlanID() billing.PlanID {
	return a.planID
}

func (a *account) FamilyID() *types.FamilyID {
	return a.familyID
}

func (a *account) Role() Role {
	return a.role
}

func New(userID types.UserID,
	currency currency.Unit,
	planID billing.PlanID,
	role Role,
	familyID *types.FamilyID) Account {
	return &account{
		userID:   userID,
		currency: currency,
		planID:   planID,
		familyID: familyID,
		role:     role,
	}
}
