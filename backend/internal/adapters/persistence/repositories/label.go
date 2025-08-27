package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/adapters/persistence/db"
	"github.com/mistribe/subtracker/internal/adapters/persistence/db/jet/app/public/model"
	. "github.com/mistribe/subtracker/internal/adapters/persistence/db/jet/app/public/table"
	"github.com/mistribe/subtracker/internal/adapters/persistence/db/models"
	"github.com/mistribe/subtracker/internal/domain/auth"
	"github.com/mistribe/subtracker/internal/domain/label"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/slicesx"

	. "github.com/go-jet/jet/v2/postgres"
)

type LabelRepository struct {
	dbContext *db.Context
}

func NewLabelRepository(dbContext *db.Context) ports.LabelRepository {
	return &LabelRepository{
		dbContext: dbContext,
	}
}

func (r LabelRepository) GetById(ctx context.Context, labelId uuid.UUID) (label.Label, error) {
	stmt := SELECT(Labels.AllColumns).
		FROM(Labels).
		WHERE(Labels.ID.EQ(UUID(labelId)))

	var row model.Labels
	if err := r.dbContext.Query(ctx, stmt, &row); err != nil {
		return nil, err
	}

	if row.ID == uuid.Nil {
		return nil, nil
	}

	lbl := models.CreateLabelFromModel(row)
	return lbl, nil
}

func (r LabelRepository) GetByIdForUser(ctx context.Context, userId string, labelId uuid.UUID) (label.Label, error) {
	stmt := SELECT(Labels.AllColumns).
		FROM(
			Labels.
				LEFT_JOIN(Families, Families.ID.EQ(Labels.OwnerFamilyID)).
				LEFT_JOIN(FamilyMembers, FamilyMembers.FamilyID.EQ(Families.ID)),
		).
		WHERE(
			Labels.ID.EQ(UUID(labelId)).
				AND(
					Labels.OwnerType.EQ(String("system")).
						OR(Labels.OwnerType.EQ(String("personal")).AND(Labels.OwnerUserID.EQ(String(userId)))).
						OR(Labels.OwnerType.EQ(String("family")).AND(FamilyMembers.UserID.EQ(String(userId)))),
				),
		)

	var row model.Labels
	if err := r.dbContext.Query(ctx, stmt, &row); err != nil {
		return nil, err
	}

	if row.ID == uuid.Nil {
		return nil, nil
	}

	lbl := models.CreateLabelFromModel(row)
	return lbl, nil
}

func (r LabelRepository) GetAll(ctx context.Context, userId string, parameters ports.LabelQueryParameters) (
	[]label.Label,
	int64,
	error) {

	baseStmt := SELECT(
		Labels.AllColumns,
		COUNT(STAR).OVER().AS("total_count"),
	).
		FROM(
			Labels.
				LEFT_JOIN(Families, Families.ID.EQ(Labels.OwnerFamilyID)).
				LEFT_JOIN(FamilyMembers, FamilyMembers.FamilyID.EQ(Families.ID)),
		).
		WHERE(
			Labels.OwnerType.EQ(String("system")).
				OR(Labels.OwnerType.EQ(String("personal")).AND(Labels.OwnerUserID.EQ(String(userId)))).
				OR(Labels.OwnerType.EQ(String("family")).AND(FamilyMembers.UserID.EQ(String(userId)))),
		)

	if parameters.SearchText != "" {
		baseStmt = baseStmt.WHERE(Labels.Name.LIKE(String(fmt.Sprintf("%%%s%%", parameters.SearchText))))
	}

	stmt := baseStmt.
		LIMIT(parameters.Limit).
		OFFSET(parameters.Offset)

	var rows []struct {
		Labels     model.Labels `json:"labels"`
		TotalCount int64        `json:"total_count"`
	}

	if err := r.dbContext.Query(ctx, stmt, &rows); err != nil {
		return nil, 0, err
	}
	if len(rows) == 0 {
		return nil, 0, nil
	}

	totalCount := rows[0].TotalCount
	labels := slicesx.Select(rows, func(row struct {
		Labels     model.Labels `json:"labels"`
		TotalCount int64        `json:"total_count"`
	}) label.Label {
		return models.CreateLabelFromModel(row.Labels)
	})

	return labels, totalCount, nil
}

func (r LabelRepository) GetSystemLabels(ctx context.Context) ([]label.Label, error) {
	stmt := SELECT(Labels.AllColumns).
		FROM(Labels).
		WHERE(
			Labels.OwnerFamilyID.IS_NULL().
				AND(Labels.OwnerUserID.IS_NULL()).
				AND(Labels.OwnerType.EQ(String("system"))),
		)

	var rows []model.Labels
	if err := r.dbContext.Query(ctx, stmt, &rows); err != nil {
		return nil, err
	}

	labels := slicesx.Select(rows, func(row model.Labels) label.Label {
		return models.CreateLabelFromModel(row)
	})
	return labels, nil
}

func (r LabelRepository) Save(ctx context.Context, labels ...label.Label) error {
	var newLabels []label.Label
	for _, lbl := range labels {
		if !lbl.IsExists() {
			newLabels = append(newLabels, lbl)
		} else {
			if err := r.update(ctx, lbl); err != nil {
				return err
			}
		}
	}
	if len(newLabels) > 0 {
		if err := r.create(ctx, newLabels); err != nil {
			return err
		}
	}

	for _, lbl := range labels {
		lbl.Clean()
	}

	return nil
}

func (r LabelRepository) create(ctx context.Context, labels []label.Label) error {
	if len(labels) == 0 {
		return nil
	}

	stmt := Labels.INSERT(
		Labels.ID,
		Labels.OwnerType,
		Labels.OwnerFamilyID,
		Labels.OwnerUserID,
		Labels.Name,
		Labels.Key,
		Labels.Color,
		Labels.CreatedAt,
		Labels.UpdatedAt,
		Labels.Etag,
	)

	for _, lbl := range labels {
		var ownerFamilyID Expression
		var ownerUserID Expression

		switch lbl.Owner().Type() {
		case auth.PersonalOwnerType:
			ownerFamilyID = NULL
			ownerUserID = String(lbl.Owner().UserId())
		case auth.FamilyOwnerType:
			ownerFamilyID = UUID(lbl.Owner().FamilyId())
			ownerUserID = NULL
		default:
			ownerFamilyID = NULL
			ownerUserID = NULL
		}

		var keyVal Expression
		if lbl.Key() != nil {
			keyVal = String(*lbl.Key())
		} else {
			keyVal = NULL
		}

		stmt = stmt.VALUES(
			UUID(lbl.Id()),
			String(lbl.Owner().Type().String()),
			ownerFamilyID,
			ownerUserID,
			String(lbl.Name()),
			keyVal,
			String(lbl.Color()),
			TimestampzT(lbl.CreatedAt()),
			TimestampzT(lbl.UpdatedAt()),
			String(lbl.ETag()),
		)
	}

	count, err := r.dbContext.Execute(ctx, stmt)
	if err != nil {
		return err
	}
	if count != int64(len(labels)) {
		return db.ErrMissMatchAffectRow
	}
	return nil
}
func (r LabelRepository) update(ctx context.Context, lbl label.Label) error {
	if !lbl.IsDirty() {
		return nil
	}

	var ownerFamilyID StringExpression
	var ownerUserID StringExpression

	switch lbl.Owner().Type() {
	case auth.PersonalOwnerType:
		ownerFamilyID = StringExp(NULL)
		ownerUserID = String(lbl.Owner().UserId())
	case auth.FamilyOwnerType:
		ownerFamilyID = StringExp(UUID(lbl.Owner().FamilyId()))
		ownerUserID = StringExp(NULL)
	default:
		ownerFamilyID = StringExp(NULL)
		ownerUserID = StringExp(NULL)
	}

	var keyVal StringExpression
	if lbl.Key() != nil {
		keyVal = String(*lbl.Key())
	} else {
		keyVal = StringExp(NULL)
	}

	stmt := Labels.UPDATE().
		SET(
			Labels.OwnerType.SET(String(lbl.Owner().Type().String())),
			Labels.OwnerFamilyID.SET(ownerFamilyID),
			Labels.OwnerUserID.SET(ownerUserID),
			Labels.Name.SET(String(lbl.Name())),
			Labels.Key.SET(keyVal),
			Labels.Color.SET(String(lbl.Color())),
			Labels.UpdatedAt.SET(TimestampzT(lbl.UpdatedAt())),
			Labels.Etag.SET(String(lbl.ETag())),
		).
		WHERE(Labels.ID.EQ(UUID(lbl.Id())))

	count, err := r.dbContext.Execute(ctx, stmt)
	if err != nil {
		return err
	}
	if count == 0 {
		return db.ErrMissMatchAffectRow
	}
	return nil
}

func (r LabelRepository) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	stmt := Labels.DELETE().
		WHERE(Labels.ID.EQ(UUID(id)))

	count, err := r.dbContext.Execute(ctx, stmt)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r LabelRepository) Exists(ctx context.Context, ids ...uuid.UUID) (bool, error) {
	if len(ids) == 0 {
		return true, nil
	}

	vals := make([]Expression, len(ids))
	for i, id := range ids {
		vals[i] = UUID(id)
	}

	stmt := SELECT(COUNT(Labels.ID).AS("count")).
		FROM(Labels).
		WHERE(Labels.ID.IN(vals...))

	var row struct {
		Count int
	}

	if err := r.dbContext.Query(ctx, stmt, &row); err != nil {
		return false, err
	}

	return row.Count == len(ids), nil
}
