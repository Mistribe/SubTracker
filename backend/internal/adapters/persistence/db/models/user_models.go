package models

import (
	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/adapters/persistence/db/jet/app/public/model"
	"github.com/mistribe/subtracker/internal/domain/currency"
	"github.com/mistribe/subtracker/internal/domain/user"
)

type UserRow struct {
	model.Users

	FamilyID                *uuid.UUID `alias:"families.id"`
	ActiveSubscriptionCount int64      `alias:"active_subscription_count"`
	CustomLabelCount        int64      `alias:"custom_label_count"`
	CustomProviderCount     int64      `alias:"custom_provider_count"`
}

func CreateUserFromJetRow(rows []UserRow) (user.User, error) {
	if len(rows) == 0 {
		return nil, nil
	}
	var familyIds []uuid.UUID
	for _, row := range rows {
		if row.FamilyID != nil {
			familyIds = append(familyIds, *row.FamilyID)
		}
	}
	id := rows[0].ID
	userCurrency, err := currency.ParseISO(rows[0].Currency)
	if err != nil {
		return nil, err
	}
	var plan *user.Plan
	if rows[0].Plan != nil {
		p, pErr := user.ParsePlan(*rows[0].Plan)
		if pErr != nil {
			return nil, pErr
		}
		plan = &p
	}

	return user.New(id,
		userCurrency,
		plan,
		familyIds), nil
}
