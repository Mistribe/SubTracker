package persistence

import (
	"context"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/auth"
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
