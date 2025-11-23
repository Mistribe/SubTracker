package repositories

import (
	"context"
	"fmt"
	"iter"

	"github.com/mistribe/subtracker/internal/adapters/persistence/db"
	. "github.com/mistribe/subtracker/internal/adapters/persistence/db/jet/app/public/table"
	"github.com/mistribe/subtracker/internal/adapters/persistence/db/models"
	"github.com/mistribe/subtracker/internal/domain/subscription"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/slicesx"
	"github.com/mistribe/subtracker/pkg/x/herd"

	. "github.com/go-jet/jet/v2/postgres"
)

const batchSize = 10

type SubscriptionRepository struct {
	dbContext *db.Context
}

func NewSubscriptionRepository(repository *db.Context) ports.SubscriptionRepository {
	return &SubscriptionRepository{
		dbContext: repository,
	}
}

func (r SubscriptionRepository) GetById(ctx context.Context, id types.SubscriptionID) (
	subscription.Subscription,
	error) {
	stmt := SELECT(
		Subscriptions.AllColumns,
		SubscriptionFamilyUsers.FamilyMemberID,
		SubscriptionFamilyUsers.SubscriptionID,
	).
		FROM(
			Subscriptions.
				LEFT_JOIN(SubscriptionFamilyUsers, SubscriptionFamilyUsers.SubscriptionID.EQ(Subscriptions.ID)),
		).
		WHERE(Subscriptions.ID.EQ(UUID(id)))

	var rows []models.SubscriptionRow

	if err := r.dbContext.Query(ctx, stmt, &rows); err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, nil
	}

	subscriptions := models.CreateSubscriptionFromJetRows(rows)
	if len(subscriptions) == 0 {
		return nil, nil
	}

	return subscriptions[0], nil
}

func (r SubscriptionRepository) GetByIdForUser(ctx context.Context, userId types.UserID, id types.SubscriptionID) (
	subscription.Subscription,
	error) {

	stmt := SELECT(
		Subscriptions.AllColumns,
		SubscriptionFamilyUsers.FamilyMemberID,
		SubscriptionFamilyUsers.SubscriptionID,
	).
		FROM(
			Subscriptions.
				LEFT_JOIN(SubscriptionFamilyUsers, SubscriptionFamilyUsers.SubscriptionID.EQ(Subscriptions.ID)),
		).
		WHERE(
			Subscriptions.ID.EQ(UUID(id)).
				AND(
					Subscriptions.OwnerType.EQ(String("family")).AND(EXISTS(
						SELECT(Families.ID).
							FROM(
								Families.
									INNER_JOIN(FamilyMembers, FamilyMembers.FamilyID.EQ(Families.ID)),
							).
							WHERE(
								Families.ID.EQ(Subscriptions.OwnerFamilyID).
									AND(FamilyMembers.UserID.EQ(String(userId.String()))),
							),
					)).
						OR(Subscriptions.OwnerType.EQ(String("personal")).AND(Subscriptions.OwnerUserID.EQ(String(userId.String())))),
				),
		)

	var rows []models.SubscriptionRow

	if err := r.dbContext.Query(ctx, stmt, &rows); err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, nil
	}

	subscriptions := models.CreateSubscriptionFromJetRows(rows)
	if len(subscriptions) == 0 {
		return nil, nil
	}

	return subscriptions[0], nil
}

func (r SubscriptionRepository) GetAll(
	ctx context.Context,
	parameters ports.SubscriptionQueryParameters) ([]subscription.Subscription, int64, error) {

	pagedSubs := SELECT(
		Subscriptions.AllColumns,
		COUNT(STAR).OVER().AS("subscription_row_with_count.total_count"),
	).
		FROM(Subscriptions).
		ORDER_BY(Subscriptions.ID).
		LIMIT(parameters.Limit).
		OFFSET(parameters.Offset).
		AsTable("s")

	stmt := SELECT(
		pagedSubs.AllColumns(),
		SubscriptionFamilyUsers.AllColumns,
		SubscriptionLabels.AllColumns,
		ProviderLabels.AllColumns,
	).
		FROM(
			pagedSubs.
				LEFT_JOIN(SubscriptionFamilyUsers,
					SubscriptionFamilyUsers.SubscriptionID.EQ(Subscriptions.ID.From(pagedSubs))).
				LEFT_JOIN(SubscriptionLabels, SubscriptionLabels.SubscriptionID.EQ(Subscriptions.ID.From(pagedSubs))).
				LEFT_JOIN(Providers, Providers.ID.EQ(Subscriptions.ProviderID.From(pagedSubs))).
				LEFT_JOIN(ProviderLabels, ProviderLabels.ProviderID.EQ(Providers.ID)),
		)

	var rows []models.SubscriptionRowWithCount

	if err := r.dbContext.Query(ctx, stmt, &rows); err != nil {
		return nil, 0, err
	}

	if len(rows) == 0 {
		return nil, 0, nil
	}

	totalCount := rows[0].TotalCount
	subscriptions := models.CreateSubscriptionFromJetRowsWithCount(rows)
	if len(subscriptions) == 0 {
		return nil, 0, nil
	}

	return subscriptions, totalCount, nil
}

func (r SubscriptionRepository) GetAllForUser(
	ctx context.Context,
	userId types.UserID,
	parameters ports.SubscriptionQueryParameters) ([]subscription.Subscription, int64, error) {

	// Build accessible providers CTE
	providers := SELECT(
		Providers.ID,
		Providers.Name,
	).FROM(Providers).
		WHERE(
			Providers.OwnerType.EQ(String("system")).
				OR(Providers.OwnerType.EQ(String("personal")).AND(Providers.OwnerUserID.EQ(String(userId.String())))).
				OR(Providers.OwnerType.EQ(String("family")).AND(EXISTS(
					SELECT(FamilyMembers.ID).
						FROM(FamilyMembers).
						WHERE(FamilyMembers.FamilyID.EQ(Providers.OwnerFamilyID).AND(FamilyMembers.UserID.EQ(String(userId.String()))))),
				)),
		).AsTable("providers")

	// Build access filter
	accessFilter := Subscriptions.OwnerType.EQ(String("system")).
		OR(Subscriptions.OwnerType.EQ(String("personal")).AND(Subscriptions.OwnerUserID.EQ(String(userId.String())))).
		OR(Subscriptions.OwnerType.EQ(String("family")).AND(EXISTS(
			SELECT(FamilyMembers.ID).
				FROM(FamilyMembers).
				WHERE(FamilyMembers.FamilyID.EQ(Subscriptions.OwnerFamilyID).AND(FamilyMembers.UserID.EQ(String(userId.String()))))),
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

	// Add from date filter if provided
	if parameters.FromDate != nil {
		searchFilter = searchFilter.AND(
			Subscriptions.StartDate.GT_EQ(TimestampzT(*parameters.FromDate)),
		)
	}

	// Add to date filter if provided
	if parameters.ToDate != nil {
		searchFilter = searchFilter.AND(
			Subscriptions.StartDate.LT_EQ(TimestampzT(*parameters.ToDate)),
		)
	}

	if len(parameters.Providers) > 0 {
		providerVals := make([]Expression, len(parameters.Providers))
		for i, provider := range parameters.Providers {
			providerVals[i] = UUID(provider)
		}
		searchFilter = searchFilter.AND(Subscriptions.ProviderID.IN(providerVals...))
	}

	// Add recurrency filter if provided
	if len(parameters.Recurrencies) > 0 {
		recurrencyVals := make([]Expression, len(parameters.Recurrencies))
		for i, rec := range parameters.Recurrencies {
			recurrencyVals[i] = String(rec.String())
		}
		searchFilter = searchFilter.AND(Subscriptions.Recurrency.IN(recurrencyVals...))
	}

	// Add filter on active only
	if !parameters.WithInactive {
		searchFilter = searchFilter.AND(
			NOW().GT_EQ(Subscriptions.StartDate).
				AND(Subscriptions.EndDate.IS_NULL().
					OR(NOW().LT_EQ(Subscriptions.EndDate))),
		)
	}

	// Build matches CTE with is_active calculation
	matches := SELECT(
		Subscriptions.AllColumns,
		providers.AllColumns(),
		CAST(
			NOW().GT_EQ(Subscriptions.StartDate).
				AND(Subscriptions.EndDate.IS_NULL().OR(NOW().LT_EQ(Subscriptions.EndDate))),
		).AS_BOOL().AS("is_active"),
	).FROM(
		Subscriptions.LEFT_JOIN(providers, Providers.ID.From(providers).EQ(Subscriptions.ProviderID)),
	).WHERE(searchFilter).AsTable("matches")

	// Build counted CTE with pagination
	counted := SELECT(
		matches.AllColumns(),
		COUNT(STAR).OVER().AS("subscription_row_with_count.total_count"),
	).FROM(matches).
		ORDER_BY(LOWER(Subscriptions.FriendlyName.From(matches)).ASC(),
			LOWER(Providers.Name.From(matches)).ASC(),
			Subscriptions.ID.From(matches).ASC()).
		LIMIT(parameters.Limit).
		OFFSET(parameters.Offset).
		AsTable("paged")

	// Final query with service users
	stmt := SELECT(
		counted.AllColumns(),
		SubscriptionFamilyUsers.AllColumns,
		SubscriptionLabels.AllColumns,
		ProviderLabels.AllColumns,
	).FROM(
		counted.
			LEFT_JOIN(SubscriptionFamilyUsers,
				SubscriptionFamilyUsers.SubscriptionID.EQ(Subscriptions.ID.From(counted))).
			LEFT_JOIN(SubscriptionLabels, SubscriptionLabels.SubscriptionID.EQ(Subscriptions.ID.From(counted))).
			LEFT_JOIN(Providers, Providers.ID.EQ(Subscriptions.ProviderID.From(counted))).
			LEFT_JOIN(ProviderLabels, ProviderLabels.ProviderID.EQ(Providers.ID)),
	)

	var rows []models.SubscriptionRowWithCount

	if err := r.dbContext.Query(ctx, stmt, &rows); err != nil {
		return nil, 0, err
	}

	if len(rows) == 0 {
		return nil, 0, nil
	}

	totalCount := rows[0].TotalCount
	subscriptions := models.CreateSubscriptionFromJetRowsWithCount(rows)
	if len(subscriptions) == 0 {
		return nil, 0, nil
	}

	return subscriptions, totalCount, nil

}

func (r SubscriptionRepository) GetAllIt(
	ctx context.Context,
	userId types.UserID, searchText string) iter.Seq[subscription.Subscription] {
	return func(yield func(subscription.Subscription) bool) {
		offset := int64(0)
		for {
			subs, _, err := r.GetAllForUser(ctx, userId, ports.SubscriptionQueryParameters{
				QueryParameters: ports.QueryParameters{
					Limit:  batchSize,
					Offset: offset,
				},
				SearchText:   searchText,
				WithInactive: true,
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

func (r SubscriptionRepository) Delete(ctx context.Context, subscriptionId types.SubscriptionID) (bool, error) {
	stmt := Subscriptions.DELETE().
		WHERE(Subscriptions.ID.EQ(UUID(subscriptionId)))

	count, err := r.dbContext.Execute(ctx, stmt)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r SubscriptionRepository) Exists(ctx context.Context, ids ...types.SubscriptionID) (bool, error) {
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
		Subscriptions.CustomPriceCurrency,
		Subscriptions.CustomPriceAmount,
		Subscriptions.OwnerType,
		Subscriptions.OwnerFamilyID,
		Subscriptions.OwnerUserID,
		Subscriptions.PayerType,
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

		var customPriceCurrencyVal, customPriceAmountVal Expression
		if sub.Price() != nil {
			customPriceCurrencyVal = String(sub.Price().Amount().Currency().String())
			customPriceAmountVal = Float(sub.Price().Amount().Value())
		} else {
			customPriceCurrencyVal = NULL
			customPriceAmountVal = NULL
		}

		var ownerFamilyIdVal, ownerUserIdVal Expression
		switch sub.Owner().Type() {
		case types.PersonalOwnerType:
			ownerFamilyIdVal = NULL
			ownerUserIdVal = String(sub.Owner().UserId().String())
		case types.FamilyOwnerType:
			ownerFamilyIdVal = UUID(sub.Owner().FamilyId())
			ownerUserIdVal = NULL
		default:
			ownerFamilyIdVal = NULL
			ownerUserIdVal = NULL
		}

		var payerTypeVal, payerMemberIdVal Expression
		if sub.Payer() != nil {
			payerTypeVal = String(sub.Payer().Type().String())
			if sub.Payer().Type() == subscription.FamilyMemberPayer {
				payerMemberIdVal = UUID(sub.Payer().MemberId())
			} else {
				payerMemberIdVal = NULL
			}
		} else {
			payerTypeVal = NULL
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
			customPriceCurrencyVal,
			customPriceAmountVal,
			String(sub.Owner().Type().String()),
			ownerFamilyIdVal,
			ownerUserIdVal,
			payerTypeVal,
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
		return db.ErrMissMatchAffectRow
	}

	// Insert service users
	allServiceUsers := herd.SelectMany(subscriptions,
		func(sub subscription.Subscription) []struct {
			SubscriptionID types.SubscriptionID
			FamilyMemberID types.FamilyMemberID
		} {
			return herd.Select(sub.FamilyUsers().Values(),
				func(u types.FamilyMemberID) struct {
					SubscriptionID types.SubscriptionID
					FamilyMemberID types.FamilyMemberID
				} {
					return struct {
						SubscriptionID types.SubscriptionID
						FamilyMemberID types.FamilyMemberID
					}{SubscriptionID: sub.Id(), FamilyMemberID: u}
				})
		})

	if len(allServiceUsers) > 0 {
		serviceUserStmt := SubscriptionFamilyUsers.INSERT(
			SubscriptionFamilyUsers.SubscriptionID,
			SubscriptionFamilyUsers.FamilyMemberID,
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

		var freeTrialStartVal, freeTrialEndVal TimestampzExpression
		if sub.FreeTrial() != nil {
			freeTrialStartVal = TimestampzT(sub.FreeTrial().StartDate())
			freeTrialEndVal = TimestampzT(sub.FreeTrial().EndDate())
		} else {
			freeTrialStartVal = TimestampzExp(NULL)
			freeTrialEndVal = TimestampzExp(NULL)
		}

		var customPriceCurrencyVal StringExpression
		var customPriceAmountVal FloatExpression
		if sub.Price() != nil {
			customPriceCurrencyVal = String(sub.Price().Amount().Currency().String())
			customPriceAmountVal = Float(sub.Price().Amount().Value())
		} else {
			customPriceCurrencyVal = StringExp(NULL)
			customPriceAmountVal = FloatExp(NULL)
		}

		var ownerFamilyIdVal, ownerUserIdVal StringExpression
		switch sub.Owner().Type() {
		case types.PersonalOwnerType:
			ownerFamilyIdVal = StringExp(NULL)
			ownerUserIdVal = String(sub.Owner().UserId().String())
		case types.FamilyOwnerType:
			ownerFamilyIdVal = StringExp(UUID(sub.Owner().FamilyId()))
			ownerUserIdVal = StringExp(NULL)
		default:
			ownerFamilyIdVal = StringExp(NULL)
			ownerUserIdVal = StringExp(NULL)
		}

		var payerTypeVal, payerMemberIdVal StringExpression
		if sub.Payer() != nil {
			payerTypeVal = String(sub.Payer().Type().String())
			if sub.Payer().Type() == subscription.FamilyMemberPayer {
				payerMemberIdVal = StringExp(UUID(sub.Payer().MemberId()))
			} else {
				payerMemberIdVal = StringExp(NULL)
			}
		} else {
			payerTypeVal = StringExp(NULL)
			payerMemberIdVal = StringExp(NULL)
		}

		var endDateVal TimestampzExpression
		if sub.EndDate() != nil {
			endDateVal = TimestampzT(*sub.EndDate())
		} else {
			endDateVal = TimestampzExp(NULL)
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
				Subscriptions.CustomPriceCurrency.SET(customPriceCurrencyVal),
				Subscriptions.CustomPriceAmount.SET(customPriceAmountVal),
				Subscriptions.OwnerType.SET(String(sub.Owner().Type().String())),
				Subscriptions.OwnerFamilyID.SET(ownerFamilyIdVal),
				Subscriptions.OwnerUserID.SET(ownerUserIdVal),
				Subscriptions.PayerType.SET(payerTypeVal),
				Subscriptions.PayerMemberID.SET(payerMemberIdVal),
				Subscriptions.StartDate.SET(TimestampzT(sub.StartDate())),
				Subscriptions.EndDate.SET(endDateVal),
				Subscriptions.Recurrency.SET(String(sub.Recurrency().String())),
				Subscriptions.CustomRecurrency.SET(customRecurrencyVal),
				Subscriptions.UpdatedAt.SET(TimestampzT(sub.UpdatedAt())),
				Subscriptions.Etag.SET(String(sub.ETag())),
			).
			WHERE(Subscriptions.ID.EQ(UUID(sub.Id())))

		count, err := r.dbContext.Execute(ctx, stmt)
		if err != nil {
			return err
		}
		if count == 0 {
			return db.ErrMissMatchAffectRow
		}
	}

	// Handle tracked changes for service users (additions & deletions, no updates required).
	if err := r.saveTrackedServiceUsersWithJet(ctx, sub.Id(), sub.FamilyUsers()); err != nil {
		return err
	}

	// Clear change tracking on successful persistence.
	sub.FamilyUsers().ClearChanges()
	return nil
}

func (r SubscriptionRepository) saveTrackedServiceUsersWithJet(
	ctx context.Context, subscriptionId types.SubscriptionID,
	serviceUsers *slicesx.Tracked[types.FamilyMemberID]) error {
	// Handle new service users
	newUsers := serviceUsers.Added()
	if len(newUsers) > 0 {
		stmt := SubscriptionFamilyUsers.INSERT(
			SubscriptionFamilyUsers.SubscriptionID,
			SubscriptionFamilyUsers.FamilyMemberID,
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
		stmt := SubscriptionFamilyUsers.DELETE().
			WHERE(
				SubscriptionFamilyUsers.SubscriptionID.EQ(UUID(subscriptionId)).
					AND(SubscriptionFamilyUsers.FamilyMemberID.EQ(UUID(userId))),
			)

		_, err := r.dbContext.Execute(ctx, stmt)
		if err != nil {
			return err
		}
	}

	return nil
}
