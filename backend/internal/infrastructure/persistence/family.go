package persistence

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/oleexo/subtracker/internal/domain/user"

	"github.com/oleexo/subtracker/internal/domain/family"
	"github.com/oleexo/subtracker/pkg/langext/option"
)

type familyModel struct {
	BaseModel
	Name             string              `gorm:"type:varchar(100);not null"`
	OwnerId          string              `gorm:"type:varchar(100);not null"`
	HaveJointAccount bool                `gorm:"type:boolean;not null;default:false"`
	Members          []familyMemberModel `gorm:"foreignKey:FamilyId;references:Id"`
}

func (f familyModel) TableName() string {
	return "families"
}

type familyMemberModel struct {
	BaseModel
	Name     string         `gorm:"type:varchar(100);not null"`
	FamilyId uuid.UUID      `gorm:"type:uuid;not null"`
	Family   familyModel    `gorm:"foreignKey:FamilyId;references:Id"`
	Email    sql.NullString `gorm:"type:varchar(100)"`
	UserId   sql.NullString `gorm:"type:varchar(100)"`
	IsKid    bool           `gorm:"type:boolean;not null;default:false"`
}

func (f familyMemberModel) TableName() string {
	return "family_members"
}

type FamilyRepository struct {
	repository *Repository
}

func NewFamilyRepository(repository *Repository) *FamilyRepository {
	return &FamilyRepository{
		repository: repository,
	}
}

func (r FamilyRepository) toFamilyMemberModel(source family.Member) familyMemberModel {
	return familyMemberModel{
		BaseModel: BaseModel{
			Id:        source.Id(),
			CreatedAt: source.CreatedAt(),
			UpdatedAt: source.UpdatedAt(),
			Etag:      source.ETag(),
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
		IsKid:    source.IsKid(),
		FamilyId: source.FamilyId(),
	}
}

func (r FamilyRepository) toFamilyModel(source *family.Family) familyModel {
	return familyModel{
		BaseModel: BaseModel{
			Id:        source.Id(),
			CreatedAt: source.CreatedAt(),
			UpdatedAt: source.UpdatedAt(),
			Etag:      source.ETag(),
		},
		Name:             source.Name(),
		OwnerId:          source.OwnerId(),
		HaveJointAccount: source.HaveJointAccount(),
	}
}

func (r FamilyRepository) toMemberEntity(source familyMemberModel) family.Member {
	var email *string
	if source.Email.Valid {
		email = &source.Email.String
	} else {
		email = nil
	}
	mbr := family.NewMemberWithoutValidation(
		source.Id,
		source.FamilyId,
		source.Name,
		email,
		source.IsKid,
		source.CreatedAt,
		source.UpdatedAt,
		true,
	)
	mbr.Clean()
	return mbr
}

func (r FamilyRepository) toFamilyEntity(source familyModel) family.Family {
	members := make([]family.Member, 0, len(source.Members))
	for _, member := range source.Members {
		members = append(members, r.toMemberEntity(member))
	}
	fam := family.NewFamilyWithoutValidation(
		source.Id,
		source.OwnerId,
		source.Name,
		source.HaveJointAccount,
		members,
		source.CreatedAt,
		source.UpdatedAt,
		true,
	)

	fam.Clean()
	return fam
}

func (r FamilyRepository) GetById(ctx context.Context, id uuid.UUID) (option.Option[family.Family], error) {
	var model familyModel
	result := r.repository.db.WithContext(ctx).
		Preload("Members").
		First(&model, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return option.None[family.Family](), nil
		}
		return option.None[family.Family](), result.Error
	}
	return option.Some(r.toFamilyEntity(model)), nil
}

func (r FamilyRepository) GetOwn(ctx context.Context) (option.Option[family.Family], error) {
	userId, ok := user.FromContext(ctx)
	if !ok {
		return option.None[family.Family](), nil
	}
	var model familyModel
	if err := r.repository.db.WithContext(ctx).Where("owner_id = ?", userId).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return option.None[family.Family](), nil
		}
		return option.None[family.Family](), err
	}
	return option.Some(r.toFamilyEntity(model)), nil
}

func (r FamilyRepository) GetAll(ctx context.Context) ([]family.Family, error) {
	userId, ok := user.FromContext(ctx)
	if !ok {
		return []family.Family{}, nil
	}
	var familyModels []familyModel
	result := r.repository.db.WithContext(ctx).
		Preload("Members").
		Where("owner_id = ?", userId).
		Find(&familyModels)
	if result.Error != nil {
		return nil, result.Error
	}
	families := make([]family.Family, 0, len(familyModels))
	for _, fam := range familyModels {
		families = append(families, r.toFamilyEntity(fam))
	}
	return families, nil
}

func (r FamilyRepository) GetAllMembers(ctx context.Context, familyId uuid.UUID) ([]family.Member, error) {
	var members []familyMemberModel
	if result := r.repository.db.WithContext(ctx).Where("family_id = ?", familyId).Find(&members); result.Error != nil {
		return nil, result.Error
	}
	result := make([]family.Member, 0, len(members))
	for _, mbr := range members {
		result = append(result, r.toMemberEntity(mbr))
	}
	return result, nil
}

func (r FamilyRepository) Save(ctx context.Context, family *family.Family) error {
	if !family.IsDirty() {
		return nil
	}
	dbMember := r.toFamilyModel(family)
	var result *gorm.DB
	if family.IsExists() {
		result = r.repository.db.WithContext(ctx).
			Omit("Members").
			Save(&dbMember)
	} else {
		result = r.repository.db.WithContext(ctx).
			Omit("Members").
			Create(&dbMember)
	}

	if result.Error != nil {
		return result.Error
	}

	if err := saveTrackedSlice(ctx, family.Members(), r.createMember, r.updateMember, r.deleteMember); err != nil {
		return err
	}
	family.Clean()
	return nil
}

func (r FamilyRepository) createMember(ctx context.Context, member family.Member) error {
	if !member.IsDirty() {
		return nil
	}
	dbMember := r.toFamilyMemberModel(member)
	if member.IsExists() {
		return family.ErrMemberAlreadyExists
	}
	result := r.repository.db.WithContext(ctx).Save(&dbMember)

	if result.Error != nil {
		return result.Error
	}
	member.Clean()
	return nil
}

func (r FamilyRepository) updateMember(ctx context.Context, member family.Member) error {
	if !member.IsDirty() {
		return nil
	}
	dbMember := r.toFamilyMemberModel(member)
	if !member.IsExists() {
		return family.ErrMemberNotAlreadyExists
	}
	result := r.repository.db.WithContext(ctx).Create(&dbMember)

	if result.Error != nil {
		return result.Error
	}
	member.Clean()
	return nil
}

func (r FamilyRepository) deleteMember(ctx context.Context, member family.Member) error {
	result := r.repository.db.WithContext(ctx).Delete(&familyMemberModel{}, member.Id())
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r FamilyRepository) DeleteMember(ctx context.Context, memberId uuid.UUID) error {
	return r.repository.db.WithContext(ctx).Delete(&familyMemberModel{}, memberId).Error
}

func (r FamilyRepository) Delete(ctx context.Context, familyId uuid.UUID) error {
	return r.repository.db.WithContext(ctx).Delete(&familyModel{}, familyId).Error
}

func (r FamilyRepository) MemberExists(ctx context.Context, familyId uuid.UUID, members ...uuid.UUID) (bool, error) {
	var count int64
	if len(members) == 1 {
		result := r.repository.db.WithContext(ctx).Model(&familyMemberModel{}).Where("family_id = ? AND id = ?",
			familyId, members[0]).Count(&count)
		if result.Error != nil {
			return false, result.Error
		}
	} else {
		result := r.repository.db.WithContext(ctx).Model(&familyMemberModel{}).Where("family_id = ? AND id IN ?",
			familyId, members).Count(&count)
		if result.Error != nil {
			return false, result.Error
		}
	}
	return count == int64(len(members)), nil
}

func (r FamilyRepository) MarkAsUpdated(ctx context.Context, familyId uuid.UUID, updatedAt time.Time) error {
	result := r.repository.db.WithContext(ctx).Model(&familyModel{}).Where("id = ?", familyId).Update("updated_at",
		updatedAt)
	return result.Error
}
