package persistence

import (
	"context"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/internal/domain/entity"
	"github.com/oleexo/subtracker/internal/domain/family"
	"github.com/oleexo/subtracker/internal/infrastructure/persistence/sql"
	"github.com/oleexo/subtracker/pkg/slicesx"
)

type FamilyRepository struct {
	dbContext *DatabaseContext
}

func NewFamilyRepository(repository *DatabaseContext) family.Repository {
	return &FamilyRepository{
		dbContext: repository,
	}
}

func (r FamilyRepository) GetById(ctx context.Context, id uuid.UUID) (family.Family, error) {
	response, err := r.dbContext.GetQueries(ctx).GetFamilyById(ctx, id)
	if err != nil {
		return nil, err
	}
	if len(response) == 0 {
		return nil, nil
	}
	families := createFamilyFromSqlcRows(response,
		func(row sql.GetFamilyByIdRow) sql.Family {
			return row.Family
		}, func(row sql.GetFamilyByIdRow) sql.FamilyMember {
			return row.FamilyMember
		},
	)
	if len(families) == 0 {
		return nil, nil
	}
	return families[0], nil
}

func (r FamilyRepository) GetAll(ctx context.Context, parameters entity.QueryParameters) ([]family.Family, int64,
	error) {
	userId, ok := auth.GetUserIdFromContext(ctx)
	if !ok {
		return []family.Family{}, 0, nil
	}

	response, count, err := r.dbContext.GetQueries(ctx).
		GetFamiliesForUser(ctx,
			&userId,
			parameters.Limit,
			parameters.Offset,
		)
	if err != nil {
		return nil, 0, err
	}
	if len(response) == 0 {
		return nil, 0, nil
	}
	families := createFamilyFromSqlcRows(response,
		func(row sql.FamilyRow) sql.Family {
			return row.Family
		},
		func(row sql.FamilyRow) *sql.FamilyMember {
			return row.Member
		},
	)
	if len(families) == 0 {
		return nil, 0, nil
	}
	return families, count, nil
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
	familyArgs := slicesx.Select(families, func(family family.Family) sql.CreateFamiliesParams {
		return sql.CreateFamiliesParams{
			ID:        family.Id(),
			Name:      family.Name(),
			OwnerID:   family.OwnerId(),
			CreatedAt: family.CreatedAt(),
			UpdatedAt: family.UpdatedAt(),
			Etag:      family.ETag(),
		}
	})
	_, err := r.dbContext.GetQueries(ctx).CreateFamilies(ctx, familyArgs)
	if err != nil {
		return err
	}

	memberArgs := slicesx.SelectMany(families, func(fam family.Family) []sql.CreateFamilyMembersParams {
		return slicesx.Select(fam.Members().Values(), func(member family.Member) sql.CreateFamilyMembersParams {
			return sql.CreateFamilyMembersParams{
				ID:        member.Id(),
				FamilyID:  fam.Id(),
				UserID:    member.UserId(),
				Name:      member.Name(),
				Type:      member.Type().String(),
				CreatedAt: member.CreatedAt(),
				UpdatedAt: member.UpdatedAt(),
				Etag:      member.ETag(),
			}
		})
	})

	if len(memberArgs) > 0 {
		_, err = r.dbContext.GetQueries(ctx).CreateFamilyMembers(ctx, memberArgs)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r FamilyRepository) update(ctx context.Context, fam family.Family) error {
	if !fam.IsDirty() {
		return nil
	}
	err := r.dbContext.GetQueries(ctx).UpdateFamily(ctx, sql.UpdateFamilyParams{
		ID:        fam.Id(),
		Name:      fam.Name(),
		OwnerID:   fam.OwnerId(),
		UpdatedAt: fam.UpdatedAt(),
	})
	if err != nil {
		return err
	}

	err = saveTrackedSlice(ctx,
		r.dbContext,
		fam.Members(),
		func(ctx context.Context, queries *sql.Queries, members []family.Member) error {
			args := slicesx.Select(members, func(member family.Member) sql.CreateFamilyMembersParams {
				return sql.CreateFamilyMembersParams{
					ID:        member.Id(),
					FamilyID:  fam.Id(),
					UserID:    member.UserId(),
					Name:      member.Name(),
					Type:      member.Type().String(),
					CreatedAt: member.CreatedAt(),
					UpdatedAt: member.UpdatedAt(),
					Etag:      member.ETag(),
				}
			})
			_, qErr := queries.CreateFamilyMembers(ctx, args)
			return qErr
		},
		func(ctx context.Context, queries *sql.Queries, member family.Member) error {
			return queries.UpdateFamilyMember(ctx, sql.UpdateFamilyMemberParams{
				ID:        member.Id(),
				FamilyID:  member.FamilyId(),
				UserID:    member.UserId(),
				Name:      member.Name(),
				Type:      member.Type().String(),
				UpdatedAt: member.UpdatedAt(),
				Etag:      member.ETag(),
			})
		},
		func(ctx context.Context, queries *sql.Queries, member family.Member) error {
			return queries.DeleteFamilyMember(ctx, member.Id())
		},
	)

	if err != nil {
		return err
	}

	fam.Members().ClearChanges()
	return err
}

func (r FamilyRepository) Delete(ctx context.Context, familyId uuid.UUID) (bool, error) {
	err := r.dbContext.GetQueries(ctx).DeleteFamily(ctx, familyId)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r FamilyRepository) MemberExists(ctx context.Context, familyId uuid.UUID, members ...uuid.UUID) (bool, error) {
	if len(members) == 0 {
		return true, nil
	}
	count, err := r.dbContext.GetQueries(ctx).IsMemberExists(ctx, sql.IsMemberExistsParams{
		FamilyID: familyId,
		Column2:  members,
	})
	if err != nil {
		return false, err
	}
	return count > int64(len(members)), nil
}

func (r FamilyRepository) Exists(ctx context.Context, ids ...uuid.UUID) (bool, error) {
	if len(ids) == 0 {
		return true, nil
	}

	count, err := r.dbContext.GetQueries(ctx).IsFamilyExists(ctx, ids)
	if err != nil {
		return false, err
	}
	return count > int64(len(ids)), nil
}

func (r FamilyRepository) IsUserMemberOfFamily(ctx context.Context, familyId uuid.UUID, userId string) (bool, error) {
	count, err := r.dbContext.GetQueries(ctx).IsMemberOfFamily(ctx, sql.IsMemberOfFamilyParams{
		FamilyID: familyId,
		UserID:   &userId,
	})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
