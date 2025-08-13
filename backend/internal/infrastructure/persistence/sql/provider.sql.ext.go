package sql

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type ProviderRow struct {
	Provider      Provider
	ProviderPlan  *ProviderPlan
	ProviderPrice *ProviderPrice
	ProviderLabel *uuid.UUID
}

func (q *Queries) GetSystemProviders(ctx context.Context) ([]ProviderRow, int64, error) {
	rows, err := q.getSystemProviders(ctx)
	if err != nil {
		return nil, 0, err
	}
	if len(rows) == 0 {
		return nil, 0, nil
	}
	totalCount := rows[0].TotalCount
	results := make([]ProviderRow, len(rows))
	for i, row := range rows {
		var providerPlan *ProviderPlan
		var providerPrice *ProviderPrice
		var labelId *uuid.UUID
		if row.ProviderPlansID != nil {
			providerPlan = &ProviderPlan{
				ID:          *row.ProviderPlansID,
				Name:        *row.ProviderPlansName,
				Description: row.ProviderPlansDescription,
				ProviderID:  *row.ProviderPlansProviderID,
				CreatedAt:   *row.ProviderPlansCreatedAt,
				UpdatedAt:   *row.ProviderPlansUpdatedAt,
				Etag:        *row.ProviderPlansEtag,
			}
		}
		if row.ProviderPricesID != nil {
			providerPrice = &ProviderPrice{
				ID:        *row.ProviderPricesID,
				StartDate: *row.ProviderPricesStartDate,
				EndDate:   row.ProviderPricesEndDate,
				Currency:  *row.ProviderPricesCurrency,
				Amount:    *row.ProviderPricesAmount,
				PlanID:    *row.ProviderPricesPlanID,
				CreatedAt: *row.ProviderPricesCreatedAt,
				UpdatedAt: *row.ProviderPricesUpdatedAt,
			}
		}

		if row.ProviderLabelsLabelID != nil {
			labelId = row.ProviderLabelsLabelID
		}

		results[i] = ProviderRow{
			Provider: Provider{
				ID:            row.ProvidersID,
				OwnerType:     row.ProvidersOwnerType,
				OwnerFamilyID: row.ProvidersOwnerFamilyID,
				OwnerUserID:   row.ProvidersOwnerUserID,
				Name:          row.ProvidersName,
				Key:           row.ProvidersKey,
				Description:   row.ProvidersDescription,
				IconUrl:       row.ProvidersIconUrl,
			},
			ProviderPlan:  providerPlan,
			ProviderPrice: providerPrice,
			ProviderLabel: labelId,
		}
	}

	return results, totalCount, nil
}

func (q *Queries) GetProviders(ctx context.Context, limit, offset int32) ([]ProviderRow, int64, error) {
	rows, err := q.getProviders(ctx, getProvidersParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, 0, err
	}
	if len(rows) == 0 {
		return nil, 0, nil
	}
	totalCount := rows[0].TotalCount
	results := make([]ProviderRow, len(rows))
	for i, row := range rows {
		var providerPlan *ProviderPlan
		var providerPrice *ProviderPrice
		var labelId *uuid.UUID
		if row.ProviderPlansID != nil {
			providerPlan = &ProviderPlan{
				ID:          *row.ProviderPlansID,
				Name:        *row.ProviderPlansName,
				Description: row.ProviderPlansDescription,
				ProviderID:  *row.ProviderPlansProviderID,
				CreatedAt:   *row.ProviderPlansCreatedAt,
				UpdatedAt:   *row.ProviderPlansUpdatedAt,
				Etag:        *row.ProviderPlansEtag,
			}
		}

		if row.ProviderPricesID != nil {
			providerPrice = &ProviderPrice{
				ID:        *row.ProviderPricesID,
				StartDate: *row.ProviderPricesStartDate,
				EndDate:   row.ProviderPricesEndDate,
				Currency:  *row.ProviderPricesCurrency,
				Amount:    *row.ProviderPricesAmount,
				PlanID:    *row.ProviderPricesPlanID,
				CreatedAt: *row.ProviderPricesCreatedAt,
				UpdatedAt: *row.ProviderPricesUpdatedAt,
				Etag:      *row.ProviderPricesEtag,
			}
		}

		if row.ProviderLabelsLabelID != nil {
			labelId = row.ProviderLabelsLabelID
		}

		results[i] = ProviderRow{
			Provider: Provider{
				ID:             row.ProvidersID,
				OwnerType:      row.ProvidersOwnerType,
				OwnerFamilyID:  row.ProvidersOwnerFamilyID,
				OwnerUserID:    row.ProvidersOwnerUserID,
				Name:           row.ProvidersName,
				Key:            row.ProvidersKey,
				Description:    row.ProvidersDescription,
				IconUrl:        row.ProvidersIconUrl,
				Url:            row.ProvidersUrl,
				PricingPageUrl: row.ProvidersPricingPageUrl,
				CreatedAt:      row.ProvidersCreatedAt,
				UpdatedAt:      row.ProvidersUpdatedAt,
				Etag:           row.ProvidersEtag,
			},
			ProviderPlan:  providerPlan,
			ProviderPrice: providerPrice,
			ProviderLabel: labelId,
		}
	}

	return results, totalCount, nil
}

func (q *Queries) GetProvidersForUser(ctx context.Context, userId string, limit, offset int32) (
	[]ProviderRow,
	int64,
	error) {
	rows, err := q.getProvidersForUser(ctx, getProvidersForUserParams{
		OwnerUserID: &userId,
		Limit:       limit,
		Offset:      offset,
	})
	if err != nil {
		return nil, 0, err
	}
	if len(rows) == 0 {
		return nil, 0, nil
	}
	totalCount := rows[0].TotalCount
	results := make([]ProviderRow, len(rows))
	for i, row := range rows {
		var providerPlan *ProviderPlan
		var providerPrice *ProviderPrice
		var labelId *uuid.UUID
		if row.ProviderPlansID != nil {
			providerPlan = &ProviderPlan{
				ID:          *row.ProviderPlansID,
				Name:        *row.ProviderPlansName,
				Description: row.ProviderPlansDescription,
				ProviderID:  *row.ProviderPlansProviderID,
				CreatedAt:   *row.ProviderPlansCreatedAt,
				UpdatedAt:   *row.ProviderPlansUpdatedAt,
				Etag:        *row.ProviderPlansEtag,
			}
		}

		if row.ProviderPricesID != nil {
			providerPrice = &ProviderPrice{
				ID:        *row.ProviderPricesID,
				StartDate: *row.ProviderPricesStartDate,
				EndDate:   row.ProviderPricesEndDate,
				Currency:  *row.ProviderPricesCurrency,
				Amount:    *row.ProviderPricesAmount,
				PlanID:    *row.ProviderPricesPlanID,
				CreatedAt: *row.ProviderPricesCreatedAt,
				UpdatedAt: *row.ProviderPricesUpdatedAt,
				Etag:      *row.ProviderPricesEtag,
			}
		}

		if row.ProviderLabelsLabelID != nil {
			labelId = row.ProviderLabelsLabelID
		}

		results[i] = ProviderRow{
			Provider: Provider{
				ID:             row.ProvidersID,
				OwnerType:      row.ProvidersOwnerType,
				OwnerFamilyID:  row.ProvidersOwnerFamilyID,
				OwnerUserID:    row.ProvidersOwnerUserID,
				Name:           row.ProvidersName,
				Key:            row.ProvidersKey,
				Description:    row.ProvidersDescription,
				IconUrl:        row.ProvidersIconUrl,
				Url:            row.ProvidersUrl,
				PricingPageUrl: row.ProvidersPricingPageUrl,
				CreatedAt:      row.ProvidersCreatedAt,
				UpdatedAt:      row.ProvidersUpdatedAt,
				Etag:           row.ProvidersEtag,
			},
			ProviderPlan:  providerPlan,
			ProviderPrice: providerPrice,
			ProviderLabel: labelId,
		}
	}

	return results, totalCount, nil
}

func (q *Queries) GetProvidersForUserWithSearch(ctx context.Context, userId, searchText string, limit, offset int32) (
	[]ProviderRow,
	int64,
	error) {
	rows, err := q.getProvidersForUserWithSearch(ctx, getProvidersForUserWithSearchParams{
		Name:        fmt.Sprintf("%%%s%%", searchText),
		OwnerUserID: &userId,
		Limit:       limit,
		Offset:      offset,
	})
	if err != nil {
		return nil, 0, err
	}
	if len(rows) == 0 {
		return nil, 0, nil
	}
	totalCount := rows[0].TotalCount
	results := make([]ProviderRow, len(rows))
	for i, row := range rows {
		var providerPlan *ProviderPlan
		var providerPrice *ProviderPrice
		var labelId *uuid.UUID
		if row.ProviderPlansID != nil {
			providerPlan = &ProviderPlan{
				ID:          *row.ProviderPlansID,
				Name:        *row.ProviderPlansName,
				Description: row.ProviderPlansDescription,
				ProviderID:  *row.ProviderPlansProviderID,
				CreatedAt:   *row.ProviderPlansCreatedAt,
				UpdatedAt:   *row.ProviderPlansUpdatedAt,
				Etag:        *row.ProviderPlansEtag,
			}
		}

		if row.ProviderPricesID != nil {
			providerPrice = &ProviderPrice{
				ID:        *row.ProviderPricesID,
				StartDate: *row.ProviderPricesStartDate,
				EndDate:   row.ProviderPricesEndDate,
				Currency:  *row.ProviderPricesCurrency,
				Amount:    *row.ProviderPricesAmount,
				PlanID:    *row.ProviderPricesPlanID,
				CreatedAt: *row.ProviderPricesCreatedAt,
				UpdatedAt: *row.ProviderPricesUpdatedAt,
				Etag:      *row.ProviderPricesEtag,
			}
		}

		if row.ProviderLabelsLabelID != nil {
			labelId = row.ProviderLabelsLabelID
		}

		results[i] = ProviderRow{
			Provider: Provider{
				ID:             row.ProvidersID,
				OwnerType:      row.ProvidersOwnerType,
				OwnerFamilyID:  row.ProvidersOwnerFamilyID,
				OwnerUserID:    row.ProvidersOwnerUserID,
				Name:           row.ProvidersName,
				Key:            row.ProvidersKey,
				Description:    row.ProvidersDescription,
				IconUrl:        row.ProvidersIconUrl,
				Url:            row.ProvidersUrl,
				PricingPageUrl: row.ProvidersPricingPageUrl,
				CreatedAt:      row.ProvidersCreatedAt,
				UpdatedAt:      row.ProvidersUpdatedAt,
				Etag:           row.ProvidersEtag,
			},
			ProviderPlan:  providerPlan,
			ProviderPrice: providerPrice,
			ProviderLabel: labelId,
		}
	}

	return results, totalCount, nil
}

func (q *Queries) GetProviderById(ctx context.Context, id uuid.UUID) ([]ProviderRow, error) {
	rows, err := q.getProviderById(ctx, id)
	if err != nil {
		return nil, err
	}
	results := make([]ProviderRow, len(rows))
	for i, row := range rows {
		var providerPlan *ProviderPlan
		var providerPrice *ProviderPrice
		var labelId *uuid.UUID
		if row.ProviderPlansID != nil {
			providerPlan = &ProviderPlan{
				ID:          *row.ProviderPlansID,
				Name:        *row.ProviderPlansName,
				Description: row.ProviderPlansDescription,
				ProviderID:  *row.ProviderPlansProviderID,
				CreatedAt:   *row.ProviderPlansCreatedAt,
				UpdatedAt:   *row.ProviderPlansUpdatedAt,
				Etag:        *row.ProviderPlansEtag,
			}
		}

		if row.ProviderPricesID != nil {
			providerPrice = &ProviderPrice{
				ID:        *row.ProviderPricesID,
				StartDate: *row.ProviderPricesStartDate,
				EndDate:   row.ProviderPricesEndDate,
				Currency:  *row.ProviderPricesCurrency,
				Amount:    *row.ProviderPricesAmount,
				PlanID:    *row.ProviderPricesPlanID,
				CreatedAt: *row.ProviderPricesCreatedAt,
				UpdatedAt: *row.ProviderPricesUpdatedAt,
				Etag:      *row.ProviderPricesEtag,
			}
		}

		if row.ProviderLabelsLabelID != nil {
			labelId = row.ProviderLabelsLabelID
		}

		results[i] = ProviderRow{
			Provider: Provider{
				ID:             row.ProvidersID,
				OwnerType:      row.ProvidersOwnerType,
				OwnerFamilyID:  row.ProvidersOwnerFamilyID,
				OwnerUserID:    row.ProvidersOwnerUserID,
				Name:           row.ProvidersName,
				Key:            row.ProvidersKey,
				Description:    row.ProvidersDescription,
				IconUrl:        row.ProvidersIconUrl,
				Url:            row.ProvidersUrl,
				PricingPageUrl: row.ProvidersPricingPageUrl,
				CreatedAt:      row.ProvidersCreatedAt,
				UpdatedAt:      row.ProvidersUpdatedAt,
				Etag:           row.ProvidersEtag,
			},
			ProviderPlan:  providerPlan,
			ProviderPrice: providerPrice,
			ProviderLabel: labelId,
		}
	}

	return results, nil
}

func (q *Queries) GetProviderByIdForUser(ctx context.Context, userId string, id uuid.UUID) ([]ProviderRow, error) {
	rows, err := q.getProviderByIdForUser(ctx, getProviderByIdForUserParams{
		OwnerUserID: &userId,
		ID:          id,
	})
	if err != nil {
		return nil, err
	}
	results := make([]ProviderRow, len(rows))
	for i, row := range rows {
		var providerPlan *ProviderPlan
		var providerPrice *ProviderPrice
		var labelId *uuid.UUID
		if row.ProviderPlansID != nil {
			providerPlan = &ProviderPlan{
				ID:          *row.ProviderPlansID,
				Name:        *row.ProviderPlansName,
				Description: row.ProviderPlansDescription,
				ProviderID:  *row.ProviderPlansProviderID,
				CreatedAt:   *row.ProviderPlansCreatedAt,
				UpdatedAt:   *row.ProviderPlansUpdatedAt,
				Etag:        *row.ProviderPlansEtag,
			}
		}

		if row.ProviderPricesID != nil {
			providerPrice = &ProviderPrice{
				ID:        *row.ProviderPricesID,
				StartDate: *row.ProviderPricesStartDate,
				EndDate:   row.ProviderPricesEndDate,
				Currency:  *row.ProviderPricesCurrency,
				Amount:    *row.ProviderPricesAmount,
				PlanID:    *row.ProviderPricesPlanID,
				CreatedAt: *row.ProviderPricesCreatedAt,
				UpdatedAt: *row.ProviderPricesUpdatedAt,
				Etag:      *row.ProviderPricesEtag,
			}
		}

		if row.ProviderLabelsLabelID != nil {
			labelId = row.ProviderLabelsLabelID
		}

		results[i] = ProviderRow{
			Provider: Provider{
				ID:             row.ProvidersID,
				OwnerType:      row.ProvidersOwnerType,
				OwnerFamilyID:  row.ProvidersOwnerFamilyID,
				OwnerUserID:    row.ProvidersOwnerUserID,
				Name:           row.ProvidersName,
				Key:            row.ProvidersKey,
				Description:    row.ProvidersDescription,
				IconUrl:        row.ProvidersIconUrl,
				Url:            row.ProvidersUrl,
				PricingPageUrl: row.ProvidersPricingPageUrl,
				CreatedAt:      row.ProvidersCreatedAt,
				UpdatedAt:      row.ProvidersUpdatedAt,
				Etag:           row.ProvidersEtag,
			},
			ProviderPlan:  providerPlan,
			ProviderPrice: providerPrice,
			ProviderLabel: labelId,
		}
	}

	return results, nil
}
