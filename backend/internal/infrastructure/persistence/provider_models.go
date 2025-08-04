package persistence

import (
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/text/currency"

	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/internal/domain/provider"
	"github.com/oleexo/subtracker/internal/infrastructure/persistence/sql"
	"github.com/oleexo/subtracker/pkg/slicesx"
)

func createProviderPriceFromSqlc(sqlcModel sql.ProviderPrice) provider.Price {
	p := provider.NewPrice(
		sqlcModel.ID,
		sqlcModel.StartDate,
		sqlcModel.EndDate,
		currency.MustParseISO(sqlcModel.Currency),
		sqlcModel.Amount,
		sqlcModel.CreatedAt,
		sqlcModel.UpdatedAt,
	)

	p.Clean()
	return p
}

func createProviderPlanFromSqlc(sqlcPlan sql.ProviderPlan, prices []provider.Price) provider.Plan {
	p := provider.NewPlan(
		sqlcPlan.ID,
		sqlcPlan.Name,
		sqlcPlan.Description,
		prices,
		sqlcPlan.CreatedAt,
		sqlcPlan.UpdatedAt,
	)

	p.Clean()
	return p
}

func createProviderFromSqlc(sqlcModel sql.Provider, plans []provider.Plan, labels []uuid.UUID) provider.Provider {
	ownerType := auth.MustParseOwnerType(sqlcModel.OwnerType)
	owner := auth.NewOwner(ownerType, sqlcModel.OwnerFamilyID, sqlcModel.OwnerUserID)

	p := provider.NewProvider(
		sqlcModel.ID,
		sqlcModel.Name,
		sqlcModel.Key,
		sqlcModel.Description,
		sqlcModel.IconUrl,
		sqlcModel.Url,
		sqlcModel.PricingPageUrl,
		labels,
		plans,
		owner,
		sqlcModel.CreatedAt,
		sqlcModel.UpdatedAt,
	)

	p.Clean()
	return p
}

// createProviderFromSqlcRows converts SQLC rows to domain providers with plans and prices
func createProviderFromSqlcRows[T any](
	rows []T,
	getProviderFunc func(T) sql.Provider,
	getPlanFunc func(T) *sql.ProviderPlan,
	getPriceFunc func(T) *sql.ProviderPrice,
	getLabelFunc func(T) *uuid.UUID) []provider.Provider {

	if len(rows) == 0 {
		return nil
	}

	providers := make(map[uuid.UUID]sql.Provider)
	planSet := make(map[uuid.UUID]struct{})
	providerPlans := make(map[uuid.UUID][]sql.ProviderPlan)
	priceSet := make(map[uuid.UUID]struct{})
	planPrices := make(map[uuid.UUID][]sql.ProviderPrice)
	labelSet := make(map[string]struct{})
	providerLabels := make(map[uuid.UUID][]uuid.UUID)

	for _, row := range rows {
		sqlcProvider := getProviderFunc(row)
		if _, ok := providers[sqlcProvider.ID]; !ok {
			providers[sqlcProvider.ID] = sqlcProvider
		}
		sqlcPlan := getPlanFunc(row)
		if sqlcPlan != nil {
			if _, ok := planSet[sqlcPlan.ID]; !ok {
				planSet[sqlcPlan.ID] = struct{}{}
				providerPlans[sqlcPlan.ProviderID] = append(providerPlans[sqlcPlan.ProviderID], *sqlcPlan)
			}
		}
		sqlcPrice := getPriceFunc(row)
		if sqlcPrice != nil {
			if _, ok := priceSet[sqlcPrice.ID]; !ok {
				priceSet[sqlcPrice.ID] = struct{}{}
				planPrices[sqlcPrice.PlanID] = append(planPrices[sqlcPrice.PlanID], *sqlcPrice)
			}
		}

		sqlcLabelId := getLabelFunc(row)
		if sqlcLabelId != nil {
			key := fmt.Sprintf("%s:%s", sqlcProvider.ID, *sqlcLabelId)
			if _, ok := labelSet[key]; !ok {
				labelSet[key] = struct{}{}
				providerLabels[sqlcProvider.ID] = append(providerLabels[sqlcProvider.ID], *sqlcLabelId)
			}
		}
	}

	results := make([]provider.Provider, len(providers))
	count := 0
	for providerId, sqlcProvider := range providers {
		var plans []provider.Plan
		sqlcPlans, planExists := providerPlans[providerId]
		if planExists {
			plans = slicesx.Select(sqlcPlans, func(sqlcPlan sql.ProviderPlan) provider.Plan {
				var prices []provider.Price
				sqlcPrices, priceExists := planPrices[sqlcPlan.ID]
				if priceExists {
					prices = slicesx.Select(sqlcPrices, func(sqlcPrice sql.ProviderPrice) provider.Price {
						return createProviderPriceFromSqlc(sqlcPrice)
					})
				}

				return createProviderPlanFromSqlc(sqlcPlan, prices)
			})
		}

		sqlcLabels, _ := providerLabels[providerId]
		results[count] = createProviderFromSqlc(sqlcProvider, plans, sqlcLabels)
		count++
	}

	return results
}
