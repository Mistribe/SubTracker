package models

import (
	"fmt"

	"github.com/mistribe/subtracker/internal/adapters/persistence/db/jet/app/public/model"
	"github.com/mistribe/subtracker/internal/domain/currency"
	"github.com/mistribe/subtracker/internal/domain/subscription"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/pkg/x"
	"github.com/mistribe/subtracker/pkg/x/herd"
)

type SubscriptionRow struct {
	model.Subscriptions
	SubscriptionFamilyUsers *model.SubscriptionFamilyUsers `json:"subscription_family_users"`
	SubscriptionLabels      *model.SubscriptionLabels      `json:"subscription_labels"`
	ProviderLabels          *model.ProviderLabels          `json:"provider_labels"`
}

type SubscriptionRowWithCount struct {
	SubscriptionRow
	TotalCount int64 `alias:"total_count"`
}

func createSubscriptionFromJet(
	jetModel model.Subscriptions,
	familyUsers []types.FamilyMemberID,
	subscriptionLabels []types.LabelID,
	providerLabels []types.LabelID) subscription.Subscription {
	var freeTrial subscription.FreeTrial
	if jetModel.FreeTrialEndDate != nil && jetModel.FreeTrialStartDate != nil {
		freeTrial = subscription.NewFreeTrial(
			*jetModel.FreeTrialStartDate,
			*jetModel.FreeTrialEndDate,
		)
	}

	var customPrice subscription.Price
	if jetModel.CustomPriceCurrency != nil && jetModel.CustomPriceAmount != nil {
		cry := currency.MustParseISO(*jetModel.CustomPriceCurrency)
		customPrice = subscription.NewPrice(currency.NewAmount(*jetModel.CustomPriceAmount, cry))
	}

	var payer subscription.Payer
	if jetModel.PayerType != nil && jetModel.OwnerFamilyID != nil {
		payerType := subscription.MustParsePayerType(*jetModel.PayerType)
		var memberID *types.FamilyMemberID
		if jetModel.PayerMemberID != nil {
			memberID = x.P(types.FamilyMemberID(*jetModel.PayerMemberID))
		}
		payer = subscription.NewPayer(payerType, types.FamilyID(*jetModel.OwnerFamilyID), memberID)
	}

	ownerType := types.MustParseOwnerType(jetModel.OwnerType)
	var ownerFamilyID *types.FamilyID
	var ownerUserID *types.UserID
	if jetModel.OwnerFamilyID != nil {
		ownerFamilyID = x.P(types.FamilyID(*jetModel.OwnerFamilyID))
	}
	if jetModel.OwnerUserID != nil {
		ownerUserID = x.P(types.UserID(*jetModel.OwnerUserID))
	}
	if jetModel.OwnerFamilyID == nil && jetModel.OwnerUserID == nil {
		panic("familyID or ownerUserID must be provided")
	}
	owner := types.NewOwner(ownerType, ownerFamilyID, ownerUserID)

	recurrency := subscription.MustParseRecurrencyType(jetModel.Recurrency)

	var labelRefs []subscription.LabelRef
	for _, subscriptionLabel := range subscriptionLabels {
		labelRefs = append(labelRefs, subscription.LabelRef{
			LabelId: subscriptionLabel,
			Source:  subscription.LabelSourceSubscription,
		})
	}

	for _, providerLabel := range providerLabels {
		labelRefs = append(labelRefs, subscription.LabelRef{
			LabelId: providerLabel,
			Source:  subscription.LabelSourceProvider,
		})
	}

	sub := subscription.NewSubscription(
		types.SubscriptionID(jetModel.ID),
		jetModel.FriendlyName,
		freeTrial,
		types.ProviderID(jetModel.ProviderID),
		customPrice,
		owner,
		payer,
		familyUsers,
		labelRefs,
		jetModel.StartDate,
		jetModel.EndDate,
		recurrency,
		jetModel.CustomRecurrency,
		jetModel.CreatedAt,
		jetModel.UpdatedAt,
	)
	sub.Clean()
	return sub
}

// CreateSubscriptionFromJetRows converts flattened jet rows to domain subscriptions
func CreateSubscriptionFromJetRows(rows []SubscriptionRow) []subscription.Subscription {
	if len(rows) == 0 {
		return nil
	}

	subscriptions := herd.NewDictionary[types.SubscriptionID, model.Subscriptions]()
	orderedIDs := herd.NewList[types.SubscriptionID]()
	familyUserSet := herd.NewSet[string]()
	subscriptionFamilyUsers := herd.NewDictionary[types.SubscriptionID, []types.FamilyMemberID]()
	subscriptionLabelSet := herd.NewSet[string]()
	subscriptionLabels := herd.NewDictionary[types.SubscriptionID, []types.LabelID]()
	providerLabelSet := herd.NewSet[string]()
	providerLabels := herd.NewDictionary[types.SubscriptionID, []types.LabelID]()

	for _, row := range rows {
		jetSubscription := row.Subscriptions
		subscriptionID := types.SubscriptionID(jetSubscription.ID)
		if _, ok := subscriptions[subscriptionID]; !ok {
			subscriptions[subscriptionID] = jetSubscription
			orderedIDs.Add(subscriptionID)
		}
		if row.SubscriptionFamilyUsers != nil {
			key := fmt.Sprintf("%s:%s", subscriptionID, *row.SubscriptionFamilyUsers)
			if !familyUserSet.Contains(key) {
				familyUserSet.Add(key)
				subscriptionFamilyUsers[subscriptionID] = append(subscriptionFamilyUsers[subscriptionID],
					types.FamilyMemberID(row.SubscriptionFamilyUsers.FamilyMemberID))
			}
		}
		if row.SubscriptionLabels != nil {
			key := fmt.Sprintf("%s:%s", subscriptionID, *row.SubscriptionLabels)
			if !subscriptionLabelSet.Contains(key) {
				subscriptionLabelSet.Add(key)
				subscriptionLabels[subscriptionID] = append(subscriptionLabels[subscriptionID],
					types.LabelID(row.SubscriptionLabels.LabelID))
			}
		}
		if row.ProviderLabels != nil {
			key := fmt.Sprintf("%s:%s", subscriptionID, *row.ProviderLabels)
			if !providerLabelSet.Contains(key) {
				providerLabelSet.Add(key)
				providerLabels[subscriptionID] = append(providerLabels[subscriptionID],
					types.LabelID(row.ProviderLabels.LabelID))
			}
		}
	}

	results := make([]subscription.Subscription, 0, orderedIDs.Len())
	for subscriptionId := range orderedIDs.It() {
		jetSubscription := subscriptions[subscriptionId]
		jetSubscriptionServiceUsers, _ := subscriptionFamilyUsers[subscriptionId]
		jetSubscriptionLabels, _ := subscriptionLabels[subscriptionId]
		jetProviderLabels, _ := providerLabels[subscriptionId]
		result := createSubscriptionFromJet(jetSubscription,
			jetSubscriptionServiceUsers,
			jetSubscriptionLabels,
			jetProviderLabels)
		results = append(results, result)
	}

	return results
}

func CreateSubscriptionFromJetRowsWithCount(rows []SubscriptionRowWithCount) []subscription.Subscription {
	simpleRows := herd.Select(rows, func(row SubscriptionRowWithCount) SubscriptionRow {
		return row.SubscriptionRow
	})

	return CreateSubscriptionFromJetRows(simpleRows)
}
