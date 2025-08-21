package persistence

import (
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/text/currency"

	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/internal/domain/provider"
	"github.com/oleexo/subtracker/internal/infrastructure/persistence/jet/app/public/model"
	"github.com/oleexo/subtracker/pkg/slicesx"
)

func createProviderPriceFromJet(jetModel model.ProviderPrices) provider.Price {
	p := provider.NewPrice(
		jetModel.ID,
		jetModel.StartDate,
		jetModel.EndDate,
		currency.MustParseISO(jetModel.Currency),
		jetModel.Amount,
		jetModel.CreatedAt,
		jetModel.UpdatedAt,
	)

	p.Clean()
	return p
}

func createProviderPlanFromJet(jetPlan model.ProviderPlans, prices []provider.Price) provider.Plan {
	p := provider.NewPlan(
		jetPlan.ID,
		jetPlan.Name,
		jetPlan.Description,
		prices,
		jetPlan.CreatedAt,
		jetPlan.UpdatedAt,
	)

	p.Clean()
	return p
}

func createProviderFromJet(jetModel model.Providers, plans []provider.Plan, labels []uuid.UUID) provider.Provider {
	ownerType := auth.MustParseOwnerType(jetModel.OwnerType)
	owner := auth.NewOwner(ownerType, jetModel.OwnerFamilyID, jetModel.OwnerUserID)

	p := provider.NewProvider(
		jetModel.ID,
		jetModel.Name,
		jetModel.Key,
		jetModel.Description,
		jetModel.IconURL,
		jetModel.URL,
		jetModel.PricingPageURL,
		labels,
		plans,
		owner,
		jetModel.CreatedAt,
		jetModel.UpdatedAt,
	)

	p.Clean()
	return p
}

// createProviderFromJetRows converts jet rows to domain providers with plans and prices
func createProviderFromJetRows(rows []struct {
	Providers      model.Providers       `json:"providers"`
	ProviderPlans  *model.ProviderPlans  `json:"provider_plans"`
	ProviderPrices *model.ProviderPrices `json:"provider_prices"`
	ProviderLabels *model.ProviderLabels `json:"provider_labels"`
}) []provider.Provider {
	if len(rows) == 0 {
		return nil
	}

	providers := make(map[uuid.UUID]model.Providers)
	planSet := make(map[uuid.UUID]struct{})
	providerPlans := make(map[uuid.UUID][]model.ProviderPlans)
	priceSet := make(map[uuid.UUID]struct{})
	planPrices := make(map[uuid.UUID][]model.ProviderPrices)
	labelSet := make(map[string]struct{})
	providerLabels := make(map[uuid.UUID][]uuid.UUID)

	for _, row := range rows {
		jetProvider := row.Providers
		if _, ok := providers[jetProvider.ID]; !ok {
			providers[jetProvider.ID] = jetProvider
		}

		if row.ProviderPlans != nil && row.ProviderPlans.ID != uuid.Nil {
			if _, ok := planSet[row.ProviderPlans.ID]; !ok {
				planSet[row.ProviderPlans.ID] = struct{}{}
				providerPlans[row.ProviderPlans.ProviderID] = append(providerPlans[row.ProviderPlans.ProviderID], *row.ProviderPlans)
			}
		}

		if row.ProviderPrices != nil && row.ProviderPrices.ID != uuid.Nil {
			if _, ok := priceSet[row.ProviderPrices.ID]; !ok {
				priceSet[row.ProviderPrices.ID] = struct{}{}
				planPrices[row.ProviderPrices.PlanID] = append(planPrices[row.ProviderPrices.PlanID], *row.ProviderPrices)
			}
		}

		if row.ProviderLabels != nil && row.ProviderLabels.LabelID != uuid.Nil {
			key := fmt.Sprintf("%s:%s", jetProvider.ID, row.ProviderLabels.LabelID)
			if _, ok := labelSet[key]; !ok {
				labelSet[key] = struct{}{}
				providerLabels[jetProvider.ID] = append(providerLabels[jetProvider.ID], row.ProviderLabels.LabelID)
			}
		}
	}

	results := make([]provider.Provider, len(providers))
	count := 0
	for providerId, jetProvider := range providers {
		var plans []provider.Plan
		jetPlans, planExists := providerPlans[providerId]
		if planExists {
			plans = slicesx.Select(jetPlans, func(jetPlan model.ProviderPlans) provider.Plan {
				var prices []provider.Price
				jetPrices, priceExists := planPrices[jetPlan.ID]
				if priceExists {
					prices = slicesx.Select(jetPrices, func(jetPrice model.ProviderPrices) provider.Price {
						return createProviderPriceFromJet(jetPrice)
					})
				}

				return createProviderPlanFromJet(jetPlan, prices)
			})
		}

		jetLabels, _ := providerLabels[providerId]
		results[count] = createProviderFromJet(jetProvider, plans, jetLabels)
		count++
	}

	return results
}

// createProviderFromJetRowsWithCount converts jet rows with count to domain providers
func createProviderFromJetRowsWithCount(rows []struct {
	Providers      model.Providers       `json:"providers"`
	ProviderPlans  *model.ProviderPlans  `json:"provider_plans"`
	ProviderPrices *model.ProviderPrices `json:"provider_prices"`
	ProviderLabels *model.ProviderLabels `json:"provider_labels"`
	TotalCount     int64                 `json:"total_count"`
}) []provider.Provider {
	// Convert to the simpler row structure
	simpleRows := slicesx.Select(rows, func(row struct {
		Providers      model.Providers       `json:"providers"`
		ProviderPlans  *model.ProviderPlans  `json:"provider_plans"`
		ProviderPrices *model.ProviderPrices `json:"provider_prices"`
		ProviderLabels *model.ProviderLabels `json:"provider_labels"`
		TotalCount     int64                 `json:"total_count"`
	}) struct {
		Providers      model.Providers       `json:"providers"`
		ProviderPlans  *model.ProviderPlans  `json:"provider_plans"`
		ProviderPrices *model.ProviderPrices `json:"provider_prices"`
		ProviderLabels *model.ProviderLabels `json:"provider_labels"`
	} {
		return struct {
			Providers      model.Providers       `json:"providers"`
			ProviderPlans  *model.ProviderPlans  `json:"provider_plans"`
			ProviderPrices *model.ProviderPrices `json:"provider_prices"`
			ProviderLabels *model.ProviderLabels `json:"provider_labels"`
		}{
			Providers:      row.Providers,
			ProviderPlans:  row.ProviderPlans,
			ProviderPrices: row.ProviderPrices,
			ProviderLabels: row.ProviderLabels,
		}
	})

	return createProviderFromJetRows(simpleRows)
}
