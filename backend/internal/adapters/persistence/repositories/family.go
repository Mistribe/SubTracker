package repositories

import (
	"context"

	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/adapters/persistence/db"
	"github.com/mistribe/subtracker/internal/adapters/persistence/db/jet/app/public/model"
	. "github.com/mistribe/subtracker/internal/adapters/persistence/db/jet/app/public/table"
	"github.com/mistribe/subtracker/internal/adapters/persistence/db/models"
	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/slicesx"

	. "github.com/go-jet/jet/v2/postgres"
)

type FamilyRepository struct {
	dbContext *db.Context
}

func NewFamilyRepository(repository *db.Context) ports.FamilyRepository {
	return &FamilyRepository{
		dbContext: repository,
	}
}

func (r FamilyRepository) GetAccountFamily(ctx context.Context, userId types.UserID) (family.Family, error) {
	userFamilies := SELECT(
		Families.AllColumns,
	).
		FROM(
			Families.
				INNER_JOIN(FamilyMembers, FamilyMembers.FamilyID.EQ(Families.ID)),
		).
		WHERE(
			FamilyMembers.UserID.EQ(String(userId.String())),
		).
		ORDER_BY(Families.ID).
		LIMIT(1).
		AsTable("uf")

	rFamilyId := Families.ID.From(userFamilies)

	stmt := SELECT(
		userFamilies.AllColumns(),
		FamilyMembers.AllColumns,
	).
		FROM(
			userFamilies.
				LEFT_JOIN(FamilyMembers, FamilyMembers.FamilyID.EQ(rFamilyId)),
		)

	var rows []struct {
		Family        model.Families       `json:"families"`
		FamilyMembers *model.FamilyMembers `json:"family_members"`
	}

	if err := r.dbContext.Query(ctx, stmt, &rows); err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, nil
	}

	// Convert rows to FamilyWithMembers format
	familyData := models.FamilyWithMembers{
		Families:      rows[0].Family,
		FamilyMembers: make([]model.FamilyMembers, 0),
	}

	for _, row := range rows {
		if row.FamilyMembers != nil && row.FamilyMembers.ID != uuid.Nil {
			familyData.FamilyMembers = append(familyData.FamilyMembers, *row.FamilyMembers)
		}
	}

	return models.CreateFamilyWithMembersFromModel(familyData)

}

func (r FamilyRepository) GetById(ctx context.Context, familyID types.FamilyID) (family.Family, error) {
	stmt := SELECT(Families.AllColumns, FamilyMembers.AllColumns).
		FROM(Families.LEFT_JOIN(FamilyMembers, FamilyMembers.FamilyID.EQ(Families.ID))).
		WHERE(Families.ID.EQ(UUID(familyID)))

	var rows []struct {
		Family        model.Families       `json:"families"`
		FamilyMembers *model.FamilyMembers `json:"family_members"`
	}

	if err := r.dbContext.Query(ctx, stmt, &rows); err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, nil
	}

	// Convert rows to FamilyWithMembers format
	familyData := models.FamilyWithMembers{
		Families:      rows[0].Family,
		FamilyMembers: make([]model.FamilyMembers, 0),
	}

	for _, row := range rows {
		if row.FamilyMembers != nil && row.FamilyMembers.ID != uuid.Nil {
			familyData.FamilyMembers = append(familyData.FamilyMembers, *row.FamilyMembers)
		}
	}

	return models.CreateFamilyWithMembersFromModel(familyData)
}

func (r FamilyRepository) Save(ctx context.Context, families ...family.Family) error {
	var newFamilies []family.Family
	for _, fam := range families {
		if !fam.IsExists() {
			newFamilies = append(newFamilies, fam)
		} else {
			if err := r.update(ctx, fam); err != nil {
				return err
			}
		}
	}

	if len(newFamilies) > 0 {
		if err := r.create(ctx, newFamilies); err != nil {
			return err
		}
	}

	for _, fam := range families {
		for _, mbr := range fam.Members().Values() {
			mbr.Clean()
		}
		fam.Clean()
	}
	return nil
}

func (r FamilyRepository) create(ctx context.Context, families []family.Family) error {
	if len(families) == 0 {
		return nil
	}

	// Insert families
	stmt := Families.INSERT(
		Families.ID,
		Families.Name,
		Families.OwnerID,
		Families.CreatedAt,
		Families.UpdatedAt,
		Families.Etag,
	)

	for _, fam := range families {
		stmt = stmt.VALUES(
			UUID(fam.Id()),
			String(fam.Name()),
			String(fam.Owner().UserId().String()),
			TimestampzT(fam.CreatedAt()),
			TimestampzT(fam.UpdatedAt()),
			String(fam.ETag()),
		)
	}

	count, err := r.dbContext.Execute(ctx, stmt)
	if err != nil {
		return err
	}
	if count != int64(len(families)) {
		return db.ErrMissMatchAffectRow
	}

	// Insert family members
	var allMembers []family.Member
	for _, fam := range families {
		for _, member := range fam.Members().Values() {
			allMembers = append(allMembers, member)
		}
	}

	if len(allMembers) > 0 {
		memberStmt := FamilyMembers.INSERT(
			FamilyMembers.ID,
			FamilyMembers.FamilyID,
			FamilyMembers.UserID,
			FamilyMembers.Name,
			FamilyMembers.Type,
			FamilyMembers.CreatedAt,
			FamilyMembers.UpdatedAt,
			FamilyMembers.Etag,
		)

		for _, member := range allMembers {
			var userID Expression
			if member.UserId() != nil {
				userID = String((*member.UserId()).String())
			} else {
				userID = NULL
			}

			memberStmt = memberStmt.VALUES(
				UUID(member.Id()),
				UUID(member.FamilyId()),
				userID,
				String(member.Name()),
				String(member.Type().String()),
				TimestampzT(member.CreatedAt()),
				TimestampzT(member.UpdatedAt()),
				String(member.ETag()),
			)
		}

		memberCount, err := r.dbContext.Execute(ctx, memberStmt)
		if err != nil {
			return err
		}
		if memberCount != int64(len(allMembers)) {
			return db.ErrMissMatchAffectRow
		}
	}

	return nil
}

func (r FamilyRepository) update(ctx context.Context, fam family.Family) error {
	if !fam.IsDirty() {
		return nil
	}

	// Update family
	stmt := Families.UPDATE().
		SET(
			Families.Name.SET(String(fam.Name())),
			Families.OwnerID.SET(String(fam.Owner().UserId().String())),
			Families.UpdatedAt.SET(TimestampzT(fam.UpdatedAt())),
			Families.Etag.SET(String(fam.ETag())),
		).
		WHERE(Families.ID.EQ(UUID(fam.Id())))

	count, err := r.dbContext.Execute(ctx, stmt)
	if err != nil {
		return err
	}
	if count == 0 {
		return db.ErrMissMatchAffectRow
	}

	// Handle member changes using go-jet operations
	err = r.saveTrackedSliceWithJet(ctx, fam.Members(), fam.Id())
	if err != nil {
		return err
	}

	fam.Members().ClearChanges()
	return nil
}

// saveTrackedSliceWithJet handles member operations using go-jet
func (r FamilyRepository) saveTrackedSliceWithJet(
	ctx context.Context, members *slicesx.Tracked[family.Member],
	familyId types.FamilyID) error {
	// Handle new members
	newMembers := members.Added()
	if len(newMembers) > 0 {
		stmt := FamilyMembers.INSERT(
			FamilyMembers.ID,
			FamilyMembers.FamilyID,
			FamilyMembers.UserID,
			FamilyMembers.Name,
			FamilyMembers.Type,
			FamilyMembers.InvitationCode,
			FamilyMembers.CreatedAt,
			FamilyMembers.UpdatedAt,
			FamilyMembers.Etag,
		)

		for _, member := range newMembers {
			var userID Expression
			if member.UserId() != nil {
				userID = String((*member.UserId()).String())
			} else {
				userID = NULL
			}
			var invitationCode Expression
			if member.InvitationCode() != nil {
				invitationCode = String(*member.InvitationCode())
			} else {
				invitationCode = NULL
			}

			stmt = stmt.VALUES(
				UUID(member.Id()),
				UUID(familyId),
				userID,
				String(member.Name()),
				String(member.Type().String()),
				invitationCode,
				TimestampzT(member.CreatedAt()),
				TimestampzT(member.UpdatedAt()),
				String(member.ETag()),
			)
		}

		_, err := r.dbContext.Execute(ctx, stmt)
		if err != nil {
			return err
		}
	}

	// Handle updated members
	updatedMembers := members.Updated()
	for _, member := range updatedMembers {
		if member.IsDirty() {
			var userID StringExpression
			if member.UserId() != nil {
				userID = String((*member.UserId()).String())
			} else {
				userID = StringExp(NULL)
			}
			var invitationCode StringExpression
			if member.InvitationCode() != nil {
				invitationCode = String(*member.InvitationCode())
			} else {
				invitationCode = StringExp(NULL)
			}

			stmt := FamilyMembers.UPDATE().
				SET(
					FamilyMembers.FamilyID.SET(UUID(member.FamilyId())),
					FamilyMembers.UserID.SET(userID),
					FamilyMembers.Name.SET(String(member.Name())),
					FamilyMembers.Type.SET(String(member.Type().String())),
					FamilyMembers.InvitationCode.SET(invitationCode),
					FamilyMembers.UpdatedAt.SET(TimestampzT(member.UpdatedAt())),
					FamilyMembers.Etag.SET(String(member.ETag())),
				).
				WHERE(FamilyMembers.ID.EQ(UUID(member.Id())))

			count, err := r.dbContext.Execute(ctx, stmt)
			if err != nil {
				return err
			}
			if count == 0 {
				return db.ErrMissMatchAffectRow
			}
		}
	}

	// Handle deleted members
	deletedMembers := members.Removed()
	for _, member := range deletedMembers {
		stmt := FamilyMembers.DELETE().
			WHERE(FamilyMembers.ID.EQ(UUID(member.Id())))

		_, err := r.dbContext.Execute(ctx, stmt)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r FamilyRepository) Delete(ctx context.Context, familyId types.FamilyID) (bool, error) {
	stmt := Families.DELETE().
		WHERE(Families.ID.EQ(UUID(familyId)))

	count, err := r.dbContext.Execute(ctx, stmt)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r FamilyRepository) MemberExists(ctx context.Context, familyId types.FamilyID,
	members ...types.FamilyMemberID) (bool, error) {
	if len(members) == 0 {
		return true, nil
	}

	vals := make([]Expression, len(members))
	for i, id := range members {
		vals[i] = UUID(id)
	}

	stmt := SELECT(COUNT(FamilyMembers.ID).AS("count")).
		FROM(FamilyMembers).
		WHERE(
			FamilyMembers.FamilyID.EQ(UUID(familyId)).
				AND(FamilyMembers.ID.IN(vals...)),
		)

	var row struct {
		Count int
	}

	if err := r.dbContext.Query(ctx, stmt, &row); err != nil {
		return false, err
	}

	return row.Count == len(members), nil
}

func (r FamilyRepository) Exists(ctx context.Context, ids ...types.FamilyID) (bool, error) {
	if len(ids) == 0 {
		return true, nil
	}

	vals := make([]Expression, len(ids))
	for i, id := range ids {
		vals[i] = UUID(id)
	}

	stmt := SELECT(COUNT(Families.ID).AS("count")).
		FROM(Families).
		WHERE(Families.ID.IN(vals...))

	var row struct {
		Count int
	}

	if err := r.dbContext.Query(ctx, stmt, &row); err != nil {
		return false, err
	}

	return row.Count == len(ids), nil
}

func (r FamilyRepository) IsUserMemberOfFamily(ctx context.Context, familyId types.FamilyID,
	userId types.UserID) (bool, error) {
	stmt := SELECT(COUNT(FamilyMembers.ID).AS("count")).
		FROM(FamilyMembers).
		WHERE(
			FamilyMembers.FamilyID.EQ(UUID(familyId)).
				AND(FamilyMembers.UserID.EQ(String(userId.String()))),
		)

	var row struct {
		Count int
	}

	if err := r.dbContext.Query(ctx, stmt, &row); err != nil {
		return false, err
	}

	return row.Count > 0, nil
}
