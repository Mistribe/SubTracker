package sql

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type FamilyRow struct {
	Family Family
	Member *FamilyMember
}

func (q *Queries) newFamilyMember(id *uuid.UUID,
	name *string,
	familyID *uuid.UUID,
	userID *string,
	memberType *string,
	createdAt *time.Time,
	updatedAt *time.Time,
	etag *string) *FamilyMember {
	if id == nil {
		return nil
	}

	return &FamilyMember{
		ID:        *id,
		Name:      *name,
		FamilyID:  *familyID,
		UserID:    userID,
		Type:      *memberType,
		CreatedAt: *createdAt,
		UpdatedAt: *updatedAt,
		Etag:      *etag,
	}
}

func (q *Queries) GetFamilyById(ctx context.Context, id uuid.UUID) ([]getFamilyByIdRow, error) {
	rows, err := q.getFamilyById(ctx, id)
	if err != nil {
		return nil, nil
	}
}

func (q *Queries) GetFamiliesForUser(ctx context.Context,
	userId *string,
	limit, offset int32) ([]FamilyRow, int64, error) {
	rows, err := q.getFamiliesForUser(ctx, getFamiliesForUserParams{
		UserID: userId,
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
	results := make([]FamilyRow, len(rows))
	for i, row := range rows {
		member := q.newFamilyMember(
			row.FamilyMembersID,
			row.FamilyMembersName,
			row.FamilyMembersFamilyID,
			row.FamilyMembersUserID,
			row.FamilyMembersType,
			row.FamilyMembersCreatedAt,
			row.FamilyMembersUpdatedAt,
			row.FamilyMembersEtag,
		)
		results[i] = FamilyRow{
			Family: Family{
				ID:        row.FamiliesID,
				Name:      row.FamiliesName,
				OwnerID:   row.FamiliesOwnerID,
				CreatedAt: row.FamiliesCreatedAt,
				UpdatedAt: row.FamiliesUpdatedAt,
				Etag:      row.FamiliesEtag,
			},
			Member: member,
		}
	}
	return results, totalCount, nil
}
