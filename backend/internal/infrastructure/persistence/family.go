package persistence

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/family"
	"github.com/oleexo/subtracker/internal/domain/label"
	"github.com/oleexo/subtracker/pkg/langext/option"
)

type memberModel struct {
	BaseModel
	Name  string         `gorm:"type:varchar(100);not null"`
	Email sql.NullString `gorm:"type:varchar(100)"`
	IsKid bool           `gorm:"type:boolean;not null;default:false"`
}

type FamilyRepository struct {
	repository *Repository
}

func NewFamilyRepository(repository *Repository) *FamilyRepository {
	return &FamilyRepository{
		repository: repository,
	}
}

func (r FamilyRepository) toModel(source *family.Member) memberModel {

	return memberModel{
		BaseModel: BaseModel{
			Id:        source.Id(),
			CreatedAt: source.CreatedAt(),
			UpdatedAt: source.UpdatedAt(),
		},
		Name: source.Name(),
		Email: option.Match(source.Email(), func(in string) sql.NullString {
			return sql.NullString{
				String: in,
				Valid:  true,
			}
		}, func() sql.NullString {
			return sql.NullString{
				Valid: false,
			}
		}),
		IsKid: source.IsKid(),
	}
}

func (r FamilyRepository) toEntity(source memberModel) family.Member {
	var email *string
	if source.Email.Valid {
		email = &source.Email.String
	} else {
		email = nil
	}
	lbl := family.NewMemberWithoutValidation(
		source.Id,
		source.Name,
		email,
		source.IsKid,
		source.CreatedAt,
		source.UpdatedAt,
		true,
	)
	lbl.Clean()
	return lbl
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
