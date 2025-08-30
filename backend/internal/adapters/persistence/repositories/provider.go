package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/adapters/persistence/db"
	. "github.com/mistribe/subtracker/internal/adapters/persistence/db/jet/app/public/table"
	"github.com/mistribe/subtracker/internal/adapters/persistence/db/models"
	"github.com/mistribe/subtracker/internal/domain/auth"
	"github.com/mistribe/subtracker/internal/domain/provider"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/slicesx"
	"github.com/mistribe/subtracker/pkg/x/collection"

	. "github.com/go-jet/jet/v2/postgres"
)

type ProviderRepository struct {
	dbContext *db.Context
}

func NewProviderRepository(dbContext *db.Context) ports.ProviderRepository {
	return &ProviderRepository{
		dbContext: dbContext,
	}
}

func (r ProviderRepository) GetById(ctx context.Context, providerId uuid.UUID) (provider.Provider, error) {
	stmt := SELECT(
		Providers.AllColumns,
		ProviderPlans.AllColumns,
		ProviderPrices.AllColumns,
		ProviderLabels.LabelID,
		ProviderLabels.ProviderID,
	).
		FROM(
			Providers.
				LEFT_JOIN(ProviderPlans, ProviderPlans.ProviderID.EQ(Providers.ID)).
				LEFT_JOIN(ProviderPrices, ProviderPrices.PlanID.EQ(ProviderPlans.ID)).
				LEFT_JOIN(ProviderLabels, ProviderLabels.ProviderID.EQ(Providers.ID)),
		).
		WHERE(Providers.ID.EQ(UUID(providerId)))

	var rows []models.ProviderRow

	if err := r.dbContext.Query(ctx, stmt, &rows); err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, nil
	}

	providers := models.CreateProviderFromJetRows(rows)
	if len(providers) == 0 {
		return nil, nil
	}

	return providers[0], nil
}

func (r ProviderRepository) GetByIdForUser(ctx context.Context, userId string, providerId uuid.UUID) (
	provider.Provider,
	error) {
	stmt := SELECT(
		Providers.AllColumns,
		ProviderPlans.AllColumns,
		ProviderPrices.AllColumns,
		ProviderLabels.LabelID,
		ProviderLabels.ProviderID,
	).
		FROM(
			Providers.
				LEFT_JOIN(ProviderPlans, ProviderPlans.ProviderID.EQ(Providers.ID)).
				LEFT_JOIN(ProviderPrices, ProviderPrices.PlanID.EQ(ProviderPlans.ID)).
				LEFT_JOIN(ProviderLabels, ProviderLabels.ProviderID.EQ(Providers.ID)).
				LEFT_JOIN(Families, Families.ID.EQ(Providers.OwnerFamilyID)).
				LEFT_JOIN(FamilyMembers, FamilyMembers.FamilyID.EQ(Families.ID)),
		).
		WHERE(
			Providers.ID.EQ(UUID(providerId)).
				AND(
					Providers.OwnerType.EQ(String("system")).
						OR(Providers.OwnerType.EQ(String("personal")).AND(Providers.OwnerUserID.EQ(String(userId)))).
						OR(Providers.OwnerType.EQ(String("family")).AND(FamilyMembers.UserID.EQ(String(userId)))),
				),
		)

	var rows []models.ProviderRow

	if err := r.dbContext.Query(ctx, stmt, &rows); err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, nil
	}

	providers := models.CreateProviderFromJetRows(rows)
	if len(providers) == 0 {
		return nil, nil
	}

	return providers[0], nil
}

func (r ProviderRepository) GetSystemProviders(ctx context.Context) ([]provider.Provider, int64, error) {
	stmt := SELECT(
		Providers.AllColumns,
		ProviderPlans.AllColumns,
		ProviderPrices.AllColumns,
		ProviderLabels.LabelID,
		ProviderLabels.ProviderID,
		COUNT(STAR).OVER().AS("total_count"),
	).
		FROM(
			Providers.
				LEFT_JOIN(ProviderPlans, ProviderPlans.ProviderID.EQ(Providers.ID)).
				LEFT_JOIN(ProviderPrices, ProviderPrices.PlanID.EQ(ProviderPlans.ID)).
				LEFT_JOIN(ProviderLabels, ProviderLabels.ProviderID.EQ(Providers.ID)),
		).
		WHERE(
			Providers.OwnerType.EQ(String("system")).
				AND(Providers.OwnerUserID.IS_NULL()).
				AND(Providers.OwnerFamilyID.IS_NULL()),
		)

	var rows []models.ProviderRowWithCount

	if err := r.dbContext.Query(ctx, stmt, &rows); err != nil {
		return nil, 0, err
	}
	if len(rows) == 0 {
		return nil, 0, nil
	}

	totalCount := rows[0].TotalCount
	providers := models.CreateProviderFromJetRowsWithCount(rows)

	if len(providers) == 0 {
		return nil, 0, nil
	}

	return providers, totalCount, nil
}

func (r ProviderRepository) GetAll(ctx context.Context, parameters ports.ProviderQueryParameters) (
	[]provider.Provider,
	int64,
	error) {

	pagedProviders := SELECT(
		Providers.AllColumns,
		COUNT(STAR).OVER().AS("provider_row_with_count.total_count"),
	).
		FROM(Providers).
		ORDER_BY(Providers.ID).
		LIMIT(parameters.Limit).
		OFFSET(parameters.Offset).
		AsTable("p")

	stmt := SELECT(
		pagedProviders.AllColumns(),
		ProviderPlans.AllColumns,
		ProviderPrices.AllColumns,
		ProviderLabels.LabelID,
		ProviderLabels.ProviderID,
	).
		FROM(
			pagedProviders.
				LEFT_JOIN(ProviderPlans, ProviderPlans.ProviderID.EQ(Providers.ID.From(pagedProviders))).
				LEFT_JOIN(ProviderPrices, ProviderPrices.PlanID.EQ(ProviderPlans.ID)).
				LEFT_JOIN(ProviderLabels, ProviderLabels.ProviderID.EQ(Providers.ID.From(pagedProviders))),
		)

	var rows []models.ProviderRowWithCount

	if err := r.dbContext.Query(ctx, stmt, &rows); err != nil {
		return nil, 0, err
	}
	if len(rows) == 0 {
		return nil, 0, nil
	}

	totalCount := rows[0].TotalCount
	providers := models.CreateProviderFromJetRowsWithCount(rows)

	if len(providers) == 0 {
		return nil, 0, nil
	}

	return providers, totalCount, nil
}

func (r ProviderRepository) GetAllForUser(
	ctx context.Context,
	userId string,
	parameters ports.ProviderQueryParameters) ([]provider.Provider, int64, error) {

	// Build base access filter
	accessFilter := Providers.OwnerType.EQ(String("system")).
		OR(Providers.OwnerType.EQ(String("personal")).AND(Providers.OwnerUserID.EQ(String(userId)))).
		OR(Providers.OwnerType.EQ(String("family")).AND(EXISTS(
			SELECT(FamilyMembers.ID).
				FROM(FamilyMembers).
				WHERE(FamilyMembers.FamilyID.EQ(Providers.OwnerFamilyID).AND(FamilyMembers.UserID.EQ(String(userId))))),
		))

	// Add search filter if provided
	searchFilter := accessFilter
	if parameters.SearchText != "" {
		searchTerm := fmt.Sprintf("%%%s%%", parameters.SearchText)
		searchFilter = accessFilter.AND(
			Providers.Name.LIKE(String(searchTerm)).
				OR(EXISTS(
					SELECT(ProviderLabels.ProviderID).
						FROM(ProviderLabels.INNER_JOIN(Labels, Labels.ID.EQ(ProviderLabels.LabelID))).
						WHERE(ProviderLabels.ProviderID.EQ(Providers.ID).AND(Labels.Name.LIKE(String(searchTerm)))),
				)),
		)
	}

	// Get matching provider IDs with pagination
	pagedProviders := SELECT(
		Providers.ID,
		COUNT(STAR).OVER().AS("provider_row_with_count.total_count"),
	).
		FROM(Providers).
		WHERE(searchFilter).
		ORDER_BY(Providers.ID).
		LIMIT(parameters.Limit).
		OFFSET(parameters.Offset).
		AsTable("paged")

	// Get full provider data with joins
	stmt := SELECT(
		Providers.AllColumns,
		ProviderPlans.AllColumns,
		ProviderPrices.AllColumns,
		ProviderLabels.LabelID,
		ProviderLabels.ProviderID,
		pagedProviders.AllColumns().Except(Providers.ID),
	).
		FROM(
			pagedProviders.
				INNER_JOIN(Providers, Providers.ID.EQ(Providers.ID.From(pagedProviders))).
				LEFT_JOIN(ProviderPlans, ProviderPlans.ProviderID.EQ(Providers.ID)).
				LEFT_JOIN(ProviderPrices, ProviderPrices.PlanID.EQ(ProviderPlans.ID)).
				LEFT_JOIN(ProviderLabels, ProviderLabels.ProviderID.EQ(Providers.ID)),
		)

	var rows []models.ProviderRowWithCount

	if err := r.dbContext.Query(ctx, stmt, &rows); err != nil {
		return nil, 0, err
	}
	if len(rows) == 0 {
		return nil, 0, nil
	}

	totalCount := rows[0].TotalCount
	providers := models.CreateProviderFromJetRowsWithCount(rows)

	if len(providers) == 0 {
		return nil, 0, nil
	}

	return providers, totalCount, nil
}

func (r ProviderRepository) Save(ctx context.Context, providers ...provider.Provider) error {
	var newProviders []provider.Provider
	for _, prov := range providers {
		if !prov.IsExists() {
			newProviders = append(newProviders, prov)
		} else {
			if err := r.update(ctx, prov); err != nil {
				return err
			}
		}
	}

	if len(newProviders) > 0 {
		if err := r.create(ctx, newProviders); err != nil {
			return err
		}
	}

	for _, prov := range providers {
		prov.Clean()
	}

	return nil
}

func (r ProviderRepository) Delete(ctx context.Context, providerId uuid.UUID) (bool, error) {
	stmt := Providers.DELETE().
		WHERE(Providers.ID.EQ(UUID(providerId)))

	count, err := r.dbContext.Execute(ctx, stmt)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r ProviderRepository) DeletePlan(ctx context.Context, planId uuid.UUID) (bool, error) {
	stmt := ProviderPlans.DELETE().
		WHERE(ProviderPlans.ID.EQ(UUID(planId)))

	count, err := r.dbContext.Execute(ctx, stmt)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r ProviderRepository) DeletePrice(ctx context.Context, priceId uuid.UUID) (bool, error) {
	stmt := ProviderPrices.DELETE().
		WHERE(ProviderPrices.ID.EQ(UUID(priceId)))

	count, err := r.dbContext.Execute(ctx, stmt)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r ProviderRepository) Exists(ctx context.Context, ids ...uuid.UUID) (bool, error) {
	if len(ids) == 0 {
		return true, nil
	}

	vals := make([]Expression, len(ids))
	for i, id := range ids {
		vals[i] = UUID(id)
	}

	stmt := SELECT(COUNT(Providers.ID).AS("count")).
		FROM(Providers).
		WHERE(Providers.ID.IN(vals...))

	var row struct {
		Count int
	}

	if err := r.dbContext.Query(ctx, stmt, &row); err != nil {
		return false, err
	}

	return row.Count == len(ids), nil
}

func (r ProviderRepository) create(ctx context.Context, providers []provider.Provider) error {
	if len(providers) == 0 {
		return nil
	}

	// Insert providers
	stmt := Providers.INSERT(
		Providers.ID,
		Providers.OwnerType,
		Providers.OwnerFamilyID,
		Providers.OwnerUserID,
		Providers.Name,
		Providers.Key,
		Providers.Description,
		Providers.IconURL,
		Providers.URL,
		Providers.PricingPageURL,
		Providers.CreatedAt,
		Providers.UpdatedAt,
		Providers.Etag,
	)

	for _, prov := range providers {
		var ownerFamilyID Expression
		var ownerUserID Expression

		switch prov.Owner().Type() {
		case auth.PersonalOwnerType:
			ownerFamilyID = NULL
			ownerUserID = String(prov.Owner().UserId())
		case auth.FamilyOwnerType:
			ownerFamilyID = UUID(prov.Owner().FamilyId())
			ownerUserID = NULL
		default:
			ownerFamilyID = NULL
			ownerUserID = NULL
		}

		var keyVal Expression
		if prov.Key() != nil {
			keyVal = String(*prov.Key())
		} else {
			keyVal = NULL
		}

		var descVal Expression
		if prov.Description() != nil {
			descVal = String(*prov.Description())
		} else {
			descVal = NULL
		}

		var iconVal Expression
		if prov.IconUrl() != nil {
			iconVal = String(*prov.IconUrl())
		} else {
			iconVal = NULL
		}

		var urlVal Expression
		if prov.Url() != nil {
			urlVal = String(*prov.Url())
		} else {
			urlVal = NULL
		}

		var pricingVal Expression
		if prov.PricingPageUrl() != nil {
			pricingVal = String(*prov.PricingPageUrl())
		} else {
			pricingVal = NULL
		}

		stmt = stmt.VALUES(
			UUID(prov.Id()),
			String(prov.Owner().Type().String()),
			ownerFamilyID,
			ownerUserID,
			String(prov.Name()),
			keyVal,
			descVal,
			iconVal,
			urlVal,
			pricingVal,
			TimestampzT(prov.CreatedAt()),
			TimestampzT(prov.UpdatedAt()),
			String(prov.ETag()),
		)
	}

	count, err := r.dbContext.Execute(ctx, stmt)
	if err != nil {
		return err
	}
	if count != int64(len(providers)) {
		return db.ErrMissMatchAffectRow
	}

	// Insert provider labels
	allLabels := collection.SelectMany(providers, func(prov provider.Provider) []struct {
		ProviderID uuid.UUID
		LabelID    uuid.UUID
	} {
		return collection.Select(prov.Labels().Values(), func(labelId uuid.UUID) struct {
			ProviderID uuid.UUID
			LabelID    uuid.UUID
		} {
			return struct {
				ProviderID uuid.UUID
				LabelID    uuid.UUID
			}{ProviderID: prov.Id(), LabelID: labelId}
		})
	})

	if len(allLabels) > 0 {
		labelStmt := ProviderLabels.INSERT(
			ProviderLabels.ProviderID,
			ProviderLabels.LabelID,
		)

		for _, label := range allLabels {
			labelStmt = labelStmt.VALUES(
				UUID(label.ProviderID),
				UUID(label.LabelID),
			)
		}

		if _, err := r.dbContext.Execute(ctx, labelStmt); err != nil {
			return err
		}
	}

	// Insert plans and prices
	for _, prov := range providers {
		if err := r.createPlansAndPrices(ctx, prov.Id(), prov.Plans().Values()); err != nil {
			return err
		}
	}

	return nil
}

func (r ProviderRepository) update(ctx context.Context, dirtyProvider provider.Provider) error {
	if !dirtyProvider.IsDirty() {
		return nil
	}

	// Update provider
	var ownerFamilyID StringExpression
	var ownerUserID StringExpression

	switch dirtyProvider.Owner().Type() {
	case auth.PersonalOwnerType:
		ownerFamilyID = StringExp(NULL)
		ownerUserID = String(dirtyProvider.Owner().UserId())
	case auth.FamilyOwnerType:
		ownerFamilyID = StringExp(UUID(dirtyProvider.Owner().FamilyId()))
		ownerUserID = StringExp(NULL)
	default:
		ownerFamilyID = StringExp(NULL)
		ownerUserID = StringExp(NULL)
	}

	var keyVal StringExpression
	if dirtyProvider.Key() != nil {
		keyVal = String(*dirtyProvider.Key())
	} else {
		keyVal = StringExp(NULL)
	}

	var descVal StringExpression
	if dirtyProvider.Description() != nil {
		descVal = String(*dirtyProvider.Description())
	} else {
		descVal = StringExp(NULL)
	}

	var iconVal StringExpression
	if dirtyProvider.IconUrl() != nil {
		iconVal = String(*dirtyProvider.IconUrl())
	} else {
		iconVal = StringExp(NULL)
	}

	var urlVal StringExpression
	if dirtyProvider.Url() != nil {
		urlVal = String(*dirtyProvider.Url())
	} else {
		urlVal = StringExp(NULL)
	}

	var pricingVal StringExpression
	if dirtyProvider.PricingPageUrl() != nil {
		pricingVal = String(*dirtyProvider.PricingPageUrl())
	} else {
		pricingVal = StringExp(NULL)
	}

	stmt := Providers.UPDATE().
		SET(
			Providers.OwnerType.SET(String(dirtyProvider.Owner().Type().String())),
			Providers.OwnerFamilyID.SET(ownerFamilyID),
			Providers.OwnerUserID.SET(ownerUserID),
			Providers.Name.SET(String(dirtyProvider.Name())),
			Providers.Key.SET(keyVal),
			Providers.Description.SET(descVal),
			Providers.IconURL.SET(iconVal),
			Providers.URL.SET(urlVal),
			Providers.PricingPageURL.SET(pricingVal),
			Providers.UpdatedAt.SET(TimestampzT(dirtyProvider.UpdatedAt())),
		).
		WHERE(Providers.ID.EQ(UUID(dirtyProvider.Id())))

	count, err := r.dbContext.Execute(ctx, stmt)
	if err != nil {
		return err
	}
	if count == 0 {
		return db.ErrMissMatchAffectRow
	}

	// Handle label changes
	if dirtyProvider.Labels().HasChanges() {
		if err := r.saveTrackedLabelsWithJet(ctx, dirtyProvider.Id(), dirtyProvider.Labels()); err != nil {
			return err
		}
		dirtyProvider.Labels().ClearChanges()
	}

	// Handle plan changes
	if dirtyProvider.Plans().HasChanges() {
		if err := r.saveTrackedPlansWithJet(ctx, dirtyProvider.Id(), dirtyProvider.Plans()); err != nil {
			return err
		}
		dirtyProvider.Plans().ClearChanges()
	}

	return nil
}

// Helper methods for handling tracked collections with go-jet

func (r ProviderRepository) saveTrackedLabelsWithJet(ctx context.Context, providerId uuid.UUID,
	labels *slicesx.Tracked[uuid.UUID]) error {
	// Handle new labels
	newLabels := labels.Added()
	if len(newLabels) > 0 {
		stmt := ProviderLabels.INSERT(
			ProviderLabels.ProviderID,
			ProviderLabels.LabelID,
		)

		for _, labelId := range newLabels {
			stmt = stmt.VALUES(
				UUID(providerId),
				UUID(labelId),
			)
		}

		_, err := r.dbContext.Execute(ctx, stmt)
		if err != nil {
			return err
		}
	}

	// Handle deleted labels
	deletedLabels := labels.Removed()
	for _, labelId := range deletedLabels {
		stmt := ProviderLabels.DELETE().
			WHERE(
				ProviderLabels.ProviderID.EQ(UUID(providerId)).
					AND(ProviderLabels.LabelID.EQ(UUID(labelId))),
			)

		_, err := r.dbContext.Execute(ctx, stmt)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r ProviderRepository) saveTrackedPlansWithJet(ctx context.Context, providerId uuid.UUID,
	plans *slicesx.Tracked[provider.Plan]) error {
	// Handle new plans
	newPlans := plans.Added()
	if len(newPlans) > 0 {
		if err := r.createPlansAndPrices(ctx, providerId, newPlans); err != nil {
			return err
		}
	}

	// Handle updated plans
	updatedPlans := plans.Updated()
	for _, plan := range updatedPlans {
		if plan.IsDirty() {
			var descVal StringExpression
			if plan.Description() != nil {
				descVal = String(*plan.Description())
			} else {
				descVal = StringExp(NULL)
			}

			stmt := ProviderPlans.UPDATE().
				SET(
					ProviderPlans.Name.SET(String(plan.Name())),
					ProviderPlans.Description.SET(descVal),
					ProviderPlans.UpdatedAt.SET(TimestampzT(plan.UpdatedAt())),
					ProviderPlans.Etag.SET(String(plan.ETag())),
				).
				WHERE(ProviderPlans.ID.EQ(UUID(plan.Id())))

			count, err := r.dbContext.Execute(ctx, stmt)
			if err != nil {
				return err
			}
			if count == 0 {
				return db.ErrMissMatchAffectRow
			}
		}

		if plan.Prices().HasChanges() {
			if err := r.saveTrackedPricesWithJet(ctx, plan.Id(), plan.Prices()); err != nil {
				return err
			}
		}
	}

	// Handle deleted plans
	deletedPlans := plans.Removed()
	for _, plan := range deletedPlans {
		stmt := ProviderPlans.DELETE().
			WHERE(ProviderPlans.ID.EQ(UUID(plan.Id())))

		_, err := r.dbContext.Execute(ctx, stmt)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r ProviderRepository) saveTrackedPricesWithJet(ctx context.Context, planId uuid.UUID,
	prices *slicesx.Tracked[provider.Price]) error {
	// Handle new prices
	newPrices := prices.Added()
	if len(newPrices) > 0 {
		stmt := ProviderPrices.INSERT(
			ProviderPrices.ID,
			ProviderPrices.PlanID,
			ProviderPrices.Currency,
			ProviderPrices.Amount,
			ProviderPrices.StartDate,
			ProviderPrices.EndDate,
			ProviderPrices.CreatedAt,
			ProviderPrices.UpdatedAt,
			ProviderPrices.Etag,
		)

		for _, price := range newPrices {
			var endDateVal Expression
			if price.EndDate() != nil {
				endDateVal = TimestampzT(*price.EndDate())
			} else {
				endDateVal = NULL
			}

			stmt = stmt.VALUES(
				UUID(price.Id()),
				UUID(planId),
				String(price.Currency().String()),
				Float(price.Amount()),
				TimestampzT(price.StartDate()),
				endDateVal,
				TimestampzT(price.CreatedAt()),
				TimestampzT(price.UpdatedAt()),
				String(price.ETag()),
			)
		}

		_, err := r.dbContext.Execute(ctx, stmt)
		if err != nil {
			return err
		}
	}

	// Handle updated prices
	updatedPrices := prices.Updated()
	for _, price := range updatedPrices {
		if price.IsDirty() {
			var endDateVal TimestampzExpression
			if price.EndDate() != nil {
				endDateVal = TimestampzT(*price.EndDate())
			} else {
				endDateVal = TimestampzExp(NULL)
			}

			stmt := ProviderPrices.UPDATE().
				SET(
					ProviderPrices.Currency.SET(String(price.Currency().String())),
					ProviderPrices.Amount.SET(Float(price.Amount())),
					ProviderPrices.StartDate.SET(TimestampzT(price.StartDate())),
					ProviderPrices.EndDate.SET(endDateVal),
					ProviderPrices.UpdatedAt.SET(TimestampzT(price.UpdatedAt())),
					ProviderPrices.Etag.SET(String(price.ETag())),
				).
				WHERE(ProviderPrices.ID.EQ(UUID(price.Id())))

			count, err := r.dbContext.Execute(ctx, stmt)
			if err != nil {
				return err
			}
			if count == 0 {
				return db.ErrMissMatchAffectRow
			}
		}
	}

	// Handle deleted prices
	deletedPrices := prices.Removed()
	for _, price := range deletedPrices {
		stmt := ProviderPrices.DELETE().
			WHERE(ProviderPrices.ID.EQ(UUID(price.Id())))

		_, err := r.dbContext.Execute(ctx, stmt)
		if err != nil {
			return err
		}
	}

	prices.ClearChanges()
	return nil
}

func (r ProviderRepository) createPlansAndPrices(ctx context.Context, providerId uuid.UUID,
	plans []provider.Plan) error {
	if len(plans) == 0 {
		return nil
	}

	// Insert plans
	planStmt := ProviderPlans.INSERT(
		ProviderPlans.ID,
		ProviderPlans.ProviderID,
		ProviderPlans.Name,
		ProviderPlans.Description,
		ProviderPlans.CreatedAt,
		ProviderPlans.UpdatedAt,
		ProviderPlans.Etag,
	)

	for _, plan := range plans {
		var descVal Expression
		if plan.Description() != nil {
			descVal = String(*plan.Description())
		} else {
			descVal = NULL
		}

		planStmt = planStmt.VALUES(
			UUID(plan.Id()),
			UUID(providerId),
			String(plan.Name()),
			descVal,
			TimestampzT(plan.CreatedAt()),
			TimestampzT(plan.UpdatedAt()),
			String(plan.ETag()),
		)
	}

	planCount, err := r.dbContext.Execute(ctx, planStmt)
	if err != nil {
		return err
	}
	if planCount != int64(len(plans)) {
		return db.ErrMissMatchAffectRow
	}

	// Insert prices for all plans
	allPrices := collection.SelectMany(plans, func(plan provider.Plan) []struct {
		PlanID uuid.UUID
		Price  provider.Price
	} {
		return collection.Select(plan.Prices().Values(), func(price provider.Price) struct {
			PlanID uuid.UUID
			Price  provider.Price
		} {
			return struct {
				PlanID uuid.UUID
				Price  provider.Price
			}{PlanID: plan.Id(), Price: price}
		})
	})

	if len(allPrices) > 0 {
		priceStmt := ProviderPrices.INSERT(
			ProviderPrices.ID,
			ProviderPrices.PlanID,
			ProviderPrices.Currency,
			ProviderPrices.Amount,
			ProviderPrices.StartDate,
			ProviderPrices.EndDate,
			ProviderPrices.CreatedAt,
			ProviderPrices.UpdatedAt,
			ProviderPrices.Etag,
		)

		for _, priceData := range allPrices {
			price := priceData.Price
			var endDateVal Expression
			if price.EndDate() != nil {
				endDateVal = TimestampzT(*price.EndDate())
			} else {
				endDateVal = NULL
			}

			priceStmt = priceStmt.VALUES(
				UUID(price.Id()),
				UUID(priceData.PlanID),
				String(price.Currency().String()),
				Float(price.Amount()),
				TimestampzT(price.StartDate()),
				endDateVal,
				TimestampzT(price.CreatedAt()),
				TimestampzT(price.UpdatedAt()),
				String(price.ETag()),
			)
		}

		priceCount, err := r.dbContext.Execute(ctx, priceStmt)
		if err != nil {
			return err
		}
		if priceCount != int64(len(allPrices)) {
			return db.ErrMissMatchAffectRow
		}
	}

	return nil
}
