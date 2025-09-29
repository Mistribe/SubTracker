package account

import (
	"github.com/mistribe/subtracker/internal/domain/currency"
	"github.com/mistribe/subtracker/internal/domain/types"
)

type ConnectedAccount interface {
	UserID() types.UserID
	PlanID() types.PlanID
	Role() types.Role
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
	planID   types.PlanID
	familyID *types.FamilyID
	role     types.Role
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

func (a *account) PlanID() types.PlanID {
	return a.planID
}

func (a *account) FamilyID() *types.FamilyID {
	return a.familyID
}

func (a *account) Role() types.Role {
	return a.role
}

func New(userID types.UserID,
	currency currency.Unit,
	planID types.PlanID,
	role types.Role,
	familyID *types.FamilyID) Account {
	return &account{
		userID:   userID,
		currency: currency,
		planID:   planID,
		familyID: familyID,
		role:     role,
	}
}
