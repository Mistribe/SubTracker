package repositories

import (
	"context"

	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/adapters/persistence/db"
	. "github.com/mistribe/subtracker/internal/adapters/persistence/db/jet/app/public/table"
	"github.com/mistribe/subtracker/internal/adapters/persistence/db/models"
	"github.com/mistribe/subtracker/internal/domain/user"
	"github.com/mistribe/subtracker/internal/ports"

	. "github.com/go-jet/jet/v2/postgres"
)

type UserRepository struct {
	dbContext *db.Context
}

func NewUserRepository(repository *db.Context) ports.UserRepository {
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

func (r UserRepository) GetUser(ctx context.Context, userId string) (user.User, error) {
	active := NOW().GT_EQ(Subscriptions.StartDate).
		AND(Subscriptions.EndDate.IS_NULL().OR(NOW().LT_EQ(Subscriptions.EndDate)))

	activeSubCount := SELECT(COUNT(Subscriptions.ID)).
		FROM(Subscriptions).
		WHERE(
			active.
				AND(
					Subscriptions.OwnerType.EQ(String("personal")).
						AND(Subscriptions.OwnerUserID.EQ(Users.ID)),
				).
				OR(
					Subscriptions.OwnerType.EQ(String("family")).AND(EXISTS(
						SELECT(FamilyMembers.ID).
							FROM(FamilyMembers).
							WHERE(
								FamilyMembers.FamilyID.EQ(Subscriptions.OwnerFamilyID).
									AND(FamilyMembers.UserID.EQ(String(userId))),
							),
					)),
				),
		)

	customLabels := SELECT(COUNT(Labels.ID)).
		FROM(Labels).
		WHERE(
			Labels.OwnerType.EQ(String("personal")).AND(Labels.OwnerUserID.EQ(Users.ID)).
				OR(
					Labels.OwnerType.EQ(String("family")).AND(EXISTS(
						SELECT(FamilyMembers.ID).
							FROM(FamilyMembers).
							WHERE(
								FamilyMembers.FamilyID.EQ(Labels.OwnerFamilyID).
									AND(FamilyMembers.UserID.EQ(String(userId))),
							),
					)),
				),
		)

	customProviders := SELECT(COUNT(Providers.ID)).
		FROM(Providers).
		WHERE(
			Providers.OwnerType.EQ(String("personal")).AND(Providers.OwnerUserID.EQ(Users.ID)).
				OR(
					Providers.OwnerType.EQ(String("family")).AND(EXISTS(
						SELECT(FamilyMembers.ID).
							FROM(FamilyMembers).
							WHERE(
								FamilyMembers.FamilyID.EQ(Providers.OwnerFamilyID).
									AND(FamilyMembers.UserID.EQ(String(userId))),
							),
					)),
				),
		)

	// todo count family members

	stmt := SELECT(
		Users.AllColumns,
		Families.ID,
		activeSubCount.AS("user_row.active_subscription_count"),
		customLabels.AS("user_row.custom_label_count"),
		customProviders.AS("user_row.custom_provider_count"),
	).
		FROM(
			Users.
				LEFT_JOIN(FamilyMembers, FamilyMembers.UserID.EQ(Users.ID)).
				LEFT_JOIN(Families, Families.ID.EQ(FamilyMembers.FamilyID)),
		).
		WHERE(Users.ID.EQ(String(userId)))

	var rows []models.UserRow
	if err := r.dbContext.Query(ctx, stmt, &rows); err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		// We need to create the default user
		if _, err := r.CreateDefault(ctx, userId); err != nil {
			return nil, err
		}
		return r.GetUser(ctx, userId)
	}

	return models.CreateUserFromJetRow(rows)
}

func (r UserRepository) CreateDefault(ctx context.Context, userId string) (user.User, error) {
	defaultUser := user.NewDefault(userId)
	return defaultUser, r.Save(ctx, defaultUser)
}

func (r UserRepository) Save(ctx context.Context, usr user.User) error {
	stmt := Users.
		INSERT(Users.ID, Users.Currency).
		VALUES(
			String(usr.Id()),
			String(usr.Currency().String()),
			String(usr.Plan().String())).
		ON_CONFLICT(Users.ID).
		DO_UPDATE(
			SET(Users.Currency.SET(Users.EXCLUDED.Currency)),
		)

	_, err := r.dbContext.Execute(ctx, stmt)
	return err
}
