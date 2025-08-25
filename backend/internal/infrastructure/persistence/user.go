package persistence

import (
	"context"
	"errors"

	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
	"golang.org/x/text/currency"

	"github.com/mistribe/subtracker/internal/domain/user"
	"github.com/mistribe/subtracker/internal/infrastructure/persistence/jet/app/public/model"

	. "github.com/go-jet/jet/v2/postgres"

	. "github.com/mistribe/subtracker/internal/infrastructure/persistence/jet/app/public/table"
)

type UserRepository struct {
	dbContext *DatabaseContext
}

func NewUserRepository(repository *DatabaseContext) user.Repository {
	return &UserRepository{dbContext: repository}
}

func (r UserRepository) GetUserFamilies(ctx context.Context, userId string) ([]uuid.UUID, error) {
	stmt := SELECT(FamilyMembers.FamilyID).
		FROM(FamilyMembers).
		WHERE(FamilyMembers.UserID.EQ(String(userId)))

	var rows []struct {
		FamilyID uuid.UUID `json:"family_id"`
	}

	if err := r.dbContext.Query(ctx, stmt, &rows); err != nil {
		return nil, err
	}

	familyIds := make([]uuid.UUID, len(rows))
	for i, row := range rows {
		familyIds[i] = row.FamilyID
	}

	return familyIds, nil
}

func (r UserRepository) GetUserProfile(ctx context.Context, userId string) (user.Profile, error) {
	stmt := SELECT(Users.AllColumns).
		FROM(Users).
		WHERE(Users.ID.EQ(String(userId)))

	var profile model.Users
	if err := r.dbContext.Query(ctx, stmt, &profile); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if profile.ID == "" {
		return nil, nil
	}

	userCurrency, err := currency.ParseISO(profile.Currency)
	if err != nil {
		return nil, err
	}

	return user.NewProfile(profile.ID, userCurrency), nil
}

func (r UserRepository) SaveProfile(ctx context.Context, profile user.Profile) error {
	stmt := Users.
		INSERT(Users.ID, Users.Currency).
		VALUES(String(profile.Id()), String(profile.Currency().String())).
		ON_CONFLICT(Users.ID).
		DO_UPDATE(
			SET(Users.Currency.SET(Users.EXCLUDED.Currency)),
		)

	_, err := r.dbContext.Execute(ctx, stmt)
	return err
}
