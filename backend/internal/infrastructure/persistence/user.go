package persistence

import (
	"context"
	dsql "database/sql"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/text/currency"

	"github.com/oleexo/subtracker/internal/domain/user"
	"github.com/oleexo/subtracker/internal/infrastructure/persistence/sql"
)

type UserRepository struct {
	dbContext *DatabaseContext
}

func NewUserRepository(repository *DatabaseContext) user.Repository {
	return &UserRepository{dbContext: repository}
}

func (r UserRepository) GetUserFamilies(ctx context.Context, userId string) ([]uuid.UUID, error) {
	return r.dbContext.GetQueries(ctx).GetUserFamilyIds(ctx, &userId)
}

func (r UserRepository) GetUserProfile(ctx context.Context, userId string) (user.Profile, error) {
	profile, err := r.dbContext.GetQueries(ctx).GetUserProfile(ctx, userId)
	if err != nil {
		if errors.Is(err, dsql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	userCurrency, err := currency.ParseISO(profile.Currency)
	if err != nil {
		return nil, err
	}

	return user.NewProfile(profile.ID, userCurrency), nil
}

func (r UserRepository) SaveProfile(ctx context.Context, profile user.Profile) error {
	err := r.dbContext.GetQueries(ctx).CreateUserProfile(ctx, sql.CreateUserProfileParams{
		ID:       profile.Id(),
		Currency: profile.Currency().String(),
	})

	return err
}
