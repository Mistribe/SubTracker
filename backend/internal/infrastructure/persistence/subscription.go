package persistence

import (
	"context"
	"fmt"
	"iter"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/internal/domain/entity"
	"github.com/oleexo/subtracker/internal/domain/subscription"
	"github.com/oleexo/subtracker/internal/infrastructure/persistence/jet/app/public/model"
	"github.com/oleexo/subtracker/pkg/slicesx"

	. "github.com/go-jet/jet/v2/postgres"

	. "github.com/oleexo/subtracker/internal/infrastructure/persistence/jet/app/public/table"
)

const batchSize = 10

type SubscriptionRepository struct {
	dbContext *DatabaseContext
}

func NewSubscriptionRepository(repository *DatabaseContext) subscription.Repository {
	return &SubscriptionRepository{
		dbContext: repository,
	}
}

func (r SubscriptionRepository) GetById(ctx context.Context, id uuid.UUID) (subscription.Subscription, error) {
	stmt := SELECT(
		Subscriptions.AllColumns,
		SubscriptionServiceUsers.FamilyMemberID,
		SubscriptionServiceUsers.SubscriptionID,
	).
		FROM(
			Subscriptions.
				LEFT_JOIN(SubscriptionServiceUsers, SubscriptionServiceUsers.SubscriptionID.EQ(Subscriptions.ID)),
		).
		WHERE(Subscriptions.ID.EQ(UUID(id)))

	var rows []struct {
		model.Subscriptions
		FamilyMemberID *uuid.UUID `alias:"subscription_service_users.family_member_id"`
		SubscriptionID *uuid.UUID `alias:"subscription_service_users.subscription_id"`
	}

	if err := r.dbContext.Query(ctx, stmt, &rows); err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, nil
	}

	subscriptions := createSubscriptionFromJetRowsFlat(rows)
	if len(subscriptions) == 0 {
		return nil, nil
	}

	return subscriptions[0], nil
}

func (r SubscriptionRepository) GetByIdForUser(ctx context.Context, userId string, id uuid.UUID) (
	subscription.Subscription,
	error) {

	stmt := SELECT(
		Subscriptions.AllColumns,
		SubscriptionServiceUsers.FamilyMemberID,
		SubscriptionServiceUsers.SubscriptionID,
	).
		FROM(
			Subscriptions.
				LEFT_JOIN(SubscriptionServiceUsers, SubscriptionServiceUsers.SubscriptionID.EQ(Subscriptions.ID)).
				LEFT_JOIN(FamilyMembers, FamilyMembers.ID.EQ(SubscriptionServiceUsers.FamilyMemberID)),
		).
		WHERE(
			Subscriptions.ID.EQ(UUID(id)).
				AND(
					Subscriptions.OwnerType.EQ(String("family")).AND(FamilyMembers.UserID.EQ(String(userId))).
						OR(Subscriptions.OwnerType.EQ(String("personal")).AND(Subscriptions.OwnerUserID.EQ(String(userId)))),
				),
		)

	var rows []struct {
		model.Subscriptions
		FamilyMemberID *uuid.UUID `alias:"subscription_service_users.family_member_id"`
		SubscriptionID *uuid.UUID `alias:"subscription_service_users.subscription_id"`
	}

	if err := r.dbContext.Query(ctx, stmt, &rows); err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, nil
	}

	subscriptions := createSubscriptionFromJetRowsFlat(rows)
	if len(subscriptions) == 0 {
		return nil, nil
	}

	return subscriptions[0], nil
}

func (r SubscriptionRepository) GetAll(
	ctx context.Context,
	parameters subscription.QueryParameters) ([]subscription.Subscription, int64, error) {

	pagedSubs := SELECT(
		Subscriptions.AllColumns,
		COUNT(STAR).OVER().AS("total_count"),
	).
		FROM(Subscriptions).
		ORDER_BY(Subscriptions.ID).
		LIMIT(parameters.Limit).
		OFFSET(parameters.Offset).
		AsTable("s")

	stmt := SELECT(
		pagedSubs.AllColumns(),
		SubscriptionServiceUsers.FamilyMemberID,
		SubscriptionServiceUsers.SubscriptionID,
	).
		FROM(
			pagedSubs.
				LEFT_JOIN(SubscriptionServiceUsers,
					SubscriptionServiceUsers.SubscriptionID.EQ(Subscriptions.ID.From(pagedSubs))),
		)

	var rows []struct {
		model.Subscriptions
		FamilyMemberID *uuid.UUID `alias:"subscription_service_users.family_member_id"`
		SubscriptionID *uuid.UUID `alias:"subscription_service_users.subscription_id"`
		TotalCount     int64      `alias:"total_count"`
		IsActive       bool       `alias:"is_active"`
	}

	if err := r.dbContext.Query(ctx, stmt, &rows); err != nil {
		return nil, 0, err
	}

	if len(rows) == 0 {
		return nil, 0, nil
	}

	totalCount := rows[0].TotalCount
	subscriptions := createSubscriptionFromJetRowsFlatWithCount(rows)
	if len(subscriptions) == 0 {
		return nil, 0, nil
	}

	return subscriptions, totalCount, nil
}

func (r SubscriptionRepository) GetAllForUser(
	ctx context.Context,
	userId string,
	parameters subscription.QueryParameters) ([]subscription.Subscription, int64, error) {

	// Build accessible providers CTE
	providers := SELECT(
		Providers.ID,
		Providers.Name,
	).FROM(Providers).
		WHERE(
			Providers.OwnerType.EQ(String("system")).
				OR(Providers.OwnerType.EQ(String("personal")).AND(Providers.OwnerUserID.EQ(String(userId)))).
				OR(Providers.OwnerType.EQ(String("family")).AND(EXISTS(
					SELECT(FamilyMembers.ID).
						FROM(FamilyMembers).
						WHERE(FamilyMembers.FamilyID.EQ(Providers.OwnerFamilyID).AND(FamilyMembers.UserID.EQ(String(userId))))),
				)),
		).AsTable("providers")

	// Build access filter
	accessFilter := Subscriptions.OwnerType.EQ(String("system")).
		OR(Subscriptions.OwnerType.EQ(String("personal")).AND(Subscriptions.OwnerUserID.EQ(String(userId)))).
		OR(Subscriptions.OwnerType.EQ(String("family")).AND(EXISTS(
			SELECT(FamilyMembers.ID).
				FROM(FamilyMembers).
				WHERE(FamilyMembers.FamilyID.EQ(Subscriptions.OwnerFamilyID).AND(FamilyMembers.UserID.EQ(String(userId))))),
		))

	// Add search filter if provided
	searchFilter := accessFilter
	if parameters.SearchText != "" {
		searchTerm := fmt.Sprintf("%%%s%%", parameters.SearchText)
		searchFilter = accessFilter.AND(
			Subscriptions.FriendlyName.LIKE(String(searchTerm)).
				OR(EXISTS(
					SELECT(Providers.ID.From(providers)).
						FROM(providers).
						WHERE(Providers.Name.From(providers).LIKE(String(searchTerm)).AND(Providers.ID.From(providers).EQ(Subscriptions.ProviderID))),
				)),
		)
	}

	// Add recurrency filter if provided
	finalFilter := searchFilter
	if len(parameters.Recurrencies) > 0 {
		recurrencyVals := make([]Expression, len(parameters.Recurrencies))
		for i, rec := range parameters.Recurrencies {
			recurrencyVals[i] = String(rec.String())
		}
		finalFilter = searchFilter.AND(Subscriptions.Recurrency.IN(recurrencyVals...))
	}

	// Build matches CTE with is_active calculation
	matches := SELECT(
		Subscriptions.AllColumns,
		CAST(
			NOW().GT_EQ(TimestampzExp(Subscriptions.StartDate)).
				AND(Subscriptions.EndDate.IS_NULL().OR(NOW().LT_EQ(TimestampzExp(Subscriptions.EndDate)))),
		).AS_BOOL().AS("is_active"),
		providers.AllColumns(),
	).FROM(
		Subscriptions.LEFT_JOIN(providers, Providers.ID.From(providers).EQ(Subscriptions.ProviderID)),
	).WHERE(finalFilter).
		AsTable("matches")

	// Build sorting
	var orderBy []OrderByClause
	if parameters.SortBy != "" && parameters.SortOrder != "" {
		switch parameters.SortBy {
		case "provider":
			if parameters.SortOrder == "ASC" {
				orderBy = append(orderBy, Providers.Name.From(matches).ASC())
			} else {
				orderBy = append(orderBy, Providers.Name.From(matches).DESC())
			}
		case "name":
			if parameters.SortOrder == "ASC" {
				orderBy = append(orderBy, Subscriptions.FriendlyName.From(matches).ASC())
			} else {
				orderBy = append(orderBy, Subscriptions.FriendlyName.From(matches).DESC())
			}
		case "price":
			if parameters.SortOrder == "ASC" {
				orderBy = append(orderBy, Subscriptions.CustomPriceAmount.From(matches).ASC())
			} else {
				orderBy = append(orderBy, Subscriptions.CustomPriceAmount.From(matches).DESC())
			}
		case "recurrency":
			if parameters.SortOrder == "ASC" {
				orderBy = append(orderBy, Subscriptions.Recurrency.From(matches).ASC())
			} else {
				orderBy = append(orderBy, Subscriptions.Recurrency.From(matches).DESC())
			}
		case "dates":
			if parameters.SortOrder == "ASC" {
				orderBy = append(orderBy, Subscriptions.StartDate.From(matches).ASC())
			} else {
				orderBy = append(orderBy, Subscriptions.StartDate.From(matches).DESC())
			}
		case "status":
			if parameters.SortOrder == "ASC" {
				orderBy = append(orderBy, BoolColumn("is_active").From(matches).ASC())
			} else {
				orderBy = append(orderBy, BoolColumn("is_active").From(matches).DESC())
			}
		}
	}
	orderBy = append(orderBy, Subscriptions.ID.From(matches).ASC())

	// Build counted CTE with pagination
	counted := SELECT(
		matches.AllColumns(),
		COUNT(STAR).OVER().AS("total_count"),
	).FROM(matches).
		ORDER_BY(orderBy...).
		LIMIT(parameters.Limit).
		OFFSET(parameters.Offset).
		AsTable("paged")

	// Final query with service users
	stmt := SELECT(
		counted.AllColumns(),
		SubscriptionServiceUsers.FamilyMemberID,
		SubscriptionServiceUsers.SubscriptionID,
	).FROM(
		counted.
			LEFT_JOIN(SubscriptionServiceUsers,
				SubscriptionServiceUsers.SubscriptionID.EQ(Subscriptions.ID.From(counted))),
	)

	var rows []struct {
		model.Subscriptions
		FamilyMemberID *uuid.UUID `alias:"subscription_service_users.family_member_id"`
		SubscriptionID *uuid.UUID `alias:"subscription_service_users.subscription_id"`
		TotalCount     int64      `alias:"total_count"`
		IsActive       bool       `alias:"is_active"`
	}

	if err := r.dbContext.Query(ctx, stmt, &rows); err != nil {
		return nil, 0, err
	}

	if len(rows) == 0 {
		return nil, 0, nil
	}

	totalCount := rows[0].TotalCount
	subscriptions := createSubscriptionFromJetRowsFlatWithCount(rows)
	if len(subscriptions) == 0 {
		return nil, 0, nil
	}

	return subscriptions, totalCount, nil
}

func (r SubscriptionRepository) GetAllIt(
	ctx context.Context,
	userId, searchText string) iter.Seq[subscription.Subscription] {
	return func(yield func(subscription.Subscription) bool) {
		offset := int64(0)
		for {
			subs, _, err := r.GetAllForUser(ctx, userId, subscription.QueryParameters{
				QueryParameters: entity.QueryParameters{
					Limit:  batchSize,
					Offset: offset,
				},
				SearchText: searchText,
			})
			if err != nil || len(subs) == 0 {
				return
			}

			for _, sub := range subs {
				if !yield(sub) {
					return
				}
			}

			if len(subs) < batchSize {
				return
			}
			offset += batchSize
		}
	}
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
	stmt := Subscriptions.DELETE().
		WHERE(Subscriptions.ID.EQ(UUID(subscriptionId)))

	count, err := r.dbContext.Execute(ctx, stmt)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r SubscriptionRepository) Exists(ctx context.Context, ids ...uuid.UUID) (bool, error) {
	if len(ids) == 0 {
		return true, nil
	}

	vals := make([]Expression, len(ids))
	for i, id := range ids {
		vals[i] = UUID(id)
	}

	stmt := SELECT(COUNT(Subscriptions.ID).AS("count")).
		FROM(Subscriptions).
		WHERE(Subscriptions.ID.IN(vals...))

	var row struct {
		Count int
	}

	if err := r.dbContext.Query(ctx, stmt, &row); err != nil {
		return false, err
	}

	return row.Count == len(ids), nil
}

func (r SubscriptionRepository) create(ctx context.Context, subscriptions []subscription.Subscription) error {
	if len(subscriptions) == 0 {
		return nil
	}

	// Insert subscriptions
	stmt := Subscriptions.INSERT(
		Subscriptions.ID,
		Subscriptions.FriendlyName,
		Subscriptions.FreeTrialStartDate,
		Subscriptions.FreeTrialEndDate,
		Subscriptions.ProviderID,
		Subscriptions.PlanID,
		Subscriptions.PriceID,
		Subscriptions.CustomPriceCurrency,
		Subscriptions.CustomPriceAmount,
		Subscriptions.OwnerType,
		Subscriptions.OwnerFamilyID,
		Subscriptions.OwnerUserID,
		Subscriptions.PayerType,
		Subscriptions.FamilyID,
		Subscriptions.PayerMemberID,
		Subscriptions.StartDate,
		Subscriptions.EndDate,
		Subscriptions.Recurrency,
		Subscriptions.CustomRecurrency,
		Subscriptions.CreatedAt,
		Subscriptions.UpdatedAt,
		Subscriptions.Etag,
	)

	for _, sub := range subscriptions {
		var friendlyNameVal Expression
		if sub.FriendlyName() != nil {
			friendlyNameVal = String(*sub.FriendlyName())
		} else {
			friendlyNameVal = NULL
		}

		var freeTrialStartVal, freeTrialEndVal Expression
		if sub.FreeTrial() != nil {
			freeTrialStartVal = TimestampzT(sub.FreeTrial().StartDate())
			freeTrialEndVal = TimestampzT(sub.FreeTrial().EndDate())
		} else {
			freeTrialStartVal = NULL
			freeTrialEndVal = NULL
		}

		var planIdVal, priceIdVal Expression
		if sub.PlanId() != nil {
			planIdVal = UUID(*sub.PlanId())
		} else {
			planIdVal = NULL
		}
		if sub.PriceId() != nil {
			priceIdVal = UUID(*sub.PriceId())
		} else {
			priceIdVal = NULL
		}

		var customPriceCurrencyVal, customPriceAmountVal Expression
		if sub.CustomPrice() != nil {
			customPriceCurrencyVal = String(sub.CustomPrice().Currency().String())
			customPriceAmountVal = Float(sub.CustomPrice().Value())
		} else {
			customPriceCurrencyVal = NULL
			customPriceAmountVal = NULL
		}

		var ownerFamilyIdVal, ownerUserIdVal Expression
		switch sub.Owner().Type() {
		case auth.PersonalOwnerType:
			ownerFamilyIdVal = NULL
			ownerUserIdVal = String(sub.Owner().UserId())
		case auth.FamilyOwnerType:
			ownerFamilyIdVal = UUID(sub.Owner().FamilyId())
			ownerUserIdVal = NULL
		default:
			ownerFamilyIdVal = NULL
			ownerUserIdVal = NULL
		}

		var payerTypeVal, familyIdVal, payerMemberIdVal Expression
		if sub.Payer() != nil {
			payerTypeVal = String(sub.Payer().Type().String())
			familyIdVal = UUID(sub.Payer().FamilyId())
			if sub.Payer().Type() == subscription.FamilyMemberPayer {
				payerMemberIdVal = UUID(sub.Payer().MemberId())
			} else {
				payerMemberIdVal = NULL
			}
		} else {
			payerTypeVal = NULL
			familyIdVal = NULL
			payerMemberIdVal = NULL
		}

		var endDateVal Expression
		if sub.EndDate() != nil {
			endDateVal = TimestampzT(*sub.EndDate())
		} else {
			endDateVal = NULL
		}

		var customRecurrencyVal Expression
		if sub.CustomRecurrency() != nil {
			customRecurrencyVal = Int32(*sub.CustomRecurrency())
		} else {
			customRecurrencyVal = NULL
		}

		stmt = stmt.VALUES(
			UUID(sub.Id()),
			friendlyNameVal,
			freeTrialStartVal,
			freeTrialEndVal,
			UUID(sub.ProviderId()),
			planIdVal,
			priceIdVal,
			customPriceCurrencyVal,
			customPriceAmountVal,
			String(sub.Owner().Type().String()),
			ownerFamilyIdVal,
			ownerUserIdVal,
			payerTypeVal,
			familyIdVal,
			payerMemberIdVal,
			TimestampzT(sub.StartDate()),
			endDateVal,
			String(sub.Recurrency().String()),
			customRecurrencyVal,
			TimestampzT(sub.CreatedAt()),
			TimestampzT(sub.UpdatedAt()),
			String(sub.ETag()),
		)
	}

	count, err := r.dbContext.Execute(ctx, stmt)
	if err != nil {
		return err
	}
	if count != int64(len(subscriptions)) {
		return ErrMissMatchAffectRow
	}

	// Insert service users
	allServiceUsers := slicesx.SelectMany(subscriptions,
		func(sub subscription.Subscription) []struct {
			SubscriptionID uuid.UUID
			FamilyMemberID uuid.UUID
		} {
			return slicesx.Select(sub.ServiceUsers().Values(),
				func(u uuid.UUID) struct {
					SubscriptionID uuid.UUID
					FamilyMemberID uuid.UUID
				} {
					return struct {
						SubscriptionID uuid.UUID
						FamilyMemberID uuid.UUID
					}{SubscriptionID: sub.Id(), FamilyMemberID: u}
				})
		})

	if len(allServiceUsers) > 0 {
		serviceUserStmt := SubscriptionServiceUsers.INSERT(
			SubscriptionServiceUsers.SubscriptionID,
			SubscriptionServiceUsers.FamilyMemberID,
		)

		for _, su := range allServiceUsers {
			serviceUserStmt = serviceUserStmt.VALUES(
				UUID(su.SubscriptionID),
				UUID(su.FamilyMemberID),
			)
		}

		if _, err := r.dbContext.Execute(ctx, serviceUserStmt); err != nil {
			return err
		}
	}

	return nil
}

func (r SubscriptionRepository) update(ctx context.Context, sub subscription.Subscription) error {
	// Persist subscription main record only when it was modified.
	if sub.IsDirty() {
		var friendlyNameVal StringExpression
		if sub.FriendlyName() != nil {
			friendlyNameVal = String(*sub.FriendlyName())
		} else {
			friendlyNameVal = StringExp(NULL)
		}

		var freeTrialStartVal, freeTrialEndVal TimestampExpression
		if sub.FreeTrial() != nil {
			freeTrialStartVal = TimestampT(sub.FreeTrial().StartDate())
			freeTrialEndVal = TimestampT(sub.FreeTrial().EndDate())
		} else {
			freeTrialStartVal = TimestampExp(NULL)
			freeTrialEndVal = TimestampExp(NULL)
		}

		var planIdVal, priceIdVal StringExpression
		if sub.PlanId() != nil {
			planIdVal = StringExp(UUID(*sub.PlanId()))
		} else {
			planIdVal = StringExp(NULL)
		}
		if sub.PriceId() != nil {
			priceIdVal = StringExp(UUID(*sub.PriceId()))
		} else {
			priceIdVal = StringExp(NULL)
		}

		var customPriceCurrencyVal StringExpression
		var customPriceAmountVal FloatExpression
		if sub.CustomPrice() != nil {
			customPriceCurrencyVal = String(sub.CustomPrice().Currency().String())
			customPriceAmountVal = Float(sub.CustomPrice().Value())
		} else {
			customPriceCurrencyVal = StringExp(NULL)
			customPriceAmountVal = FloatExp(NULL)
		}

		var ownerFamilyIdVal, ownerUserIdVal StringExpression
		switch sub.Owner().Type() {
		case auth.PersonalOwnerType:
			ownerFamilyIdVal = StringExp(NULL)
			ownerUserIdVal = String(sub.Owner().UserId())
		case auth.FamilyOwnerType:
			ownerFamilyIdVal = StringExp(UUID(sub.Owner().FamilyId()))
			ownerUserIdVal = StringExp(NULL)
		default:
			ownerFamilyIdVal = StringExp(NULL)
			ownerUserIdVal = StringExp(NULL)
		}

		var payerTypeVal, familyIdVal, payerMemberIdVal StringExpression
		if sub.Payer() != nil {
			payerTypeVal = String(sub.Payer().Type().String())
			familyIdVal = StringExp(UUID(sub.Payer().FamilyId()))
			if sub.Payer().Type() == subscription.FamilyMemberPayer {
				payerMemberIdVal = StringExp(UUID(sub.Payer().MemberId()))
			} else {
				payerMemberIdVal = StringExp(NULL)
			}
		} else {
			payerTypeVal = StringExp(NULL)
			familyIdVal = StringExp(NULL)
			payerMemberIdVal = StringExp(NULL)
		}

		var endDateVal TimestampExpression
		if sub.EndDate() != nil {
			endDateVal = TimestampT(*sub.EndDate())
		} else {
			endDateVal = TimestampExp(NULL)
		}

		var customRecurrencyVal IntegerExpression
		if sub.CustomRecurrency() != nil {
			customRecurrencyVal = Int32(*sub.CustomRecurrency())
		} else {
			customRecurrencyVal = IntExp(NULL)
		}

		stmt := Subscriptions.UPDATE().
			SET(
				Subscriptions.FriendlyName.SET(friendlyNameVal),
				Subscriptions.FreeTrialStartDate.SET(freeTrialStartVal),
				Subscriptions.FreeTrialEndDate.SET(freeTrialEndVal),
				Subscriptions.ProviderID.SET(UUID(sub.ProviderId())),
				Subscriptions.PlanID.SET(planIdVal),
				Subscriptions.PriceID.SET(priceIdVal),
				Subscriptions.CustomPriceCurrency.SET(customPriceCurrencyVal),
				Subscriptions.CustomPriceAmount.SET(customPriceAmountVal),
				Subscriptions.OwnerType.SET(String(sub.Owner().Type().String())),
				Subscriptions.OwnerFamilyID.SET(ownerFamilyIdVal),
				Subscriptions.OwnerUserID.SET(ownerUserIdVal),
				Subscriptions.PayerType.SET(payerTypeVal),
				Subscriptions.FamilyID.SET(familyIdVal),
				Subscriptions.PayerMemberID.SET(payerMemberIdVal),
				Subscriptions.StartDate.SET(TimestampT(sub.StartDate())),
				Subscriptions.EndDate.SET(endDateVal),
				Subscriptions.Recurrency.SET(String(sub.Recurrency().String())),
				Subscriptions.CustomRecurrency.SET(customRecurrencyVal),
				Subscriptions.UpdatedAt.SET(TimestampT(sub.UpdatedAt())),
				Subscriptions.Etag.SET(String(sub.ETag())),
			).
			WHERE(Subscriptions.ID.EQ(UUID(sub.Id())))

		count, err := r.dbContext.Execute(ctx, stmt)
		if err != nil {
			return err
		}
		if count == 0 {
			return ErrMissMatchAffectRow
		}
	}

	// Handle tracked changes for service users (additions & deletions, no updates required).
	if err := r.saveTrackedServiceUsersWithJet(ctx, sub.Id(), sub.ServiceUsers()); err != nil {
		return err
	}

	// Clear change tracking on successful persistence.
	sub.ServiceUsers().ClearChanges()
	return nil
}

func (r SubscriptionRepository) saveTrackedServiceUsersWithJet(ctx context.Context, subscriptionId uuid.UUID,
	serviceUsers *slicesx.Tracked[uuid.UUID]) error {
	// Handle new service users
	newUsers := serviceUsers.Added()
	if len(newUsers) > 0 {
		stmt := SubscriptionServiceUsers.INSERT(
			SubscriptionServiceUsers.SubscriptionID,
			SubscriptionServiceUsers.FamilyMemberID,
		)

		for _, userId := range newUsers {
			stmt = stmt.VALUES(
				UUID(subscriptionId),
				UUID(userId),
			)
		}

		_, err := r.dbContext.Execute(ctx, stmt)
		if err != nil {
			return err
		}
	}

	// Handle deleted service users
	deletedUsers := serviceUsers.Removed()
	for _, userId := range deletedUsers {
		stmt := SubscriptionServiceUsers.DELETE().
			WHERE(
				SubscriptionServiceUsers.SubscriptionID.EQ(UUID(subscriptionId)).
					AND(SubscriptionServiceUsers.FamilyMemberID.EQ(UUID(userId))),
			)

		_, err := r.dbContext.Execute(ctx, stmt)
		if err != nil {
			return err
		}
	}

	return nil
}
