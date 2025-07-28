package persistence

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/oleexo/subtracker/internal/domain/user"

	"github.com/oleexo/subtracker/internal/domain/family"
	"github.com/oleexo/subtracker/pkg/langext/option"
)

type FamilyRepository struct {
	repository *DatabaseContext
}

func NewFamilyRepository(repository *DatabaseContext) *FamilyRepository {
	return &FamilyRepository{
		repository: repository,
	}
}

func (r FamilyRepository) GetById(ctx context.Context, id uuid.UUID) (option.Option[family.Family], error) {
	var model familySqlModel
	result := r.repository.db.WithContext(ctx).
		Preload("Members").
		First(&model, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return option.None[family.Family](), nil
		}
		return option.None[family.Family](), result.Error
	}
	fam := newFamily(model)
	fam.Clean()
	return option.Some(fam), nil
}

func (r FamilyRepository) GetOwn(ctx context.Context) (option.Option[family.Family], error) {
	userId, ok := user.FromContext(ctx)
	if !ok {
		return option.None[family.Family](), nil
	}
	var model familySqlModel
	if err := r.repository.db.WithContext(ctx).Where("owner_id = ?", userId).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return option.None[family.Family](), nil
		}
		return option.None[family.Family](), err
	}

	fam := newFamily(model)
	return option.Some(fam), nil
}

func (r FamilyRepository) GetAll(ctx context.Context, size, page int) ([]family.Family, error) {
	userId, ok := user.FromContext(ctx)
	if !ok {
		return []family.Family{}, nil
	}
	var familyModels []familySqlModel
	result := r.repository.db.WithContext(ctx).
		Preload("Members").
		Where("owner_id = ?", userId).
		Offset((page - 1) * size).Limit(size).
		Find(&familyModels)
	if result.Error != nil {
		return nil, result.Error
	}
	families := make([]family.Family, 0, len(familyModels))
	for _, model := range familyModels {
		fam := newFamily(model)
		families = append(families, fam)
	}
	return families, nil
}

func (r FamilyRepository) GetAllCount(ctx context.Context) (int64, error) {
	userId, ok := user.FromContext(ctx)
	if !ok {
		return 0, nil
	}
	var count int64
	result := r.repository.db.WithContext(ctx).
		Model(&familySqlModel{}).
		Where("owner_id = ?", userId).
		Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

func (r FamilyRepository) GetAllMembers(ctx context.Context, familyId uuid.UUID) ([]family.Member, error) {
	var members []familyMemberSqlModel
	if result := r.repository.db.WithContext(ctx).Where("family_id = ?", familyId).Find(&members); result.Error != nil {
		return nil, result.Error
	}
	result := make([]family.Member, 0, len(members))
	for _, model := range members {
		mbr := newFamilyMember(model)
		result = append(result, mbr)
	}
	return result, nil
}

func (r FamilyRepository) Save(ctx context.Context, dirtyFamily family.Family) error {
	if !dirtyFamily.IsDirty() {
		return nil
	}
	dbMember := newFamilySqlModel(dirtyFamily)
	var result *gorm.DB
	if dirtyFamily.IsExists() {
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

	if err := saveTrackedSlice(ctx,
		r.repository.db,
		dirtyFamily.Members(), func(mbr family.Member) familyMemberSqlModel {
			return newFamilyMemberSqlModel(mbr)
		}); err != nil {
		return err
	}
	dirtyFamily.Clean()
	return nil
}

func (r FamilyRepository) DeleteMember(ctx context.Context, memberId uuid.UUID) error {
	return r.repository.db.WithContext(ctx).Delete(&familyMemberSqlModel{}, memberId).Error
}

func (r FamilyRepository) Delete(ctx context.Context, familyId uuid.UUID) error {
	return r.repository.db.WithContext(ctx).Delete(&familySqlModel{}, familyId).Error
}

func (r FamilyRepository) MemberExists(ctx context.Context, familyId uuid.UUID, members ...uuid.UUID) (bool, error) {
	var count int64
	if len(members) == 1 {
		result := r.repository.db.WithContext(ctx).Model(&familyMemberSqlModel{}).Where("family_id = ? AND id = ?",
			familyId, members[0]).Count(&count)
		if result.Error != nil {
			return false, result.Error
		}
	} else {
		result := r.repository.db.WithContext(ctx).Model(&familyMemberSqlModel{}).Where("family_id = ? AND id IN ?",
			familyId, members).Count(&count)
		if result.Error != nil {
			return false, result.Error
		}
	}
	return count == int64(len(members)), nil
}
