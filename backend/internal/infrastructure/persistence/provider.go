package persistence

import (
	"context"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/internal/domain/entity"
	"github.com/oleexo/subtracker/internal/domain/provider"
	"github.com/oleexo/subtracker/internal/infrastructure/persistence/sql"
	"github.com/oleexo/subtracker/pkg/slicesx"
)

type ProviderRepository struct {
	dbContext *DatabaseContext
}

func NewProviderRepository(repository *DatabaseContext) provider.Repository {
	return &ProviderRepository{dbContext: repository}
}

func (r ProviderRepository) GetById(ctx context.Context, providerId uuid.UUID) (provider.Provider, error) {
	_, ok := auth.GetUserIdFromContext(ctx)
	if !ok {
		return nil, nil
	}
	response, err := r.dbContext.GetQueries(ctx).GetProviderById(ctx, providerId)
	if err != nil {
		return nil, err
	}
	if len(response) == 0 {
		return nil, nil
	}
	providers := createProviderFromSqlcRows(response,
		func(row sql.GetProviderByIdRow) sql.Provider {
			return row.Provider
		},
		func(row sql.GetProviderByIdRow) sql.ProviderPlan {
			return row.ProviderPlan
		},
		func(row sql.GetProviderByIdRow) sql.ProviderPrice {
			return row.ProviderPrice
		},
	)

	if len(providers) == 0 {
		return nil, nil
	}

	return providers[0], nil
}

func (r ProviderRepository) GetSystemProviders(ctx context.Context) ([]provider.Provider, error) {
	response, err := r.dbContext.GetQueries(ctx).GetSystemProviders(ctx)
	if err != nil {
		return nil, err
	}
	if len(response) == 0 {
		return nil, nil
	}
	providers := createProviderFromSqlcRows(response,
		func(row sql.GetSystemProvidersRow) sql.Provider {
			return row.Provider
		},
		func(row sql.GetSystemProvidersRow) sql.ProviderPlan {
			return row.ProviderPlan
		},
		func(row sql.GetSystemProvidersRow) sql.ProviderPrice {
			return row.ProviderPrice
		},
	)

	if len(providers) == 0 {
		return nil, nil
	}

	return providers, nil
}

func (r ProviderRepository) GetAll(ctx context.Context, parameters entity.QueryParameters) (
	[]provider.Provider,
	int64,
	error) {
	_, ok := auth.GetUserIdFromContext(ctx)
	if !ok {
		return nil, 0, nil
	}

	response, err := r.dbContext.GetQueries(ctx).GetProviders(ctx, sql.GetProvidersParams{
		Limit:  parameters.Limit,
		Offset: parameters.Offset,
	})
	if err != nil {
		return nil, 0, err
	}
	if len(response) == 0 {
		return nil, 0, nil
	}
	providers := createProviderFromSqlcRows(response,
		func(row sql.GetProvidersRow) sql.Provider {
			return row.Provider
		},
		func(row sql.GetProvidersRow) sql.ProviderPlan {
			return row.ProviderPlan
		},
		func(row sql.GetProvidersRow) sql.ProviderPrice {
			return row.ProviderPrice
		},
	)

	if len(providers) == 0 {
		return nil, 0, nil
	}

	return providers, response[0].TotalCount, nil
}

func (r ProviderRepository) Save(ctx context.Context, providers ...provider.Provider) error {
	if !dirtyProvider.IsDirty() {
		return nil
	}

	if dirtyProvider.IsExists() {
		return r.update(ctx, dirtyProvider)
	}
	return r.create(ctx, dirtyProvider)
}

func (r ProviderRepository) Delete(ctx context.Context, providerId uuid.UUID) (bool, error) {
	result := r.dbContext.db.WithContext(ctx).Delete(&ProviderSqlModel{}, providerId)
	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

func (r ProviderRepository) DeletePlan(ctx context.Context, planId uuid.UUID) (bool, error) {
	result := r.dbContext.db.WithContext(ctx).Delete(&providerPlanSqlModel{}, planId)
	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

func (r ProviderRepository) DeletePrice(ctx context.Context, priceId uuid.UUID) (bool, error) {
	result := r.dbContext.db.WithContext(ctx).Delete(&providerPriceSqlModel{}, priceId)
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
	result := r.dbContext.db.WithContext(ctx).
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
	result := r.dbContext.db.WithContext(ctx).
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
	sqlModels := slicesx.Select(labelIds, func(labelId uuid.UUID) providerLabelSqlModel {
		return newProviderLabelSqlModel(providerId, labelId)
	})
	result := r.dbContext.db.WithContext(ctx).Create(&sqlModels)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r ProviderRepository) createPlans(ctx context.Context, providerId uuid.UUID, plans []provider.Plan) error {
	planSqlModels := slicesx.Select(plans, func(plan provider.Plan) providerPlanSqlModel {
		return newProviderPlanSqlModel(providerId, plan)
	})

	result := r.dbContext.db.WithContext(ctx).
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

	result = r.dbContext.db.WithContext(ctx).Create(&priceSqlModels)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r ProviderRepository) update(ctx context.Context, dirtyProvider provider.Provider) error {
	sqlModel := newProviderSqlModel(dirtyProvider)

	result := r.dbContext.db.WithContext(ctx).
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
		r.dbContext.db,
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
		r.dbContext.db,
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
		r.dbContext.db,
		prices,
		func(price provider.Price) providerPriceSqlModel {
			return newProviderPriceSqlModel(planId, price)
		},
	)
}
