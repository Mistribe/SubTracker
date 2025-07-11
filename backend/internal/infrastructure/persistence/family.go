package persistence

import (
	"context"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/application/core/option"
	"github.com/oleexo/subtracker/internal/domain/family"
)

type FamilyRepository struct {
}

func NewFamilyRepository() *FamilyRepository {
	return &FamilyRepository{}
}

func (r FamilyRepository) Get(ctx context.Context, id uuid.UUID) (option.Option[family.Member], error) {
	//TODO implement me
	panic("implement me")
}

func (r FamilyRepository) GetAll(ctx context.Context) ([]family.Member, error) {
	//TODO implement me
	panic("implement me")
}

func (r FamilyRepository) Save(ctx context.Context, member family.Member) error {
	//TODO implement me
	panic("implement me")
}

func (r FamilyRepository) Delete(ctx context.Context, id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (r FamilyRepository) Exists(ctx context.Context, members ...uuid.UUID) (bool, error) {
	//TODO implement me
	panic("implement me")
}
