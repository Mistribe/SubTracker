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
	var newProviders []provider.Provider
	for _, prov := range providers {
		if !prov.IsExists() {
			newProviders = append(newProviders, prov)
		} else {
			if err := r.update(ctx, prov); err != nil {
				return err
			}
		}
	}

	if len(newProviders) > 0 {
		if err := r.create(ctx, newProviders); err != nil {
			return err
		}
	}

	for _, prov := range providers {
		prov.Clean()
	}

	return nil
}

func (r ProviderRepository) Delete(ctx context.Context, providerId uuid.UUID) (bool, error) {
	err := r.dbContext.GetQueries(ctx).DeleteProvider(ctx, providerId)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r ProviderRepository) DeletePlan(ctx context.Context, planId uuid.UUID) (bool, error) {
	err := r.dbContext.GetQueries(ctx).DeleteProviderPlan(ctx, planId)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r ProviderRepository) DeletePrice(ctx context.Context, priceId uuid.UUID) (bool, error) {
	err := r.dbContext.GetQueries(ctx).DeleteProviderPrice(ctx, priceId)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r ProviderRepository) Exists(ctx context.Context, ids ...uuid.UUID) (bool, error) {
	if len(ids) == 0 {
		return true, nil
	}

	count, err := r.dbContext.GetQueries(ctx).IsProviderExists(ctx, ids)
	if err != nil {
		return false, err
	}
	return count > int64(len(ids)), nil
}

func (r ProviderRepository) create(ctx context.Context, providers []provider.Provider) error {
	args := slicesx.Select(providers, func(prov provider.Provider) sql.CreateProvidersParams {
		model := sql.CreateProvidersParams{
			ID:             prov.Id(),
			OwnerType:      prov.Owner().Type().String(),
			Name:           prov.Name(),
			Key:            prov.Key(),
			Description:    prov.Description(),
			IconUrl:        prov.IconUrl(),
			Url:            prov.Url(),
			PricingPageUrl: prov.PricingPageUrl(),
			CreatedAt:      prov.CreatedAt(),
			UpdatedAt:      prov.UpdatedAt(),
			Etag:           prov.ETag(),
		}

		switch prov.Owner().Type() {
		case auth.PersonalOwnerType:
			userId := prov.Owner().UserId()
			model.OwnerUserID = &userId
		case auth.FamilyOwnerType:
			familyId := prov.Owner().FamilyId()
			model.OwnerFamilyID = &familyId
		}

		return model
	})

	if _, err := r.dbContext.GetQueries(ctx).CreateProviders(ctx, args); err != nil {
		return err
	}

	labelArgs := slicesx.SelectMany(providers, func(prov provider.Provider) []sql.CreateProviderLabelsParams {
		return slicesx.Select(prov.Labels().Values(), func(labelId uuid.UUID) sql.CreateProviderLabelsParams {
			return sql.CreateProviderLabelsParams{
				ProviderID: prov.Id(),
				LabelID:    labelId,
			}
		})
	})

	if len(labelArgs) > 0 {
		if _, err := r.dbContext.GetQueries(ctx).CreateProviderLabels(ctx, labelArgs); err != nil {
			return err
		}
	}

	planArgs := slicesx.SelectMany(providers, func(prov provider.Provider) []sql.CreateProviderPlansParams {
		return slicesx.Select(prov.Plans().Values(), func(plan provider.Plan) sql.CreateProviderPlansParams {
			return sql.CreateProviderPlansParams{
				ID:          plan.Id(),
				ProviderID:  prov.Id(),
				Name:        plan.Name(),
				Description: plan.Description(),
				CreatedAt:   plan.CreatedAt(),
				UpdatedAt:   plan.UpdatedAt(),
				Etag:        plan.ETag(),
			}
		})
	})

	if len(planArgs) > 0 {
		if _, err := r.dbContext.GetQueries(ctx).CreateProviderPlans(ctx, planArgs); err != nil {
			return err
		}
	}

	priceArgs := slicesx.SelectMany(providers, func(prov provider.Provider) []sql.CreateProviderPricesParams {
		return slicesx.SelectMany(prov.Plans().Values(), func(plan provider.Plan) []sql.CreateProviderPricesParams {
			return slicesx.Select(plan.Prices().Values(), func(price provider.Price) sql.CreateProviderPricesParams {
				return sql.CreateProviderPricesParams{
					ID:        price.Id(),
					PlanID:    plan.Id(),
					Currency:  price.Currency().String(),
					Amount:    price.Amount(),
					StartDate: price.StartDate(),
					EndDate:   price.EndDate(),
					CreatedAt: price.CreatedAt(),
					UpdatedAt: price.UpdatedAt(),
					Etag:      price.ETag(),
				}
			})
		})
	})

	if len(priceArgs) > 0 {
		if _, err := r.dbContext.GetQueries(ctx).CreateProviderPrices(ctx, priceArgs); err != nil {
			return err
		}
	}

	return nil
}

func (r ProviderRepository) update(ctx context.Context, dirtyProvider provider.Provider) error {
	model := sql.UpdateProviderParams{
		ID:             dirtyProvider.Id(),
		OwnerType:      dirtyProvider.Owner().Type().String(),
		Name:           dirtyProvider.Name(),
		Key:            dirtyProvider.Key(),
		Description:    dirtyProvider.Description(),
		IconUrl:        dirtyProvider.IconUrl(),
		Url:            dirtyProvider.Url(),
		PricingPageUrl: dirtyProvider.PricingPageUrl(),
		UpdatedAt:      dirtyProvider.UpdatedAt(),
	}

	switch dirtyProvider.Owner().Type() {
	case auth.PersonalOwnerType:
		userId := dirtyProvider.Owner().UserId()
		model.OwnerUserID = &userId
	case auth.FamilyOwnerType:
		familyId := dirtyProvider.Owner().FamilyId()
		model.OwnerFamilyID = &familyId
	}

	if err := r.dbContext.GetQueries(ctx).UpdateProvider(ctx, model); err != nil {
		return err
	}

	if dirtyProvider.Labels().HasChanges() {
		if err := saveTrackedSlice(ctx, r.dbContext, dirtyProvider.Labels(),
			func(ctx context.Context, queries *sql.Queries, entities []uuid.UUID) error {
				args := slicesx.Select(entities, func(labelId uuid.UUID) sql.CreateProviderLabelsParams {
					return sql.CreateProviderLabelsParams{
						LabelID:    labelId,
						ProviderID: dirtyProvider.Id(),
					}
				})
				_, qErr := queries.CreateProviderLabels(ctx, args)
				return qErr
			},
			func(ctx context.Context, queries *sql.Queries, entity uuid.UUID) error {
				panic("can't update provider_labels relation")
			},
			func(ctx context.Context, queries *sql.Queries, entity uuid.UUID) error {
				return queries.DeleteProviderLabel(ctx, sql.DeleteProviderLabelParams{
					ProviderID: dirtyProvider.Id(),
					LabelID:    entity,
				})
			}); err != nil {
			return err
		}

		dirtyProvider.Labels().ClearChanges()
	}

	if dirtyProvider.Plans().HasChanges() {
		if err := saveTrackedSlice(ctx,
			r.dbContext,
			dirtyProvider.Plans(),
			func(ctx context.Context, queries *sql.Queries, plans []provider.Plan) error {
				args := slicesx.Select(plans, func(plan provider.Plan) sql.CreateProviderPlansParams {
					return sql.CreateProviderPlansParams{
						ID:          plan.Id(),
						ProviderID:  dirtyProvider.Id(),
						Name:        plan.Name(),
						Description: plan.Description(),
						CreatedAt:   plan.CreatedAt(),
						UpdatedAt:   plan.UpdatedAt(),
						Etag:        plan.ETag(),
					}
				})
				_, qErr := queries.CreateProviderPlans(ctx, args)
				return qErr
			},
			func(ctx context.Context, queries *sql.Queries, entity provider.Plan) error {
				return queries.UpdateProviderPlan(ctx, sql.UpdateProviderPlanParams{
					ID:          entity.Id(),
					Name:        entity.Name(),
					Description: entity.Description(),
					UpdatedAt:   entity.UpdatedAt(),
					Etag:        entity.ETag(),
					ProviderID:  dirtyProvider.Id(),
				})
			},
			func(ctx context.Context, queries *sql.Queries, entity provider.Plan) error {
				return queries.DeleteProviderPlan(ctx, entity.Id())
			},
		); err != nil {
			return err
		}

		for updatedPlan := range dirtyProvider.Plans().Updated() {
			if updatedPlan.Prices().HasChanges() {
				if err := r.savePrices(ctx, updatedPlan.Id(), updatedPlan.Prices()); err != nil {
					return err
				}
			}
		}

		for updatePlan := range dirtyProvider.Plans().Added() {
			if err := r.savePrices(ctx, updatePlan.Id(), updatePlan.Prices()); err != nil {
				return err
			}
		}

		dirtyProvider.Plans().ClearChanges()
	}

	return nil
}

func (r ProviderRepository) savePrices(ctx context.Context,
	planId uuid.UUID,
	trackedPrices *slicesx.Tracked[provider.Price]) error {
	err := saveTrackedSlice(ctx,
		r.dbContext,
		trackedPrices,
		func(ctx context.Context, queries *sql.Queries, entities []provider.Price) error {
			args := slicesx.Select(entities, func(price provider.Price) sql.CreateProviderPricesParams {
				return sql.CreateProviderPricesParams{
					ID:        price.Id(),
					PlanID:    planId,
					Currency:  price.Currency().String(),
					Amount:    price.Amount(),
					StartDate: price.StartDate(),
					EndDate:   price.EndDate(),
					CreatedAt: price.CreatedAt(),
					UpdatedAt: price.UpdatedAt(),
					Etag:      price.ETag(),
				}

			})
			_, qErr := queries.CreateProviderPrices(ctx, args)
			return qErr
		},
		func(ctx context.Context, queries *sql.Queries, entity provider.Price) error {
			return queries.UpdateProviderPrice(ctx, sql.UpdateProviderPriceParams{
				ID:        entity.Id(),
				Currency:  entity.Currency().String(),
				Amount:    entity.Amount(),
				StartDate: entity.StartDate(),
				EndDate:   entity.EndDate(),
				UpdatedAt: entity.UpdatedAt(),
				Etag:      entity.ETag(),
				PlanID:    planId,
			})
		},
		func(ctx context.Context, queries *sql.Queries, entity provider.Price) error {
			return queries.DeleteProviderPrice(ctx, entity.Id())
		},
	)

	if err != nil {
		return err
	}

	trackedPrices.ClearChanges()
	return nil
}
