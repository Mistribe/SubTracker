package persistence

import (
	"context"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/family"
	"github.com/oleexo/subtracker/pkg/langext/option"
)

type FamilyRepository struct {
	members map[uuid.UUID]family.Member
}

func NewFamilyRepository() *FamilyRepository {
	return &FamilyRepository{
		members: make(map[uuid.UUID]family.Member),
	}
}

func (r FamilyRepository) Get(ctx context.Context, id uuid.UUID) (option.Option[family.Member], error) {
	if member, ok := r.members[id]; ok {
		return option.Some(member), nil
	}
	return option.None[family.Member](), nil
}

func (r FamilyRepository) GetAll(ctx context.Context) ([]family.Member, error) {
	members := make([]family.Member, 0, len(r.members))
	for _, member := range r.members {
		members = append(members, member)
	}
	return members, nil
}

func (r FamilyRepository) Save(ctx context.Context, member family.Member) error {
	r.members[member.Id()] = member
	return nil
}

func (r FamilyRepository) Delete(ctx context.Context, id uuid.UUID) error {
	delete(r.members, id)
	return nil
}

func (r FamilyRepository) Exists(ctx context.Context, members ...uuid.UUID) (bool, error) {
	for _, id := range members {
		if _, exists := r.members[id]; !exists {
			return false, nil
		}
	}
	return true, nil
}
