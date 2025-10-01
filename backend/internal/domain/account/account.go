package account

import (
	"time"

	"github.com/mistribe/subtracker/internal/domain/currency"
	"github.com/mistribe/subtracker/internal/domain/entity"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/pkg/langext/option"
)

type ConnectedAccount interface {
	UserID() types.UserID
	PlanID() types.PlanID
	Role() types.Role
}

type Account interface {
	ConnectedAccount
	entity.Entity[string]
	entity.ETagEntity

	Currency() option.Option[currency.Unit]
	SetCurrency(newCurrency currency.Unit)
	FamilyID() *types.FamilyID
}

type account struct {
	*entity.Base[string]

	currency *currency.Unit
	planID   types.PlanID
	familyID *types.FamilyID
	role     types.Role
}

func (a *account) UserID() types.UserID {
	return types.UserID(a.Base.Id())
}

func (a *account) Currency() option.Option[currency.Unit] {
	return option.New(a.currency)
}

func (a *account) SetCurrency(newCurrency currency.Unit) {
	a.currency = &newCurrency
	a.SetAsDirty()
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

func (a *account) Equal(other Account) bool {
	if other == nil {
		return false
	}

	return a.ETag() == other.ETag()
}

func (a *account) ETagFields() []interface{} {
	fields := []interface{}{
		a.PlanID().String(),
		a.Role().String(),
	}

	a.Currency().IfSome(func(c currency.Unit) {
		fields = append(fields, c.String())
	})

	if a.FamilyID() != nil {
		fields = append(fields, a.FamilyID().String())
	}

	return fields
}
func (a *account) ETag() string {
	return entity.CalculateETag(a)
}

func New(
	userID types.UserID,
	currency *currency.Unit,
	planID types.PlanID,
	role types.Role,
	familyID *types.FamilyID,
	createdAt time.Time,
	updatedAt time.Time) Account {
	return &account{
		Base:     entity.NewBase[string](userID.String(), createdAt, updatedAt, true, false),
		currency: currency,
		planID:   planID,
		familyID: familyID,
		role:     role,
	}
}
