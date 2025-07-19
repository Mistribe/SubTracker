package persistence

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/oleexo/subtracker/internal/domain/label"
	"github.com/oleexo/subtracker/pkg/langext/option"
)

type labelModel struct {
	BaseModel
	Name      string `gorm:"type:varchar(100);not null"`
	Color     string `gorm:"type:varchar(20);not null"`
	IsDefault bool   `gorm:"type:boolean;not null;default:false;index"`
}

func (l labelModel) TableName() string {
	return "labels"
}

type LabelRepository struct {
	repository *Repository
}

func NewLabelRepository(repository *Repository) *LabelRepository {
	return &LabelRepository{
		repository: repository,
	}
}

func (r LabelRepository) toModel(source *label.Label) labelModel {
	return labelModel{
		BaseModel: BaseModel{
			Id:        source.Id(),
			CreatedAt: source.CreatedAt(),
			UpdatedAt: source.UpdatedAt(),
			Etag:      source.ETag(),
		},
		Name:      source.Name(),
		Color:     source.Color(),
		IsDefault: source.IsDefault(),
	}
}

func (r LabelRepository) toEntity(source labelModel) label.Label {
	lbl := label.NewLabelWithoutValidation(
		source.Id,
		source.Name,
		source.IsDefault,
		source.Color,
		source.CreatedAt,
		source.UpdatedAt,
		true,
	)
	lbl.Clean()
	return lbl
}

func (r LabelRepository) Get(ctx context.Context, id uuid.UUID) (option.Option[label.Label], error) {
	var model labelModel
	if err := r.repository.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return option.None[label.Label](), nil
		}
		return option.None[label.Label](), err
	}
	return option.Some(r.toEntity(model)), nil
}

func (r LabelRepository) GetAll(ctx context.Context, withDefault bool) ([]label.Label, error) {
	var labels []labelModel
	query := r.repository.db.WithContext(ctx)
	if !withDefault {
		query = query.Where("is_default = ?", false)
	}
	if result := query.Find(&labels); result.Error != nil {
		return nil, result.Error
	}
	result := make([]label.Label, 0, len(labels))
	for _, lbl := range labels {
		result = append(result, r.toEntity(lbl))
	}
	return result, nil
}

func (r LabelRepository) GetDefaults(ctx context.Context) ([]label.Label, error) {
	var labels []labelModel
	if result := r.repository.db.WithContext(ctx).Where("is_default = ?", true).Find(&labels); result.Error != nil {
		return nil, result.Error
	}
	result := make([]label.Label, 0, len(labels))
	for _, lbl := range labels {
		result = append(result, r.toEntity(lbl))
	}
	return result, nil
}

func (r LabelRepository) Save(ctx context.Context, lbl *label.Label) error {
	if lbl.IsDirty() == false {
		return nil
	}
	if lbl.IsExists() {
		return r.update(ctx, lbl)
	}
	return r.create(ctx, lbl)
}

func (r LabelRepository) create(ctx context.Context, lbl *label.Label) error {
	dbLabel := r.toModel(lbl)
	result := r.repository.db.WithContext(ctx).Create(&dbLabel)
	if result.Error != nil {
		return result.Error
	}
	lbl.Clean()
	return nil
}
func (r LabelRepository) update(ctx context.Context, lbl *label.Label) error {
	dbLabel := r.toModel(lbl)
	result := r.repository.db.WithContext(ctx).Save(dbLabel)
	if result.Error != nil {
		return result.Error
	}
	lbl.Clean()
	return nil
}

func (r LabelRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if result := r.repository.db.WithContext(ctx).Delete(&labelModel{}, id); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r LabelRepository) Exists(ctx context.Context, ids ...uuid.UUID) (bool, error) {
	var count int64
	if len(ids) == 1 {
		result := r.repository.db.WithContext(ctx).Model(&labelModel{}).Where("id = ?", ids[0]).Count(&count)
		if result.Error != nil {
			return false, result.Error
		}
	} else {
		result := r.repository.db.WithContext(ctx).Model(&labelModel{}).Where("id IN ?", ids).Count(&count)
		if result.Error != nil {
			return false, result.Error
		}
	}
	return count == int64(len(ids)), nil
}
