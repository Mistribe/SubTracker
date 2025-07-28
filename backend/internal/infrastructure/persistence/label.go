package persistence

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/oleexo/subtracker/internal/domain/label"
	"github.com/oleexo/subtracker/internal/domain/user"
	"github.com/oleexo/subtracker/pkg/langext/option"
)

type LabelRepository struct {
	repository *DatabaseContext
}

func NewLabelRepository(repository *DatabaseContext) *LabelRepository {
	return &LabelRepository{
		repository: repository,
	}
}

func (r LabelRepository) GetById(ctx context.Context, labelId uuid.UUID) (option.Option[label.Label], error) {
	var model labelSqlModel
	if err := r.repository.db.WithContext(ctx).First(&model, labelId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return option.None[label.Label](), nil
		}
		return option.None[label.Label](), err
	}
	lbl := newLabel(model)
	lbl.Clean()
	return option.Some(lbl), nil
}

func (r LabelRepository) GetAll(ctx context.Context, parameters label.QueryParameters) ([]label.Label, error) {
	userId, ok := user.FromContext(ctx)
	if !ok {
		return nil, nil
	}
	var labels []labelSqlModel
	query := r.repository.db.WithContext(ctx)

	if parameters.WithDefaults {
		query = query.Where("owner_id = ? OR (owner_id is null AND is_default = ?)", userId, true)
	} else {
		query = query.Where("owner_id = ?", userId)
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
	userId, ok := user.FromContext(ctx)
	if !ok {
		return 0, nil
	}
	var labels []labelSqlModel
	query := r.repository.db.WithContext(ctx).Model(&labelSqlModel{})

	if parameters.WithDefaults {
		query = query.Where("owner_id = ? OR (owner_id is null AND is_default = ?)", userId, true)
	} else {
		query = query.Where("owner_id = ?", userId)
	}
	var count int64
	if result := query.Count(&count); result.Error != nil {
		return 0, result.Error
	}
	result := make([]label.Label, 0, len(labels))
	for _, lbl := range labels {
		result = append(result, newLabel(lbl))
	}
	return count, nil
}

func (r LabelRepository) GetDefaults(ctx context.Context) ([]label.Label, error) {
	var labels []labelSqlModel
	if result := r.repository.db.WithContext(ctx).Where("is_default = ?", true).Find(&labels); result.Error != nil {
		return nil, result.Error
	}
	result := make([]label.Label, 0, len(labels))
	for _, model := range labels {
		lbl := newLabel(model)
		lbl.Clean()
		result = append(result)
	}
	return result, nil
}

func (r LabelRepository) Save(ctx context.Context, lbl label.Label) error {
	if lbl.IsDirty() == false {
		return nil
	}
	if lbl.IsExists() {
		return r.update(ctx, lbl)
	}
	return r.create(ctx, lbl)
}

func (r LabelRepository) create(ctx context.Context, lbl label.Label) error {
	dbLabel := newLabelSqlModel(lbl)
	result := r.repository.db.WithContext(ctx).Create(&dbLabel)
	if result.Error != nil {
		return result.Error
	}
	lbl.Clean()
	return nil
}
func (r LabelRepository) update(ctx context.Context, lbl label.Label) error {
	dbLabel := newLabelSqlModel(lbl)
	result := r.repository.db.WithContext(ctx).Save(dbLabel)
	if result.Error != nil {
		return result.Error
	}
	lbl.Clean()
	return nil
}

func (r LabelRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if result := r.repository.db.WithContext(ctx).Delete(&labelSqlModel{}, id); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r LabelRepository) Exists(ctx context.Context, ids ...uuid.UUID) (bool, error) {
	var count int64
	if len(ids) == 1 {
		result := r.repository.db.WithContext(ctx).Model(&labelSqlModel{}).Where("id = ?", ids[0]).Count(&count)
		if result.Error != nil {
			return false, result.Error
		}
	} else {
		result := r.repository.db.WithContext(ctx).Model(&labelSqlModel{}).Where("id IN ?", ids).Count(&count)
		if result.Error != nil {
			return false, result.Error
		}
	}
	return count == int64(len(ids)), nil
}
