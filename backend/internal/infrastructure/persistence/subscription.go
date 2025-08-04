package persistence

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/internal/domain/entity"
	"github.com/oleexo/subtracker/internal/domain/subscription"
	"github.com/oleexo/subtracker/internal/infrastructure/persistence/sql"
	"github.com/oleexo/subtracker/pkg/slicesx"
)

type SubscriptionRepository struct {
	dbContext *DatabaseContext
}

func NewSubscriptionRepository(repository *DatabaseContext) subscription.Repository {
	return &SubscriptionRepository{
		dbContext: repository,
	}
}

func (r SubscriptionRepository) GetById(ctx context.Context, id uuid.UUID) (subscription.Subscription, error) {
	response, err := r.dbContext.GetQueries(ctx).GetSubscriptionById(ctx, id)
	if err != nil {
		return nil, err
	}

	if len(response) == 0 {
		return nil, nil
	}

	subscriptions := createSubscriptionFromSqlcRows(response,
		func(row sql.SubscriptionRow) sql.Subscription {
			return row.Subscription
		},
		func(row sql.SubscriptionRow) *sql.SubscriptionServiceUser {
			return row.ServiceUser
		},
	)
	if len(subscriptions) == 0 {
		return nil, nil
	}

	return subscriptions[0], nil
}

func (r SubscriptionRepository) GetAll(ctx context.Context,
	parameters entity.QueryParameters) ([]subscription.Subscription, int64, error) {
	response, count, err := r.dbContext.GetQueries(ctx).GetSubscriptions(ctx, parameters.Limit, parameters.Offset)
	if err != nil {
		return nil, 0, err
	}

	if len(response) == 0 {
		return nil, 0, nil
	}

	subscriptions := createSubscriptionFromSqlcRows(response,
		func(row sql.SubscriptionRow) sql.Subscription {
			return row.Subscription
		},
		func(row sql.SubscriptionRow) *sql.SubscriptionServiceUser {
			return row.ServiceUser
		},
	)
	if len(subscriptions) == 0 {
		return nil, 0, nil
	}

	return subscriptions, count, nil
}

func (r SubscriptionRepository) Save(ctx context.Context, subscriptions ...subscription.Subscription) error {
	var newSubscriptions []subscription.Subscription
	for _, sub := range subscriptions {
		if !sub.IsExists() {
			newSubscriptions = append(newSubscriptions, sub)
		} else {
			if err := r.update(ctx, sub); err != nil {
				return err
			}
		}
	}

	if len(newSubscriptions) > 0 {
		if err := r.create(ctx, newSubscriptions); err != nil {
			return err
		}
	}

	for _, sub := range subscriptions {
		sub.Clean()
	}

	return nil
}

func (r SubscriptionRepository) Delete(ctx context.Context, subscriptionId uuid.UUID) (bool, error) {
	err := r.dbContext.GetQueries(ctx).DeleteSubscription(ctx, subscriptionId)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r SubscriptionRepository) Exists(ctx context.Context, ids ...uuid.UUID) (bool, error) {
	count, err := r.dbContext.GetQueries(ctx).IsSubscriptionExists(ctx, ids)
	if err != nil {
		return false, err
	}
	return count > int64(len(ids)), nil
}

func (r SubscriptionRepository) create(ctx context.Context, subscriptions []subscription.Subscription) error {
	if len(subscriptions) == 0 {
		return nil
	}

	subscriptionArgs := slicesx.Select(subscriptions,
		func(sub subscription.Subscription) sql.CreateSubscriptionsParams {
			args := sql.CreateSubscriptionsParams{
				ID:               sub.Id(),
				FriendlyName:     sub.FriendlyName(),
				ProviderID:       sub.ProviderId(),
				PlanID:           sub.PlanId(),
				PriceID:          sub.PriceId(),
				OwnerType:        sub.Owner().Type().String(),
				StartDate:        sub.StartDate(),
				EndDate:          sub.EndDate(),
				Recurrency:       sub.Recurrency().String(),
				CustomRecurrency: sub.CustomRecurrency(),
				CreatedAt:        sub.CreatedAt(),
				UpdatedAt:        sub.UpdatedAt(),
				Etag:             sub.ETag(),
			}

			if sub.FreeTrial() != nil {
				freeTrialStartDate := sub.FreeTrial().StartDate()
				args.FreeTrialStartDate = &freeTrialStartDate
				freeTrialEndDate := sub.FreeTrial().EndDate()
				args.FreeTrialEndDate = &freeTrialEndDate
			}

			if sub.CustomPrice() != nil {
				customPriceAmount := sub.CustomPrice().Amount()
				args.CustomPriceAmount = &customPriceAmount
				customPriceCurrency := sub.CustomPrice().Currency().String()
				args.CustomPriceCurrency = &customPriceCurrency
			}

			switch sub.Owner().Type() {
			case auth.PersonalOwnerType:
				userId := sub.Owner().UserId()
				args.OwnerUserID = &userId
			case auth.FamilyOwnerType:
				familyId := sub.Owner().FamilyId()
				args.OwnerFamilyID = &familyId
			}

			if sub.Payer() != nil {
				payerType := sub.Payer().Type().String()
				args.PayerType = &payerType
				familyId := sub.Payer().FamilyId()
				args.FamilyID = &familyId
				switch sub.Payer().Type() {
				case subscription.FamilyMemberPayer:
					payerMemberId := sub.Payer().MemberId()
					args.PayerMemberID = &payerMemberId
				}
			}

			return args
		})

	if _, err := r.dbContext.GetQueries(ctx).CreateSubscriptions(ctx, subscriptionArgs); err != nil {
		return err
	}

	serverUserArgs := slicesx.SelectMany(subscriptions,
		func(sub subscription.Subscription) []sql.CreateSubscriptionServiceUsersParams {
			return slicesx.Select(sub.ServiceUsers().Values(),
				func(u uuid.UUID) sql.CreateSubscriptionServiceUsersParams {
					return sql.CreateSubscriptionServiceUsersParams{
						SubscriptionID: sub.Id(),
						FamilyMemberID: u,
					}
				})
		})

	if len(serverUserArgs) > 0 {
		if _, err := r.dbContext.GetQueries(ctx).CreateSubscriptionServiceUsers(ctx, serverUserArgs); err != nil {
			return err
		}
	}

	return nil
}

func (r SubscriptionRepository) update(ctx context.Context, sub subscription.Subscription) error {
	// Persist subscription main record only when it was modified.
	if sub.IsDirty() {
		if err := r.dbContext.GetQueries(ctx).UpdateSubscription(ctx, sql.UpdateSubscriptionParams{
			ID:           sub.Id(),
			FriendlyName: sub.FriendlyName(),
			FreeTrialStartDate: func() *time.Time {
				if ft := sub.FreeTrial(); ft != nil {
					start := ft.StartDate()
					return &start
				}
				return nil
			}(),
			FreeTrialEndDate: func() *time.Time {
				if ft := sub.FreeTrial(); ft != nil {
					end := ft.EndDate()
					return &end
				}
				return nil
			}(),
			ProviderID: sub.ProviderId(),
			PlanID:     sub.PlanId(),
			PriceID:    sub.PriceId(),
			CustomPriceCurrency: func() *string {
				if cp := sub.CustomPrice(); cp != nil {
					cur := cp.Currency().String()
					return &cur
				}
				return nil
			}(),
			CustomPriceAmount: func() *float64 {
				if cp := sub.CustomPrice(); cp != nil {
					amt := cp.Amount()
					return &amt
				}
				return nil
			}(),
			OwnerType: sub.Owner().Type().String(),
			OwnerFamilyID: func() *uuid.UUID {
				if sub.Owner().Type() == auth.FamilyOwnerType {
					fid := sub.Owner().FamilyId()
					return &fid
				}
				return nil
			}(),
			OwnerUserID: func() *string {
				if sub.Owner().Type() == auth.PersonalOwnerType {
					uid := sub.Owner().UserId()
					return &uid
				}
				return nil
			}(),
			PayerType: func() *string {
				if payer := sub.Payer(); payer != nil {
					pt := payer.Type().String()
					return &pt
				}
				return nil
			}(),
			FamilyID: func() *uuid.UUID {
				if payer := sub.Payer(); payer != nil {
					fid := payer.FamilyId()
					return &fid
				}
				return nil
			}(),
			PayerMemberID: func() *uuid.UUID {
				if payer := sub.Payer(); payer != nil && payer.Type() == subscription.FamilyMemberPayer {
					mid := payer.MemberId()
					return &mid
				}
				return nil
			}(),
			StartDate:        sub.StartDate(),
			EndDate:          sub.EndDate(),
			Recurrency:       sub.Recurrency().String(),
			CustomRecurrency: sub.CustomRecurrency(),
			UpdatedAt:        sub.UpdatedAt(),
			Etag:             sub.ETag(),
		}); err != nil {
			return err
		}
	}

	// Handle tracked changes for service users (additions & deletions, no updates required).
	if err := saveTrackedSlice(
		ctx,
		r.dbContext,
		sub.ServiceUsers(),
		// Create new links
		func(ctx context.Context, queries *sql.Queries, users []uuid.UUID) error {
			if len(users) == 0 {
				return nil
			}
			args := slicesx.Select(users, func(u uuid.UUID) sql.CreateSubscriptionServiceUsersParams {
				return sql.CreateSubscriptionServiceUsersParams{
					SubscriptionID: sub.Id(),
					FamilyMemberID: u,
				}
			})
			_, err := queries.CreateSubscriptionServiceUsers(ctx, args)
			return err
		},
		// Update not required – pass nil.
		nil,
		// Delete removed links
		func(ctx context.Context, queries *sql.Queries, user uuid.UUID) error {
			return queries.DeleteSubscriptionServiceUser(ctx, sql.DeleteSubscriptionServiceUserParams{
				SubscriptionID: sub.Id(),
				FamilyMemberID: user,
			})
		},
	); err != nil {
		return err
	}

	// Clear change tracking on successful persistence.
	sub.ServiceUsers().ClearChanges()
	return nil
}
