package persistence

import (
	"context"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/internal/domain/family"
	"github.com/oleexo/subtracker/internal/domain/label"
	"github.com/oleexo/subtracker/internal/infrastructure/persistence/sql"
	"github.com/oleexo/subtracker/pkg/slicesx"
)

type LabelRepository struct {
	dbContext   *DatabaseContext
	authService auth.Service
}

func NewLabelRepository(dbContext *DatabaseContext,
	authService auth.Service) label.Repository {
	return &LabelRepository{
		dbContext:   dbContext,
		authService: authService,
	}
}

func (r LabelRepository) GetById(ctx context.Context, labelId uuid.UUID) (label.Label, error) {
	response, err := r.dbContext.GetQueries(ctx).GetLabelById(ctx, labelId)
	if err != nil {
		return nil, err
	}

	lbl := createLabelFromSqlc(response)
	lbl.Clean()
	return lbl, nil
}

func (r LabelRepository) GetAll(ctx context.Context, parameters label.QueryParameters) ([]label.Label, int64, error) {
	userId, ok := auth.GetUserIdFromContext(ctx)
	if !ok {
		return nil, 0, nil
	}

	var families []uuid.UUID
	if parameters.Owners.Contains(auth.FamilyOwnerType) {
		if parameters.FamilyId != nil {
			if r.authService.IsInFamily(ctx, *parameters.FamilyId) {
				return nil, 0, family.ErrFamilyNotFound
			}
			families = append(families, *parameters.FamilyId)
		} else {
			families = r.authService.MustGetFamilies(ctx)
		}
	}

	response, err := r.dbContext.GetQueries(ctx).GetLabels(ctx, sql.GetLabelsParams{
		Column1:     parameters.Owners.Strings(),
		OwnerUserID: &userId,
		Column3:     families,
		Limit:       parameters.Limit,
		Offset:      parameters.Offset,
	})
	if err != nil {
		return nil, 0, err
	}

	if len(response) == 0 {
		return nil, 0, nil
	}

	labels := slicesx.Select(response, func(model sql.GetLabelsRow) label.Label {
		return createLabelFromSqlc(model.Label)
	})

	return labels, response[0].TotalCount, nil
}

func (r LabelRepository) GetSystemLabels(ctx context.Context) ([]label.Label, error) {
	response, err := r.dbContext.GetQueries(ctx).GetSystemLabels(ctx)
	if err != nil {
		return nil, err
	}
	labels := slicesx.Select(response, func(model sql.Label) label.Label {
		return createLabelFromSqlc(model)
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
	args := slicesx.Select(labels, func(lbl label.Label) sql.CreateLabelsParams {
		params := sql.CreateLabelsParams{
			ID:            lbl.Id(),
			OwnerType:     lbl.Owner().Type().String(),
			OwnerFamilyID: nil,
			OwnerUserID:   nil,
			Name:          lbl.Name(),
			Key:           lbl.Key(),
			Color:         lbl.Color(),
			CreatedAt:     lbl.CreatedAt(),
			UpdatedAt:     lbl.UpdatedAt(),
			Etag:          lbl.ETag(),
		}
		switch lbl.Owner().Type() {
		case auth.PersonalOwnerType:
			userId := lbl.Owner().UserId()
			params.OwnerUserID = &userId
		case auth.FamilyOwnerType:
			familyId := lbl.Owner().FamilyId()
			params.OwnerFamilyID = &familyId
		}
		return params
	})
	_, err := r.dbContext.GetQueries(ctx).CreateLabels(ctx, args)
	if err != nil {
		return err
	}
	return nil
}
func (r LabelRepository) update(ctx context.Context, lbl label.Label) error {
	params := sql.UpdateLabelParams{
		ID:        lbl.Id(),
		OwnerType: lbl.Owner().Type().String(),
		Name:      lbl.Name(),
		Key:       lbl.Key(),
		Color:     lbl.Color(),
		UpdatedAt: lbl.UpdatedAt(),
		Etag:      lbl.ETag(),
	}

	switch lbl.Owner().Type() {
	case auth.PersonalOwnerType:
		userId := lbl.Owner().UserId()
		params.OwnerUserID = &userId
	case auth.FamilyOwnerType:
		familyId := lbl.Owner().FamilyId()
		params.OwnerFamilyID = &familyId
	}

	return r.dbContext.GetQueries(ctx).UpdateLabel(ctx, params)
}

func (r LabelRepository) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	err := r.dbContext.GetQueries(ctx).DeleteLabel(ctx, id)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r LabelRepository) Exists(ctx context.Context, ids ...uuid.UUID) (bool, error) {
	count, err := r.dbContext.GetQueries(ctx).IsLabelExists(ctx, ids)
	if err != nil {
		return false, err
	}
	return count == int64(len(ids)), nil
}
