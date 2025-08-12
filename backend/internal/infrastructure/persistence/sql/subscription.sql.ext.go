package sql

import (
	"context"

	"github.com/google/uuid"
)

type SubscriptionRow struct {
	Subscription Subscription
	ServiceUser  *SubscriptionServiceUser
}

func (q *Queries) GetSubscriptionsForUser(ctx context.Context, userId string, limit, offset int32) (
	[]SubscriptionRow,
	int64,
	error) {
	rows, err := q.getSubscriptionsForUser(ctx, getSubscriptionsForUserParams{
		UserID: &userId,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, 0, err
	}
	if len(rows) == 0 {
		return nil, 0, nil
	}
	totalCount := rows[0].TotalCount
	results := make([]SubscriptionRow, len(rows))
	for i, row := range rows {
		var serviceUser *SubscriptionServiceUser
		if row.SubscriptionServiceUsersFamilyMemberID != nil {
			serviceUser = &SubscriptionServiceUser{
				FamilyMemberID: *row.SubscriptionServiceUsersFamilyMemberID,
				SubscriptionID: row.SubscriptionsID,
			}
		}
		results[i] = SubscriptionRow{
			Subscription: Subscription{
				ID:                  row.SubscriptionsID,
				OwnerType:           row.SubscriptionsOwnerType,
				OwnerFamilyID:       row.SubscriptionsOwnerFamilyID,
				OwnerUserID:         row.SubscriptionsOwnerUserID,
				FriendlyName:        row.SubscriptionsFriendlyName,
				FreeTrialStartDate:  row.SubscriptionsFreeTrialStartDate,
				FreeTrialEndDate:    row.SubscriptionsFreeTrialEndDate,
				ProviderID:          row.SubscriptionsProviderID,
				PlanID:              row.SubscriptionsPlanID,
				PriceID:             row.SubscriptionsPriceID,
				FamilyID:            row.SubscriptionsFamilyID,
				PayerType:           row.SubscriptionsPayerType,
				PayerMemberID:       row.SubscriptionsPayerMemberID,
				StartDate:           row.SubscriptionsStartDate,
				EndDate:             row.SubscriptionsEndDate,
				Recurrency:          row.SubscriptionsRecurrency,
				CustomRecurrency:    row.SubscriptionsCustomRecurrency,
				CustomPriceCurrency: row.SubscriptionsCustomPriceCurrency,
				CustomPriceAmount:   row.SubscriptionsCustomPriceAmount,
				CreatedAt:           row.SubscriptionsCreatedAt,
				UpdatedAt:           row.SubscriptionsUpdatedAt,
				Etag:                row.SubscriptionsEtag,
			},
			ServiceUser: serviceUser,
		}
	}

	return results, totalCount, nil
}

func (q *Queries) GetSubscriptions(ctx context.Context, limit, offset int32) (
	[]SubscriptionRow,
	int64,
	error) {
	rows, err := q.getSubscriptions(ctx, getSubscriptionsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, 0, err
	}
	if len(rows) == 0 {
		return nil, 0, nil
	}
	totalCount := rows[0].TotalCount
	results := make([]SubscriptionRow, len(rows))
	for i, row := range rows {
		var serviceUser *SubscriptionServiceUser
		if row.SubscriptionServiceUsersFamilyMemberID != nil {
			serviceUser = &SubscriptionServiceUser{
				FamilyMemberID: *row.SubscriptionServiceUsersFamilyMemberID,
				SubscriptionID: row.SubscriptionsID,
			}
		}
		results[i] = SubscriptionRow{
			Subscription: Subscription{
				ID:                  row.SubscriptionsID,
				OwnerType:           row.SubscriptionsOwnerType,
				OwnerFamilyID:       row.SubscriptionsOwnerFamilyID,
				OwnerUserID:         row.SubscriptionsOwnerUserID,
				FriendlyName:        row.SubscriptionsFriendlyName,
				FreeTrialStartDate:  row.SubscriptionsFreeTrialStartDate,
				FreeTrialEndDate:    row.SubscriptionsFreeTrialEndDate,
				ProviderID:          row.SubscriptionsProviderID,
				PlanID:              row.SubscriptionsPlanID,
				PriceID:             row.SubscriptionsPriceID,
				FamilyID:            row.SubscriptionsFamilyID,
				PayerType:           row.SubscriptionsPayerType,
				PayerMemberID:       row.SubscriptionsPayerMemberID,
				StartDate:           row.SubscriptionsStartDate,
				EndDate:             row.SubscriptionsEndDate,
				Recurrency:          row.SubscriptionsRecurrency,
				CustomRecurrency:    row.SubscriptionsCustomRecurrency,
				CustomPriceCurrency: row.SubscriptionsCustomPriceCurrency,
				CustomPriceAmount:   row.SubscriptionsCustomPriceAmount,
				CreatedAt:           row.SubscriptionsCreatedAt,
				UpdatedAt:           row.SubscriptionsUpdatedAt,
				Etag:                row.SubscriptionsEtag,
			},
			ServiceUser: serviceUser,
		}
	}

	return results, totalCount, nil
}

func (q *Queries) GetSubscriptionById(ctx context.Context, id uuid.UUID) ([]SubscriptionRow, error) {
	rows, err := q.getSubscriptionById(ctx, id)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, nil
	}

	results := make([]SubscriptionRow, len(rows))
	for i, row := range rows {
		var serviceUser *SubscriptionServiceUser
		if row.SubscriptionServiceUsersFamilyMemberID != nil {
			serviceUser = &SubscriptionServiceUser{
				FamilyMemberID: *row.SubscriptionServiceUsersFamilyMemberID,
				SubscriptionID: row.SubscriptionsID,
			}
		}
		results[i] = SubscriptionRow{
			Subscription: Subscription{
				ID:                  row.SubscriptionsID,
				OwnerType:           row.SubscriptionsOwnerType,
				OwnerFamilyID:       row.SubscriptionsOwnerFamilyID,
				OwnerUserID:         row.SubscriptionsOwnerUserID,
				FriendlyName:        row.SubscriptionsFriendlyName,
				FreeTrialStartDate:  row.SubscriptionsFreeTrialStartDate,
				FreeTrialEndDate:    row.SubscriptionsFreeTrialEndDate,
				ProviderID:          row.SubscriptionsProviderID,
				PlanID:              row.SubscriptionsPlanID,
				PriceID:             row.SubscriptionsPriceID,
				FamilyID:            row.SubscriptionsFamilyID,
				PayerType:           row.SubscriptionsPayerType,
				PayerMemberID:       row.SubscriptionsPayerMemberID,
				StartDate:           row.SubscriptionsStartDate,
				EndDate:             row.SubscriptionsEndDate,
				Recurrency:          row.SubscriptionsRecurrency,
				CustomRecurrency:    row.SubscriptionsCustomRecurrency,
				CustomPriceCurrency: row.SubscriptionsCustomPriceCurrency,
				CustomPriceAmount:   row.SubscriptionsCustomPriceAmount,
				CreatedAt:           row.SubscriptionsCreatedAt,
				UpdatedAt:           row.SubscriptionsUpdatedAt,
				Etag:                row.SubscriptionsEtag,
			},
			ServiceUser: serviceUser,
		}
	}
	return results, nil
}

func (q *Queries) GetSubscriptionByIdForUser(ctx context.Context, id uuid.UUID, userId string) (
	[]SubscriptionRow,
	error) {
	rows, err := q.getSubscriptionByIdForUser(ctx, getSubscriptionByIdForUserParams{
		UserID: &userId,
		ID:     id,
	})
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, nil
	}

	results := make([]SubscriptionRow, len(rows))
	for i, row := range rows {
		var serviceUser *SubscriptionServiceUser
		if row.SubscriptionServiceUsersFamilyMemberID != nil {
			serviceUser = &SubscriptionServiceUser{
				FamilyMemberID: *row.SubscriptionServiceUsersFamilyMemberID,
				SubscriptionID: row.SubscriptionsID,
			}
		}
		results[i] = SubscriptionRow{
			Subscription: Subscription{
				ID:                  row.SubscriptionsID,
				OwnerType:           row.SubscriptionsOwnerType,
				OwnerFamilyID:       row.SubscriptionsOwnerFamilyID,
				OwnerUserID:         row.SubscriptionsOwnerUserID,
				FriendlyName:        row.SubscriptionsFriendlyName,
				FreeTrialStartDate:  row.SubscriptionsFreeTrialStartDate,
				FreeTrialEndDate:    row.SubscriptionsFreeTrialEndDate,
				ProviderID:          row.SubscriptionsProviderID,
				PlanID:              row.SubscriptionsPlanID,
				PriceID:             row.SubscriptionsPriceID,
				FamilyID:            row.SubscriptionsFamilyID,
				PayerType:           row.SubscriptionsPayerType,
				PayerMemberID:       row.SubscriptionsPayerMemberID,
				StartDate:           row.SubscriptionsStartDate,
				EndDate:             row.SubscriptionsEndDate,
				Recurrency:          row.SubscriptionsRecurrency,
				CustomRecurrency:    row.SubscriptionsCustomRecurrency,
				CustomPriceCurrency: row.SubscriptionsCustomPriceCurrency,
				CustomPriceAmount:   row.SubscriptionsCustomPriceAmount,
				CreatedAt:           row.SubscriptionsCreatedAt,
				UpdatedAt:           row.SubscriptionsUpdatedAt,
				Etag:                row.SubscriptionsEtag,
			},
			ServiceUser: serviceUser,
		}
	}
	return results, nil
}
