package persistence

import (
	"github.com/google/uuid"
	"golang.org/x/text/currency"

	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/internal/domain/provider"
	"github.com/oleexo/subtracker/internal/infrastructure/persistence/sql"
)

func createProviderPriceFromSqlc(sqlcModel sql.ProviderPrice) provider.Price {
	return provider.NewPrice(
		sqlcModel.ID,
		sqlcModel.StartDate,
		sqlcModel.EndDate,
		currency.MustParseISO(sqlcModel.Currency),
		sqlcModel.Amount,
		sqlcModel.CreatedAt,
		sqlcModel.UpdatedAt,
	)
}

func createProviderFromSqlc(sqlcModel sql.Provider) provider.Provider {
	ownerType := auth.MustParseOwnerType(sqlcModel.OwnerType)
	owner := auth.NewOwner(ownerType, sqlcModel.OwnerFamilyID, sqlcModel.OwnerUserID)

	return provider.NewProvider(
		sqlcModel.ID,
		sqlcModel.Name,
		sqlcModel.Key,
		sqlcModel.Description,
		sqlcModel.IconUrl,
		sqlcModel.Url,
		sqlcModel.PricingPageUrl,
		nil,
		nil,
		owner,
		sqlcModel.CreatedAt,
		sqlcModel.UpdatedAt,
	)
}

// aggregateProviderData organizes rows into a structured map for easier processing
func aggregateProviderData[T any](rows []T,
	getProviderFunc func(T) sql.Provider,
	getPlanFunc func(T) sql.ProviderPlan,
	getPriceFunc func(T) sql.ProviderPrice) (map[uuid.UUID]map[uuid.UUID][]provider.Price, sql.Provider) {

	providerMap := make(map[uuid.UUID]map[uuid.UUID][]provider.Price)
	var lastProvider sql.Provider

	for _, row := range rows {
		sqlcProvider := getProviderFunc(row)
		sqlcPlan := getPlanFunc(row)
		sqlcPrice := getPriceFunc(row)

		lastProvider = sqlcProvider

		// Initialize provider map if not exists
		if providerMap[sqlcProvider.ID] == nil {
			providerMap[sqlcProvider.ID] = make(map[uuid.UUID][]provider.Price)
		}

		// Add price if valid (not null from LEFT JOIN)
		if sqlcPrice.ID != uuid.Nil {
			price := createProviderPriceFromSqlc(sqlcPrice)
			providerMap[sqlcProvider.ID][sqlcPlan.ID] = append(providerMap[sqlcProvider.ID][sqlcPlan.ID], price)
		}

		// Initialize plan entry if it doesn't exist (for plans without prices)
		if sqlcPlan.ID != uuid.Nil {
			if _, exists := providerMap[sqlcProvider.ID][sqlcPlan.ID]; !exists {
				providerMap[sqlcProvider.ID][sqlcPlan.ID] = []provider.Price{}
			}
		}
	}

	return providerMap, lastProvider
}

// buildPlansFromMap converts plan data from the aggregated map to domain plans
func buildPlansFromMap[T any](planMap map[uuid.UUID][]provider.Price, rows []T,
	getPlanFunc func(T) sql.ProviderPlan) []provider.Plan {
	if len(planMap) == 0 {
		return nil
	}

	plans := make([]provider.Plan, 0, len(planMap))
	for planID, prices := range planMap {
		// Find the plan data from rows
		var planData sql.ProviderPlan
		for _, row := range rows {
			if getPlanFunc(row).ID == planID {
				planData = getPlanFunc(row)
				break
			}
		}

		plan := provider.NewPlan(
			planData.ID,
			planData.Name,
			planData.Description,
			prices,
			planData.CreatedAt,
			planData.UpdatedAt,
		)
		plans = append(plans, plan)
	}

	return plans
}

// createProviderWithPlans creates a complete provider with its plans
func createProviderWithPlans(baseProvider sql.Provider, plans []provider.Plan) provider.Provider {
	domainProvider := createProviderFromSqlc(baseProvider)

	return provider.NewProvider(
		domainProvider.Id(),
		domainProvider.Name(),
		domainProvider.Key(),
		domainProvider.Description(),
		domainProvider.IconUrl(),
		domainProvider.Url(),
		domainProvider.PricingPageUrl(),
		nil, // labels - handled separately if needed
		plans,
		domainProvider.Owner(),
		domainProvider.CreatedAt(),
		domainProvider.UpdatedAt(),
	)
}

// createProviderFromSqlcRows converts SQLC rows to domain providers with plans and prices
func createProviderFromSqlcRows[T any](rows []T,
	getProviderFunc func(T) sql.Provider,
	getPlanFunc func(T) sql.ProviderPlan,
	getPriceFunc func(T) sql.ProviderPrice) []provider.Provider {

	if len(rows) == 0 {
		return nil
	}

	// Aggregate data from rows
	providerMap, lastProvider := aggregateProviderData(rows, getProviderFunc, getPlanFunc, getPriceFunc)

	// Convert to domain providers
	providers := make([]provider.Provider, 0, len(providerMap))
	for _, planMap := range providerMap {
		plans := buildPlansFromMap(planMap, rows, getPlanFunc)
		domainProvider := createProviderWithPlans(lastProvider, plans)
		providers = append(providers, domainProvider)
	}

	return providers
}
