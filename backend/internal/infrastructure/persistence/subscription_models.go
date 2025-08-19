package persistence

import (
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/text/currency"

	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/pkg/slicesx"

	"github.com/oleexo/subtracker/internal/domain/subscription"
	"github.com/oleexo/subtracker/internal/infrastructure/persistence/sql"
)

func createSubscriptionFromSqlc(sqlcModel sql.Subscription, serviceUsers []uuid.UUID) subscription.Subscription {
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
		serviceUsers,
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
	getSubscriptionServiceUserFunc func(T) *sql.SubscriptionServiceUser) []subscription.Subscription {
	if len(rows) == 0 {
		return nil
	}

	subscriptions := make(map[uuid.UUID]sql.Subscription)
	orderedIDs := make([]uuid.UUID, 0, len(rows))
	serviceUserSet := slicesx.NewSet[string]()
	subscriptionServiceUsers := make(map[uuid.UUID][]uuid.UUID)

	for _, row := range rows {
		sqlcSubscription := getSubscriptionFunc(row)
		if _, ok := subscriptions[sqlcSubscription.ID]; !ok {
			subscriptions[sqlcSubscription.ID] = sqlcSubscription
			orderedIDs = append(orderedIDs, sqlcSubscription.ID)
		}
		sqlcSubscriptionServiceUser := getSubscriptionServiceUserFunc(row)
		if sqlcSubscriptionServiceUser != nil {
			key := fmt.Sprintf("%s:%s", sqlcSubscription.ID, sqlcSubscriptionServiceUser.FamilyMemberID)
			if !serviceUserSet.Contains(key) {
				serviceUserSet.Add(key)
				subscriptionServiceUsers[sqlcSubscription.ID] = append(subscriptionServiceUsers[sqlcSubscription.ID],
					sqlcSubscriptionServiceUser.FamilyMemberID)
			}
		}
	}

	results := make([]subscription.Subscription, 0, len(orderedIDs))
	for _, subscriptionId := range orderedIDs {
		sqlcSubscription := subscriptions[subscriptionId]
		sqlcSubscriptionServiceUsers, _ := subscriptionServiceUsers[subscriptionId]
		results = append(results, createSubscriptionFromSqlc(sqlcSubscription, sqlcSubscriptionServiceUsers))
	}

	return results
}
