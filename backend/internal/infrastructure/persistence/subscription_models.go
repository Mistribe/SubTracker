package persistence

import (
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/text/currency"

	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/internal/domain/subscription"
	"github.com/oleexo/subtracker/internal/infrastructure/persistence/jet/app/public/model"
	"github.com/oleexo/subtracker/pkg/slicesx"
)

func createSubscriptionFromJet(jetModel model.Subscriptions, serviceUsers []uuid.UUID) subscription.Subscription {
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

func createSubscriptionFromJetRows(rows []struct {
	Subscriptions            model.Subscriptions             `json:"subscriptions"`
	SubscriptionServiceUsers *model.SubscriptionServiceUsers `json:"subscription_service_users"`
}) []subscription.Subscription {
	if len(rows) == 0 {
		return nil
	}

	subscriptions := make(map[uuid.UUID]model.Subscriptions)
	orderedIDs := make([]uuid.UUID, 0, len(rows))
	serviceUserSet := slicesx.NewSet[string]()
	subscriptionServiceUsers := make(map[uuid.UUID][]uuid.UUID)

	for _, row := range rows {
		jetSubscription := row.Subscriptions
		if _, ok := subscriptions[jetSubscription.ID]; !ok {
			subscriptions[jetSubscription.ID] = jetSubscription
			orderedIDs = append(orderedIDs, jetSubscription.ID)
		}
		if row.SubscriptionServiceUsers != nil && row.SubscriptionServiceUsers.FamilyMemberID != uuid.Nil {
			key := fmt.Sprintf("%s:%s", jetSubscription.ID, row.SubscriptionServiceUsers.FamilyMemberID)
			if !serviceUserSet.Contains(key) {
				serviceUserSet.Add(key)
				subscriptionServiceUsers[jetSubscription.ID] = append(subscriptionServiceUsers[jetSubscription.ID],
					row.SubscriptionServiceUsers.FamilyMemberID)
			}
		}
	}

	results := make([]subscription.Subscription, 0, len(orderedIDs))
	for _, subscriptionId := range orderedIDs {
		jetSubscription := subscriptions[subscriptionId]
		jetSubscriptionServiceUsers, _ := subscriptionServiceUsers[subscriptionId]
		results = append(results, createSubscriptionFromJet(jetSubscription, jetSubscriptionServiceUsers))
	}

	return results
}

func createSubscriptionFromJetRowsWithCount(rows []struct {
	Subscriptions            model.Subscriptions             `json:"subscriptions"`
	SubscriptionServiceUsers *model.SubscriptionServiceUsers `json:"subscription_service_users"`
	TotalCount               int64                           `json:"total_count"`
}) []subscription.Subscription {
	// Convert to the simpler row structure
	simpleRows := slicesx.Select(rows, func(row struct {
		Subscriptions            model.Subscriptions             `json:"subscriptions"`
		SubscriptionServiceUsers *model.SubscriptionServiceUsers `json:"subscription_service_users"`
		TotalCount               int64                           `json:"total_count"`
	}) struct {
		Subscriptions            model.Subscriptions             `json:"subscriptions"`
		SubscriptionServiceUsers *model.SubscriptionServiceUsers `json:"subscription_service_users"`
	} {
		return struct {
			Subscriptions            model.Subscriptions             `json:"subscriptions"`
			SubscriptionServiceUsers *model.SubscriptionServiceUsers `json:"subscription_service_users"`
		}{
			Subscriptions:            row.Subscriptions,
			SubscriptionServiceUsers: row.SubscriptionServiceUsers,
		}
	})

	return createSubscriptionFromJetRows(simpleRows)
}

// createSubscriptionFromJetRowsFlat converts flattened jet rows to domain subscriptions
func createSubscriptionFromJetRowsFlat(rows []struct {
	model.Subscriptions
	FamilyMemberID *uuid.UUID `alias:"subscription_service_users.family_member_id"`
	SubscriptionID *uuid.UUID `alias:"subscription_service_users.subscription_id"`
}) []subscription.Subscription {
	if len(rows) == 0 {
		return nil
	}

	subscriptions := make(map[uuid.UUID]model.Subscriptions)
	orderedIDs := make([]uuid.UUID, 0, len(rows))
	serviceUserSet := slicesx.NewSet[string]()
	subscriptionServiceUsers := make(map[uuid.UUID][]uuid.UUID)

	for _, row := range rows {
		jetSubscription := row.Subscriptions
		if _, ok := subscriptions[jetSubscription.ID]; !ok {
			subscriptions[jetSubscription.ID] = jetSubscription
			orderedIDs = append(orderedIDs, jetSubscription.ID)
		}
		if row.FamilyMemberID != nil && *row.FamilyMemberID != uuid.Nil {
			key := fmt.Sprintf("%s:%s", jetSubscription.ID, *row.FamilyMemberID)
			if !serviceUserSet.Contains(key) {
				serviceUserSet.Add(key)
				subscriptionServiceUsers[jetSubscription.ID] = append(subscriptionServiceUsers[jetSubscription.ID],
					*row.FamilyMemberID)
			}
		}
	}

	results := make([]subscription.Subscription, 0, len(orderedIDs))
	for _, subscriptionId := range orderedIDs {
		jetSubscription := subscriptions[subscriptionId]
		jetSubscriptionServiceUsers, _ := subscriptionServiceUsers[subscriptionId]
		results = append(results, createSubscriptionFromJet(jetSubscription, jetSubscriptionServiceUsers))
	}

	return results
}

// createSubscriptionFromJetRowsFlatWithCount converts flattened jet rows with count to domain subscriptions
func createSubscriptionFromJetRowsFlatWithCount(rows []struct {
	model.Subscriptions
	FamilyMemberID *uuid.UUID `alias:"subscription_service_users.family_member_id"`
	SubscriptionID *uuid.UUID `alias:"subscription_service_users.subscription_id"`
	TotalCount     int64      `alias:"total_count"`
	IsActive       bool       `alias:"is_active"`
}) []subscription.Subscription {
	// Convert to the simpler row structure
	simpleRows := slicesx.Select(rows, func(row struct {
		model.Subscriptions
		FamilyMemberID *uuid.UUID `alias:"subscription_service_users.family_member_id"`
		SubscriptionID *uuid.UUID `alias:"subscription_service_users.subscription_id"`
		TotalCount     int64      `alias:"total_count"`
		IsActive       bool       `alias:"is_active"`
	}) struct {
		model.Subscriptions
		FamilyMemberID *uuid.UUID `alias:"subscription_service_users.family_member_id"`
		SubscriptionID *uuid.UUID `alias:"subscription_service_users.subscription_id"`
	} {
		return struct {
			model.Subscriptions
			FamilyMemberID *uuid.UUID `alias:"subscription_service_users.family_member_id"`
			SubscriptionID *uuid.UUID `alias:"subscription_service_users.subscription_id"`
		}{
			Subscriptions:  row.Subscriptions,
			FamilyMemberID: row.FamilyMemberID,
			SubscriptionID: row.SubscriptionID,
		}
	})

	return createSubscriptionFromJetRowsFlat(simpleRows)
}
