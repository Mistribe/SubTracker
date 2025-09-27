package account

import (
	"github.com/mistribe/subtracker/internal/domain/billing"
	"github.com/mistribe/subtracker/internal/domain/currency"
	"github.com/mistribe/subtracker/internal/domain/family"
)

type ConnectedAccount interface {
	UserID() UserID
	PlanID() billing.PlanID
	Role() Role
}

type Account interface {
	ConnectedAccount

	Id() string
	Currency() currency.Unit
	SetCurrency(newCurrency currency.Unit)
	FamilyID() *family.TypeID
}

type account struct {
	userID   UserID
	currency currency.Unit
	planID   billing.PlanID
	familyID *family.TypeID
	role     Role
}

func (a *account) Id() string {
	return a.userID.String()
}

func (a *account) UserID() UserID {
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

func (a *account) FamilyID() *family.TypeID {
	return a.familyID
}

func (a *account) Role() Role {
	return a.role
}

func New(userID UserID,
	currency currency.Unit,
	planID billing.PlanID,
	role Role,
	familyID *family.TypeID) Account {
	return &account{
		userID:   userID,
		currency: currency,
		planID:   planID,
		familyID: familyID,
		role:     role,
	}
}
