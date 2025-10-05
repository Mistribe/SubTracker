package models

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/adapters/persistence/db/jet/app/public/model"
	"github.com/mistribe/subtracker/internal/domain/provider"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/pkg/x"
	"github.com/mistribe/subtracker/pkg/x/herd"
)

func createProviderFromJet(jetModel model.Providers, labels []types.LabelID) provider.Provider {
	ownerType := types.MustParseOwnerType(jetModel.OwnerType)
	var ownerFamilyID *types.FamilyID
	var ownerUserID *types.UserID
	if jetModel.OwnerFamilyID != nil {
		ownerFamilyID = x.P(types.FamilyID(*jetModel.OwnerFamilyID))
	}
	if jetModel.OwnerUserID != nil {
		ownerUserID = x.P(types.UserID(*jetModel.OwnerUserID))
	}
	owner := types.NewOwner(ownerType, ownerFamilyID, ownerUserID)

	p := provider.NewProvider(
		types.ProviderID(jetModel.ID),
		jetModel.Name,
		jetModel.Key,
		jetModel.Description,
		jetModel.IconURL,
		jetModel.URL,
		jetModel.PricingPageURL,
		labels,
		owner,
		jetModel.CreatedAt,
		jetModel.UpdatedAt,
	)

	p.Clean()
	return p
}

type ProviderRow struct {
	Providers      model.Providers       `json:"providers"`
	ProviderLabels *model.ProviderLabels `json:"provider_labels"`
}

type ProviderRowWithCount struct {
	ProviderRow
	TotalCount int64 `alias:"total_count" json:"total_count"`
}

// createProviderFromJetRows converts jet rows to domain providers with plans and prices
func CreateProviderFromJetRows(rows []ProviderRow) []provider.Provider {
	if len(rows) == 0 {
		return nil
	}

	providers := herd.NewDictionary[types.ProviderID, model.Providers]()
	labelSet := herd.NewSet[string]()
	providerLabels := herd.NewDictionary[types.ProviderID, []types.LabelID]()

	for _, row := range rows {
		jetProvider := row.Providers
		providerID := types.ProviderID(jetProvider.ID)
		if _, ok := providers[providerID]; !ok {
			providers[providerID] = jetProvider
		}

		if row.ProviderLabels != nil && row.ProviderLabels.LabelID != uuid.Nil {
			key := fmt.Sprintf("%s:%s", jetProvider.ID, row.ProviderLabels.LabelID)
			if _, ok := labelSet[key]; !ok {
				labelSet[key] = struct{}{}
				providerLabels[providerID] = append(providerLabels[providerID],
					types.LabelID(row.ProviderLabels.LabelID))
			}
		}
	}

	results := make([]provider.Provider, len(providers))
	count := 0
	for providerId, jetProvider := range providers {

		jetLabels, _ := providerLabels[providerId]
		results[count] = createProviderFromJet(jetProvider, jetLabels)
		count++
	}

	return results
}

func CreateProviderFromJetRowsWithCount(rows []ProviderRowWithCount) []provider.Provider {
	simpleRows := herd.Select(rows, func(row ProviderRowWithCount) ProviderRow {
		return ProviderRow{
			Providers:      row.Providers,
			ProviderLabels: row.ProviderLabels,
		}
	})

	return CreateProviderFromJetRows(simpleRows)
}
