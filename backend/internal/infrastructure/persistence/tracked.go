package persistence

import (
	"context"

	"gorm.io/gorm"

	"github.com/oleexo/subtracker/pkg/slicesx"
)

func saveTrackedSlice[TEntity comparable, TSqlModel any](
	ctx context.Context,
	db *gorm.DB,
	s *slicesx.Tracked[TEntity],
	mapFunc func(TEntity) TSqlModel,
	omit ...string,
) error {
	var model TSqlModel

	var addedSqlModels []TSqlModel
	for add := range s.Added() {
		addedSqlModels = append(addedSqlModels, mapFunc(add))
	}
	baseRequest := db.WithContext(ctx).Model(&model)
	if len(omit) > 0 {
		baseRequest = baseRequest.Omit(omit...)
	}
	if err := baseRequest.Create(&addedSqlModels).Error; err != nil {
		return err
	}
	for remove := range s.Removed() {
		model = mapFunc(remove)
		if err := baseRequest.Delete(&model).Error; err != nil {
			return err
		}
	}

	for update := range s.Updated() {
		model = mapFunc(update)
		if err := baseRequest.Save(&model).Error; err != nil {
			return err
		}
	}

	return nil
}
