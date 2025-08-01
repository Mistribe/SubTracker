package persistence

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/internal/domain/entity"
	"github.com/oleexo/subtracker/internal/domain/family"
)

type FamilyRepository struct {
	repository *DatabaseContext
}

func NewFamilyRepository(repository *DatabaseContext) family.Repository {
	return &FamilyRepository{
		repository: repository,
	}
}

func (r FamilyRepository) GetById(ctx context.Context, id uuid.UUID) (family.Family, error) {
	var model FamilySqlModel
	result := r.repository.db.WithContext(ctx).
		Preload("Members").
		First(&model, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	fam := newFamily(model)
	fam.Clean()
	return fam, nil
}

func (r FamilyRepository) GetAll(ctx context.Context, parameters entity.QueryParameters) ([]family.Family, error) {
	userId, ok := auth.GetUserIdFromContext(ctx)
	if !ok {
		return []family.Family{}, nil
	}
	var familyModels []FamilySqlModel
	query := r.repository.db.WithContext(ctx).
		Preload("Members").
		Where("owner_id = ?", userId)
	if parameters.Offset > 0 {
		query = query.Offset(parameters.Offset)
	}
	if parameters.Limit > 0 {
		query = query.Limit(parameters.Limit)
	}

	result := query.Find(&familyModels)
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
	userId, ok := auth.GetUserIdFromContext(ctx)
	if !ok {
		return 0, nil
	}
	var count int64
	result := r.repository.db.WithContext(ctx).
		Model(&FamilySqlModel{}).
		Where("owner_id = ?", userId).
		Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

func (r FamilyRepository) GetAllMembers(ctx context.Context, familyId uuid.UUID) ([]family.Member, error) {
	var members []FamilyMemberSqlModel
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
		dirtyFamily.Members(), func(mbr family.Member) FamilyMemberSqlModel {
			return newFamilyMemberSqlModel(mbr)
		}); err != nil {
		return err
	}
	dirtyFamily.Clean()
	return nil
}

func (r FamilyRepository) DeleteMember(ctx context.Context, memberId uuid.UUID) (bool, error) {
	result := r.repository.db.WithContext(ctx).Delete(&FamilyMemberSqlModel{}, memberId)
	if result.Error != nil {
		return false, result.Error
	}

	if result.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

func (r FamilyRepository) Delete(ctx context.Context, familyId uuid.UUID) (bool, error) {
	result := r.repository.db.WithContext(ctx).Delete(&FamilySqlModel{}, familyId)
	if result.Error != nil {
		return false, result.Error
	}

	if result.RowsAffected == 0 {
		return false, nil
	}

	return true, nil
}

func (r FamilyRepository) MemberExists(ctx context.Context, familyId uuid.UUID, members ...uuid.UUID) (bool, error) {
	var count int64
	if len(members) == 1 {
		result := r.repository.db.WithContext(ctx).Model(&FamilyMemberSqlModel{}).Where("family_id = ? AND id = ?",
			familyId, members[0]).Count(&count)
		if result.Error != nil {
			return false, result.Error
		}
	} else {
		result := r.repository.db.WithContext(ctx).Model(&FamilyMemberSqlModel{}).Where("family_id = ? AND id IN ?",
			familyId, members).Count(&count)
		if result.Error != nil {
			return false, result.Error
		}
	}
	return count == int64(len(members)), nil
}

func (r FamilyRepository) Exists(ctx context.Context, ids ...uuid.UUID) (bool, error) {
	var count int64
	if len(ids) == 1 {
		result := r.repository.db.WithContext(ctx).Model(&FamilySqlModel{}).Where("id = ?", ids[0]).Count(&count)
		if result.Error != nil {
			return false, result.Error
		}
	} else {
		result := r.repository.db.WithContext(ctx).Model(&FamilySqlModel{}).Where("id IN ?", ids).Count(&count)
		if result.Error != nil {
			return false, result.Error
		}
	}

	return count == int64(len(ids)), nil
}

func (r FamilyRepository) IsUserMemberOfFamily(ctx context.Context, familyId uuid.UUID, userId string) (bool, error) {
	var count int64
	result := r.repository.db.WithContext(ctx).
		Model(&FamilyMemberSqlModel{}).
		Where("family_id = ? AND user_id = ?", familyId, userId).
		Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}
