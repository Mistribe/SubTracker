package persistence

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/internal/domain/family"
	"github.com/oleexo/subtracker/internal/domain/label"
)

type LabelRepository struct {
	labelRepository *DatabaseContext
	authService     auth.Service
}

func NewLabelRepository(labelRepository *DatabaseContext,
	authService auth.Service) label.Repository {
	return &LabelRepository{
		labelRepository: labelRepository,
		authService:     authService,
	}
}

func (r LabelRepository) GetById(ctx context.Context, labelId uuid.UUID) (label.Label, error) {
	var model LabelSqlModel
	if err := r.labelRepository.db.WithContext(ctx).First(&model, labelId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	lbl := newLabel(model)
	lbl.Clean()
	return lbl, nil
}

func (r LabelRepository) GetAll(ctx context.Context, parameters label.QueryParameters) ([]label.Label, error) {
	userId, ok := auth.GetUserIdFromContext(ctx)
	if !ok {
		return nil, nil
	}
	var labels []LabelSqlModel
	query := r.labelRepository.db.WithContext(ctx)

	for i, ownerType := range parameters.Owners {
		switch ownerType {
		case auth.PersonalOwnerType:
			if i == 0 {
				query = query.Where("owner_type = ? AND owner_user_id = ?", auth.PersonalOwnerType, userId)
			} else {
				query = query.Or("owner_type = ? AND owner_user_id = ?", auth.PersonalOwnerType, userId)
			}
		case auth.SystemOwnerType:
			if i == 0 {
				query = query.Where("owner_type = ?", auth.SystemOwnerType)
			} else {
				query = query.Or("owner_type = ?", auth.SystemOwnerType)
			}
		case auth.FamilyOwnerType:
			var families []uuid.UUID
			if parameters.FamilyId != nil {
				if r.authService.IsInFamily(ctx, *parameters.FamilyId) {
					return nil, family.ErrFamilyNotFound
				}
				families = append(families, *parameters.FamilyId)
			} else {
				families = r.authService.MustGetFamilies(ctx)
			}
			if len(families) > 0 {
				if i == 0 {
					query = query.Where("owner_type = ? AND owner_family_id IN ?", auth.FamilyOwnerType, families)
				} else {
					query = query.Or("owner_type = ? AND owner_family_id IN ?", auth.FamilyOwnerType, families)
				}
			}
		}
	}

	if parameters.Offset > 0 {
		query = query.Offset(parameters.Offset)
	}
	if parameters.Limit > 0 {
		query = query.Limit(parameters.Limit)
	}
	if result := query.Find(&labels); result.Error != nil {
		return nil, result.Error
	}
	result := make([]label.Label, 0, len(labels))
	for _, model := range labels {
		lbl := newLabel(model)
		lbl.Clean()
		result = append(result, lbl)
	}
	return result, nil
}

func (r LabelRepository) GetAllCount(ctx context.Context, parameters label.QueryParameters) (int64, error) {
	userId, ok := auth.GetUserIdFromContext(ctx)
	if !ok {
		return 0, nil
	}
	query := r.labelRepository.db.WithContext(ctx).Model(&LabelSqlModel{})

	for i, ownerType := range parameters.Owners {
		switch ownerType {
		case auth.PersonalOwnerType:
			if i == 0 {
				query = query.Where("owner_type = ? AND owner_user_id = ?", auth.PersonalOwnerType, userId)
			} else {
				query = query.Or("owner_type = ? AND owner_user_id = ?", auth.PersonalOwnerType, userId)
			}
		case auth.SystemOwnerType:
			if i == 0 {
				query = query.Where("owner_type = ?", auth.SystemOwnerType)
			} else {
				query = query.Or("owner_type = ?", auth.SystemOwnerType)
			}
		case auth.FamilyOwnerType:
			families := r.authService.MustGetFamilies(ctx)
			if len(families) > 0 {
				if i == 0 {
					query = query.Where("owner_type = ? AND owner_family_id IN ?", auth.FamilyOwnerType, families)
				} else {
					query = query.Or("owner_type = ? AND owner_family_id IN ?", auth.FamilyOwnerType, families)
				}
			}
		}
	}
	var count int64
	if result := query.Count(&count); result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

func (r LabelRepository) GetSystemLabels(ctx context.Context) ([]label.Label, error) {
	var labelModels []LabelSqlModel
	query := r.labelRepository.db.WithContext(ctx).Where("owner_type = ?", auth.SystemOwnerType)
	result := query.Find(&labelModels)
	if result.Error != nil {
		return nil, result.Error
	}
	labels := make([]label.Label, 0, len(labelModels))
	for _, model := range labelModels {
		lbl := newLabel(model)
		lbl.Clean()
		labels = append(labels, lbl)
	}
	return labels, nil
}

func (r LabelRepository) Save(ctx context.Context, lbl label.Label) error {
	if !lbl.IsDirty() {
		return nil
	}
	if lbl.IsExists() {
		return r.update(ctx, lbl)
	}
	return r.create(ctx, lbl)
}

func (r LabelRepository) create(ctx context.Context, lbl label.Label) error {
	dbLabel := newLabelSqlModel(lbl)
	result := r.labelRepository.db.WithContext(ctx).Create(&dbLabel)
	if result.Error != nil {
		return result.Error
	}
	lbl.Clean()
	return nil
}
func (r LabelRepository) update(ctx context.Context, lbl label.Label) error {
	dbLabel := newLabelSqlModel(lbl)
	result := r.labelRepository.db.WithContext(ctx).Save(dbLabel)
	if result.Error != nil {
		return result.Error
	}
	lbl.Clean()
	return nil
}

func (r LabelRepository) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	result := r.labelRepository.db.WithContext(ctx).Delete(&LabelSqlModel{}, id)
	if result.Error != nil {
		return false, result.Error
	}

	if result.RowsAffected == 0 {
		return false, nil
	}

	return true, nil
}

func (r LabelRepository) Exists(ctx context.Context, ids ...uuid.UUID) (bool, error) {
	var count int64
	if len(ids) == 1 {
		result := r.labelRepository.db.WithContext(ctx).Model(&LabelSqlModel{}).Where("id = ?", ids[0]).Count(&count)
		if result.Error != nil {
			return false, result.Error
		}
	} else {
		result := r.labelRepository.db.WithContext(ctx).Model(&LabelSqlModel{}).Where("id IN ?", ids).Count(&count)
		if result.Error != nil {
			return false, result.Error
		}
	}
	return count == int64(len(ids)), nil
}
