package persistence

import (
	"context"
	dsql "database/sql"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/text/currency"

	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/internal/infrastructure/persistence/sql"
)

type AuthenticationRepository struct {
	dbContext *DatabaseContext
}

func NewAuthenticationRepository(repository *DatabaseContext) auth.Repository {
	return &AuthenticationRepository{dbContext: repository}
}

func (r AuthenticationRepository) GetUserFamilies(ctx context.Context, userId string) ([]uuid.UUID, error) {
	return r.dbContext.GetQueries(ctx).GetUserFamilyIds(ctx, &userId)
}

func (r AuthenticationRepository) GetUserProfile(ctx context.Context, userId string) (auth.UserProfile, error) {
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

	return auth.NewUserProfile(profile.ID, userCurrency), nil
}

func (r AuthenticationRepository) SaveProfile(ctx context.Context, profile auth.UserProfile) error {
	err := r.dbContext.GetQueries(ctx).CreateUserProfile(ctx, sql.CreateUserProfileParams{
		ID:       profile.Id(),
		Currency: profile.Currency().String(),
	})

	return err
}
