package models

import (
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/text/currency"

	"github.com/mistribe/subtracker/internal/adapters/persistence/db/jet/app/public/model"
	"github.com/mistribe/subtracker/internal/domain/auth"
	"github.com/mistribe/subtracker/internal/domain/subscription"
	"github.com/mistribe/subtracker/pkg/slicesx"
)

type SubscriptionRow struct {
	model.Subscriptions
	SubscriptionServiceUsers *model.SubscriptionServiceUsers `json:"subscription_service_users"`
	SubscriptionLabels       *model.SubscriptionLabels       `json:"subscription_labels"`
	ProviderLabels           *model.ProviderLabels           `json:"provider_labels"`
}

type SubscriptionRowWithCount struct {
	SubscriptionRow
	TotalCount int64 `alias:"total_count"`
}

func createSubscriptionFromJet(
	jetModel model.Subscriptions,
	serviceUsers []uuid.UUID,
	subscriptionLabels []uuid.UUID,
	providerLabels []uuid.UUID) subscription.Subscription {
	var freeTrial subscription.FreeTrial
	if jetModel.FreeTrialEndDate != nil && jetModel.FreeTrialStartDate != nil {
		freeTrial = subscription.NewFreeTrial(
			*jetModel.FreeTrialStartDate,
			*jetModel.FreeTrialEndDate,
		)
	}

	var customPrice subscription.CustomPrice
	if jetModel.CustomPriceCurrency != nil && jetModel.CustomPriceAmount != nil {
		cry := currency.MustParseISO(*jetModel.CustomPriceCurrency)
		customPrice = subscription.NewCustomPrice(
			*jetModel.CustomPriceAmount,
			cry,
		)
	}

	var payer subscription.Payer
	if jetModel.PayerType != nil && jetModel.FamilyID != nil {
		payerType := subscription.MustParsePayerType(*jetModel.PayerType)
		payer = subscription.NewPayer(payerType, *jetModel.FamilyID, jetModel.PayerMemberID)
	}

	ownerType := auth.MustParseOwnerType(jetModel.OwnerType)
	owner := auth.NewOwner(ownerType, jetModel.OwnerFamilyID, jetModel.OwnerUserID)

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
		jetModel.ID,
		jetModel.FriendlyName,
		freeTrial,
		jetModel.ProviderID,
		jetModel.PlanID,
		jetModel.PriceID,
		customPrice,
		owner,
		payer,
		serviceUsers,
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

	subscriptions := make(map[uuid.UUID]model.Subscriptions)
	orderedIDs := make([]uuid.UUID, 0, len(rows))
	serviceUserSet := slicesx.NewSet[string]()
	subscriptionServiceUsers := make(map[uuid.UUID][]uuid.UUID)
	subscriptionLabelSet := slicesx.NewSet[string]()
	subscriptionLabels := make(map[uuid.UUID][]uuid.UUID)
	providerLabelSet := slicesx.NewSet[string]()
	providerLabels := make(map[uuid.UUID][]uuid.UUID)

	for _, row := range rows {
		jetSubscription := row.Subscriptions
		if _, ok := subscriptions[jetSubscription.ID]; !ok {
			subscriptions[jetSubscription.ID] = jetSubscription
			orderedIDs = append(orderedIDs, jetSubscription.ID)
		}
		if row.SubscriptionServiceUsers != nil {
			key := fmt.Sprintf("%s:%s", jetSubscription.ID, *row.SubscriptionServiceUsers)
			if !serviceUserSet.Contains(key) {
				serviceUserSet.Add(key)
				subscriptionServiceUsers[jetSubscription.ID] = append(subscriptionServiceUsers[jetSubscription.ID],
					row.SubscriptionServiceUsers.FamilyMemberID)
			}
		}
		if row.SubscriptionLabels != nil {
			key := fmt.Sprintf("%s:%s", jetSubscription.ID, *row.SubscriptionLabels)
			if !subscriptionLabelSet.Contains(key) {
				subscriptionLabelSet.Add(key)
				subscriptionLabels[jetSubscription.ID] = append(subscriptionLabels[jetSubscription.ID],
					row.SubscriptionLabels.LabelID)
			}
		}
		if row.ProviderLabels != nil {
			key := fmt.Sprintf("%s:%s", jetSubscription.ID, *row.ProviderLabels)
			if !providerLabelSet.Contains(key) {
				providerLabelSet.Add(key)
				providerLabels[jetSubscription.ID] = append(providerLabels[jetSubscription.ID],
					row.ProviderLabels.LabelID)
			}
		}
	}

	results := make([]subscription.Subscription, 0, len(orderedIDs))
	for _, subscriptionId := range orderedIDs {
		jetSubscription := subscriptions[subscriptionId]
		jetSubscriptionServiceUsers, _ := subscriptionServiceUsers[subscriptionId]
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
	simpleRows := slicesx.Select(rows, func(row SubscriptionRowWithCount) SubscriptionRow {
		return row.SubscriptionRow
	})

	return CreateSubscriptionFromJetRows(simpleRows)
}
