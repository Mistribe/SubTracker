package repositories

import (
	"context"
	"fmt"

	"github.com/mistribe/subtracker/internal/adapters/persistence/db"
	. "github.com/mistribe/subtracker/internal/adapters/persistence/db/jet/app/public/table"
	"github.com/mistribe/subtracker/internal/adapters/persistence/db/models"
	"github.com/mistribe/subtracker/internal/domain/provider"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/slicesx"
	"github.com/mistribe/subtracker/pkg/x/herd"

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

func (r ProviderRepository) GetById(ctx context.Context, providerId types.ProviderID) (provider.Provider, error) {
	stmt := SELECT(
		Providers.AllColumns,
		ProviderLabels.LabelID,
		ProviderLabels.ProviderID,
	).
		FROM(
			Providers.
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

func (r ProviderRepository) GetByIdForUser(ctx context.Context, userId types.UserID, providerId types.ProviderID) (
	provider.Provider,
	error) {
	stmt := SELECT(
		Providers.AllColumns,
		ProviderLabels.LabelID,
		ProviderLabels.ProviderID,
	).
		FROM(
			Providers.
				LEFT_JOIN(ProviderLabels, ProviderLabels.ProviderID.EQ(Providers.ID)).
				LEFT_JOIN(Families, Families.ID.EQ(Providers.OwnerFamilyID)).
				LEFT_JOIN(FamilyMembers, FamilyMembers.FamilyID.EQ(Families.ID)),
		).
		WHERE(
			Providers.ID.EQ(UUID(providerId)).
				AND(
					Providers.OwnerType.EQ(String("system")).
						OR(Providers.OwnerType.EQ(String("personal")).AND(Providers.OwnerUserID.EQ(String(userId.String())))).
						OR(Providers.OwnerType.EQ(String("family")).AND(FamilyMembers.UserID.EQ(String(userId.String())))),
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
	// Use a subselect to compute total once based on distinct providers, then join labels
	base := SELECT(
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
		COUNT(Providers.ID).OVER().AS("total_count"),
	).
		FROM(Providers).
		WHERE(Providers.OwnerType.EQ(String("system"))).
		AsTable("bp")

	stmt := SELECT(
		base.AllColumns(),
		ProviderLabels.LabelID,
		ProviderLabels.ProviderID,
	).
		FROM(
			base.
				LEFT_JOIN(ProviderLabels, ProviderLabels.ProviderID.EQ(Providers.ID.From(base))),
		).
		ORDER_BY(Providers.ID.From(base).ASC())

	var rows []models.ProviderRowWithCount
	if err := r.dbContext.Query(ctx, stmt, &rows); err != nil {
		return nil, 0, err
	}
	if len(rows) == 0 {
		return nil, 0, nil
	}

	providers := models.CreateProviderFromJetRowsWithCount(rows)
	if len(providers) == 0 {
		return nil, 0, nil
	}

	totalCount := rows[0].TotalCount
	if totalCount == 0 && len(providers) > 0 {
		totalCount = int64(len(providers))
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
		ProviderLabels.LabelID,
		ProviderLabels.ProviderID,
	).
		FROM(
			pagedProviders.
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
	userId types.UserID,
	parameters ports.ProviderQueryParameters) ([]provider.Provider, int64, error) {

	// Build base access filter
	accessFilter := Providers.OwnerType.EQ(String("system")).
		OR(Providers.OwnerType.EQ(String("personal")).AND(Providers.OwnerUserID.EQ(String(userId.String())))).
		OR(Providers.OwnerType.EQ(String("family")).AND(EXISTS(
			SELECT(FamilyMembers.ID).
				FROM(FamilyMembers).
				WHERE(FamilyMembers.FamilyID.EQ(Providers.OwnerFamilyID).AND(FamilyMembers.UserID.EQ(String(userId.String()))))),
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
		ProviderLabels.LabelID,
		ProviderLabels.ProviderID,
		pagedProviders.AllColumns().Except(Providers.ID),
	).
		FROM(
			pagedProviders.
				INNER_JOIN(Providers, Providers.ID.EQ(Providers.ID.From(pagedProviders))).
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

func (r ProviderRepository) Delete(ctx context.Context, providerId types.ProviderID) (bool, error) {
	stmt := Providers.DELETE().
		WHERE(Providers.ID.EQ(UUID(providerId)))

	count, err := r.dbContext.Execute(ctx, stmt)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r ProviderRepository) Exists(ctx context.Context, ids ...types.ProviderID) (bool, error) {
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
		case types.PersonalOwnerType:
			ownerFamilyID = NULL
			ownerUserID = String(prov.Owner().UserId().String())
		case types.FamilyOwnerType:
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
	allLabels := herd.SelectMany(providers, func(prov provider.Provider) []struct {
		ProviderID types.ProviderID
		LabelID    types.LabelID
	} {
		return herd.Select(prov.Labels().Values(), func(labelId types.LabelID) struct {
			ProviderID types.ProviderID
			LabelID    types.LabelID
		} {
			return struct {
				ProviderID types.ProviderID
				LabelID    types.LabelID
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
	case types.PersonalOwnerType:
		ownerFamilyID = StringExp(NULL)
		ownerUserID = String(dirtyProvider.Owner().UserId().String())
	case types.FamilyOwnerType:
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

	return nil
}

// Helper methods for handling tracked collections with go-jet

func (r ProviderRepository) saveTrackedLabelsWithJet(
	ctx context.Context, providerId types.ProviderID,
	labels *slicesx.Tracked[types.LabelID]) error {
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
