package persistence

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/internal/domain/entity"
	"github.com/oleexo/subtracker/internal/domain/provider"
	"github.com/oleexo/subtracker/pkg/slicesx"
)

type ProviderRepository struct {
	repository *DatabaseContext
}

func NewProviderRepository(repository *DatabaseContext) *ProviderRepository {
	return &ProviderRepository{repository: repository}
}

func (r ProviderRepository) GetById(ctx context.Context, providerId uuid.UUID) (provider.Provider, error) {
	_, ok := auth.GetUserIdFromContext(ctx)
	if !ok {
		return nil, nil
	}
	var model ProviderSqlModel
	result := r.repository.db.WithContext(ctx).
		Preload("Labels").
		Preload("Plans").
		Preload("Plans.Prices").
		// todo add filter by owner
		First(&model, providerId)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	p := newProvider(model)
	p.Clean()
	return p, nil
}

func (r ProviderRepository) GetAll(ctx context.Context, parameters entity.QueryParameters) (
	[]provider.Provider,
	error) {
	_, ok := auth.GetUserIdFromContext(ctx)
	if !ok {
		return nil, nil
	}

	var providerSqlModels []ProviderSqlModel
	query := r.repository.db.WithContext(ctx).
		Preload("Labels").
		Preload("Plans").
		Preload("Plans.Prices")
	// todo add filter by owner

	if parameters.Offset >= 0 {
		query = query.Offset(parameters.Offset)
	}
	if parameters.Limit > 0 {
		query = query.Limit(parameters.Limit)
	}
	result := query.Find(&providerSqlModels)
	if result.Error != nil {
		return nil, result.Error
	}
	providers := make([]provider.Provider, 0, len(providerSqlModels))
	for _, model := range providerSqlModels {
		p := newProvider(model)
		p.Clean()
		providers = append(providers, p)
	}
	return providers, nil
}

func (r ProviderRepository) GetAllCount(ctx context.Context) (int64, error) {
	var count int64
	result := r.repository.db.WithContext(ctx).
		Model(&ProviderSqlModel{}).
		// todo add filter by owner
		Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

func (r ProviderRepository) Save(ctx context.Context, dirtyProvider provider.Provider) error {
	if !dirtyProvider.IsDirty() {
		return nil
	}

	if dirtyProvider.IsExists() {
		return r.update(ctx, dirtyProvider)
	}
	return r.create(ctx, dirtyProvider)
}

func (r ProviderRepository) Delete(ctx context.Context, providerId uuid.UUID) (bool, error) {
	result := r.repository.db.WithContext(ctx).Delete(&ProviderSqlModel{}, providerId)
	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

func (r ProviderRepository) DeletePlan(ctx context.Context, planId uuid.UUID) (bool, error) {
	result := r.repository.db.WithContext(ctx).Delete(&providerPlanSqlModel{}, planId)
	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

func (r ProviderRepository) DeletePrice(ctx context.Context, priceId uuid.UUID) (bool, error) {
	result := r.repository.db.WithContext(ctx).Delete(&providerPriceSqlModel{}, priceId)
	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

func (r ProviderRepository) Exists(ctx context.Context, ids ...uuid.UUID) (bool, error) {
	var count int64
	result := r.repository.db.WithContext(ctx).
		Model(&ProviderSqlModel{}).
		Where("id IN ?", ids).
		Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count == int64(len(ids)), nil
}

func (r ProviderRepository) create(ctx context.Context, newProvider provider.Provider) error {
	sqlModel := newProviderSqlModel(newProvider)
	result := r.repository.db.WithContext(ctx).
		Omit("Labels", "Plans").
		Create(&sqlModel)
	if result.Error != nil {
		return result.Error
	}
	newProvider.Clean()
	if newProvider.Labels().HasChanges() {
		if err := r.createLabels(ctx, newProvider.Id(), newProvider.Labels().Values()); err != nil {
			return err
		}
		newProvider.Labels().ClearChanges()
	}

	if newProvider.Plans().HasChanges() {
		if err := r.createPlans(ctx, newProvider.Id(), newProvider.Plans().Values()); err != nil {
			return err
		}
		newProvider.Plans().ClearChanges()
	}

	return nil
}

func (r ProviderRepository) createLabels(ctx context.Context, providerId uuid.UUID, labelIds []uuid.UUID) error {
	sqlModels := slicesx.Map(labelIds, func(labelId uuid.UUID) providerLabelSqlModel {
		return newProviderLabelSqlModel(providerId, labelId)
	})
	result := r.repository.db.WithContext(ctx).Create(&sqlModels)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r ProviderRepository) createPlans(ctx context.Context, providerId uuid.UUID, plans []provider.Plan) error {
	planSqlModels := slicesx.Map(plans, func(plan provider.Plan) providerPlanSqlModel {
		return newProviderPlanSqlModel(providerId, plan)
	})

	result := r.repository.db.WithContext(ctx).
		Omit("Prices").
		Create(&planSqlModels)
	if result.Error != nil {
		return result.Error
	}

	priceSqlModels := slicesx.SelectMany(plans, func(plan provider.Plan) []providerPriceSqlModel {
		models := make([]providerPriceSqlModel, plan.Prices().Len())
		for i, price := range plan.Prices().Values() {
			models[i] = newProviderPriceSqlModel(plan.Id(), price)
		}
		return models
	})

	result = r.repository.db.WithContext(ctx).Create(&priceSqlModels)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r ProviderRepository) update(ctx context.Context, dirtyProvider provider.Provider) error {
	sqlModel := newProviderSqlModel(dirtyProvider)

	result := r.repository.db.WithContext(ctx).
		Omit("Labels", "Plans").
		Save(&sqlModel)
	if result.Error != nil {
		return result.Error
	}

	if dirtyProvider.Labels().HasChanges() {
		if err := r.updateLabels(ctx, dirtyProvider.Id(), dirtyProvider.Labels()); err != nil {
			return err
		}

		dirtyProvider.Labels().ClearChanges()
	}

	if dirtyProvider.Plans().HasChanges() {
		if err := r.updatePlans(ctx, dirtyProvider.Id(), dirtyProvider.Plans()); err != nil {
			return err
		}

		dirtyProvider.Plans().ClearChanges()
	}

	return nil
}

func (r ProviderRepository) updateLabels(
	ctx context.Context,
	providerId uuid.UUID,
	labels *slicesx.Tracked[uuid.UUID]) error {
	return saveTrackedSlice(ctx,
		r.repository.db,
		labels,
		func(labelId uuid.UUID) providerLabelSqlModel {
			return newProviderLabelSqlModel(providerId, labelId)
		},
	)
}

func (r ProviderRepository) updatePlans(
	ctx context.Context,
	providerId uuid.UUID,
	plans *slicesx.Tracked[provider.Plan]) error {
	if err := saveTrackedSlice(ctx,
		r.repository.db,
		plans,
		func(plan provider.Plan) providerPlanSqlModel {
			return newProviderPlanSqlModel(providerId, plan)
		}); err != nil {
		return err
	}

	for _, plan := range plans.Values() {
		if plan.IsDirty() && plan.Prices().HasChanges() {
			if err := r.updatePrices(ctx, plan.Id(), plan.Prices()); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r ProviderRepository) updatePrices(
	ctx context.Context,
	planId uuid.UUID,
	prices *slicesx.Tracked[provider.Price]) error {
	return saveTrackedSlice(ctx,
		r.repository.db,
		prices,
		func(price provider.Price) providerPriceSqlModel {
			return newProviderPriceSqlModel(planId, price)
		},
	)
}
