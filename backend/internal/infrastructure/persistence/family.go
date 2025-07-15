package persistence

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/oleexo/subtracker/internal/domain/family"
	"github.com/oleexo/subtracker/pkg/langext/option"
)

type familyMemberModel struct {
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

func (r FamilyRepository) toModel(source *family.Member) familyMemberModel {

	return familyMemberModel{
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

func (r FamilyRepository) toEntity(source familyMemberModel) family.Member {
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
	var model familyMemberModel
	if err := r.repository.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return option.None[family.Member](), nil
		}
		return option.None[family.Member](), err
	}
	return option.Some(r.toEntity(model)), nil
}

func (r FamilyRepository) GetAll(ctx context.Context) ([]family.Member, error) {
	var members []familyMemberModel
	if result := r.repository.db.WithContext(ctx).Find(&members); result.Error != nil {
		return nil, result.Error
	}
	result := make([]family.Member, 0, len(members))
	for _, mbr := range members {
		result = append(result, r.toEntity(mbr))
	}
	return result, nil
}

func (r FamilyRepository) Save(ctx context.Context, member family.Member) error {
	if !member.IsDirty() {
		return nil
	}
	dbMember := r.toModel(&member)
	var result *gorm.DB
	if member.IsExists() {
		result = r.repository.db.WithContext(ctx).Save(&dbMember)
	} else {
		result = r.repository.db.WithContext(ctx).Create(&dbMember)
	}

	if result.Error != nil {
		return result.Error
	}
	member.Clean()
	return nil
}

func (r FamilyRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.repository.db.WithContext(ctx).Delete(&familyMemberModel{}, id).Error
}

func (r FamilyRepository) Exists(ctx context.Context, members ...uuid.UUID) (bool, error) {
	var count int64
	if len(members) == 1 {
		result := r.repository.db.WithContext(ctx).Model(&familyMemberModel{}).Where("id = ?", members[0]).Count(&count)
		if result.Error != nil {
			return false, result.Error
		}
	} else {
		result := r.repository.db.WithContext(ctx).Model(&familyMemberModel{}).Where("id IN ?", members).Count(&count)
		if result.Error != nil {
			return false, result.Error
		}
	}
	return count == int64(len(members)), nil
}
