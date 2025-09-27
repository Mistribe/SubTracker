package repositories

import (
	"context"

	"github.com/mistribe/subtracker/internal/adapters/persistence/db"
	. "github.com/mistribe/subtracker/internal/adapters/persistence/db/jet/app/public/table"
	"github.com/mistribe/subtracker/internal/adapters/persistence/db/models"
	"github.com/mistribe/subtracker/internal/domain/account"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/x"

	. "github.com/go-jet/jet/v2/postgres"
)

type AccountRepository struct {
	dbContext *db.Context
}

func NewAccountRepository(repository *db.Context) ports.AccountRepository {
	return &AccountRepository{dbContext: repository}
}

func (r AccountRepository) GetFamily(ctx context.Context, userId types.UserID) (*types.FamilyID, error) {
	stmt := SELECT(Accounts.FamilyID).
		FROM(Accounts).
		WHERE(Accounts.ID.EQ(String(userId.String()))).
		LIMIT(1)

	var rows []models.AccountRow
	if err := r.dbContext.Query(ctx, stmt, &rows); err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, nil
	}

	familyID := types.FamilyID(*rows[0].FamilyID)
	return &familyID, nil
}

func (r AccountRepository) GetById(ctx context.Context, userId types.UserID) (account.Account, error) {
	stmt := SELECT(Accounts.AllColumns).
		FROM(Accounts).
		WHERE(Accounts.ID.EQ(String(userId.String()))).
		LIMIT(1)

	var rows []models.AccountRow
	if err := r.dbContext.Query(ctx, stmt, &rows); err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, nil
	}

	return models.CreateAccountFromJetRow(rows[0])
}

func (r AccountRepository) Save(ctx context.Context, acc account.Account) error {
	stmt := Accounts.
		INSERT(
			Accounts.ID,
			Accounts.Currency,
			Accounts.Plan,
			Accounts.Role,
			Accounts.FamilyID,
		).
		VALUES(
			String(acc.Id()),
			String(acc.Currency().String()),
			String(acc.PlanID().String()),
			String(acc.Role().String()),
			x.TernaryFunc(
				acc.FamilyID() == nil,
				func() Expression {
					return NULL
				},
				func() Expression {
					return String(acc.FamilyID().String())
				}),
		)

	_, err := r.dbContext.Execute(ctx, stmt)
	return err
}
