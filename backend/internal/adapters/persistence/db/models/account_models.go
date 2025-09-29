package models

import (
	"github.com/mistribe/subtracker/internal/adapters/persistence/db/jet/app/public/model"
	"github.com/mistribe/subtracker/internal/domain/account"
	"github.com/mistribe/subtracker/internal/domain/currency"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/pkg/x"
)

type AccountRow struct {
	model.Accounts
}

func CreateAccountFromJetRow(row AccountRow) (account.Account, error) {
	id := row.ID
	userCurrency, err := currency.ParseISO(row.Currency)
	if err != nil {
		return nil, err
	}
	var plan types.PlanID
	if row.Plan != nil {
		p, pErr := types.ParsePlan(*row.Plan)
		if pErr != nil {
			return nil, pErr
		}
		plan = p
	} else {
		plan = types.PlanFree
	}
	var familyID *types.FamilyID
	if row.FamilyID != nil {
		familyID = x.P(types.FamilyID(*row.FamilyID))
	}
	role := types.ParseRoleOrDefault(row.Role, types.RoleUser)

	return account.New(types.UserID(id),
		userCurrency,
		plan,
		role,
		familyID), nil
}
