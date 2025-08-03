package persistence

import (
	"github.com/google/uuid"
	"golang.org/x/text/currency"

	"github.com/oleexo/subtracker/internal/domain/auth"

	"github.com/oleexo/subtracker/internal/domain/subscription"
	"github.com/oleexo/subtracker/internal/infrastructure/persistence/sql"
)

func createSubscriptionFromSqlc(sqlcModel sql.Subscription) subscription.Subscription {
	var freeTrial subscription.FreeTrial
	if sqlcModel.FreeTrialEndDate != nil && sqlcModel.FreeTrialStartDate != nil {
		freeTrial = subscription.NewFreeTrial(
			*sqlcModel.FreeTrialStartDate,
			*sqlcModel.FreeTrialEndDate,
		)
	}

	var customPrice subscription.CustomPrice
	if sqlcModel.CustomPriceCurrency != nil && sqlcModel.CustomPriceAmount != nil {
		cry := currency.MustParseISO(*sqlcModel.CustomPriceCurrency)
		customPrice = subscription.NewCustomPrice(
			*sqlcModel.CustomPriceAmount,
			cry,
		)
	}

	var payer subscription.Payer
	if sqlcModel.PayerType != nil && sqlcModel.FamilyID != nil {
		payerType := subscription.MustParsePayerType(*sqlcModel.PayerType)
		payer = subscription.NewPayer(payerType, *sqlcModel.FamilyID, sqlcModel.PayerMemberID)
	}

	ownerType := auth.MustParseOwnerType(sqlcModel.OwnerType)
	owner := auth.NewOwner(ownerType, sqlcModel.OwnerFamilyID, sqlcModel.OwnerUserID)

	recurrency := subscription.MustParseRecurrencyType(sqlcModel.Recurrency)

	sub := subscription.NewSubscription(
		sqlcModel.ID,
		sqlcModel.FriendlyName,
		freeTrial,
		sqlcModel.ProviderID,
		sqlcModel.PlanID,
		sqlcModel.PriceID,
		customPrice,
		owner,
		payer,
		nil,
		sqlcModel.StartDate,
		sqlcModel.EndDate,
		recurrency,
		sqlcModel.CustomRecurrency,
		sqlcModel.CreatedAt,
		sqlcModel.UpdatedAt,
	)
	sub.Clean()
	return sub
}

func createSubscriptionFromSqlcRows[T any](rows []T,
	getSubscriptionFunc func(T) sql.Subscription,
	getSubscriptionServiceUserFunc func(T) sql.SubscriptionServiceUser) []subscription.Subscription {
	if len(rows) == 0 {
		return nil
	}

	subscriptionMap := make(map[uuid.UUID][]uuid.UUID)
	var lastSubscription sql.Subscription

	for _, row := range rows {
		sqlcSubscription := getSubscriptionFunc(row)
		sqlcServiceUser := getSubscriptionServiceUserFunc(row)

		lastSubscription = sqlcSubscription // Keep reference for subscription data

		// Add service user if valid (not null from LEFT JOIN)
		if sqlcServiceUser.SubscriptionID != uuid.Nil {
			subscriptionMap[sqlcSubscription.ID] = append(subscriptionMap[sqlcSubscription.ID],
				sqlcServiceUser.FamilyMemberID)
		}
	}

	// Convert to domain subscriptions
	subscriptions := make([]subscription.Subscription, 0, len(subscriptionMap))
	for id, serviceUsers := range subscriptionMap {
		if id == lastSubscription.ID {
			sub := createSubscriptionFromSqlc(lastSubscription)
			sub.SetServiceUsers(serviceUsers)
			sub.Clean()
			subscriptions = append(subscriptions, sub)
		}
	}

	return subscriptions
}
