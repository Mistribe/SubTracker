package persistence

import (
	"context"

	"github.com/google/uuid"
	"github.com/oleexo/subtracker/internal/domain/auth"
)

type AuthenticationRepository struct {
	repository *DatabaseContext
}

func NewAuthenticationRepository(repository *DatabaseContext) auth.Repository {
	return &AuthenticationRepository{repository: repository}
}

func (r AuthenticationRepository) GetUserFamilies(ctx context.Context, userId string) ([]uuid.UUID, error) {
	var familyIds []uuid.UUID
	result := r.repository.db.WithContext(ctx).
		Model(&FamilyMemberSqlModel{}).
		Where("user_id = ?", userId).
		Select("id").
		Scan(&familyIds)
	if result.Error != nil {
		return nil, result.Error
	}

	return familyIds, nil
}
