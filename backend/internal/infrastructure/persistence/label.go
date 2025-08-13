package persistence

import (
	"context"
	dsql "database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/internal/domain/label"
	"github.com/oleexo/subtracker/internal/infrastructure/persistence/sql"
	"github.com/oleexo/subtracker/pkg/slicesx"
)

type LabelRepository struct {
	dbContext *DatabaseContext
}

func NewLabelRepository(dbContext *DatabaseContext) label.Repository {
	return &LabelRepository{
		dbContext: dbContext,
	}
}

func (r LabelRepository) GetById(ctx context.Context, labelId uuid.UUID) (label.Label, error) {
	response, err := r.dbContext.GetQueries(ctx).GetLabelById(ctx, labelId)
	if err != nil {
		if errors.Is(err, dsql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	lbl := createLabelFromSqlc(response)
	lbl.Clean()
	return lbl, nil
}

func (r LabelRepository) GetByIdForUser(ctx context.Context, userId string, labelId uuid.UUID) (label.Label, error) {
	response, err := r.dbContext.GetQueries(ctx).GetLabelByIdForUser(ctx, sql.GetLabelByIdForUserParams{
		ID:          labelId,
		OwnerUserID: &userId,
	})
	if err != nil {
		if errors.Is(err, dsql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	lbl := createLabelFromSqlc(response.Label)
	lbl.Clean()
	return lbl, nil
}

func (r LabelRepository) GetAll(ctx context.Context, userId string, parameters label.QueryParameters) (
	[]label.Label,
	int64,
	error) {
	var labels []label.Label
	var totalCount int64
	if parameters.SearchText != "" {
		response, err := r.dbContext.GetQueries(ctx).GetLabelsWithSearchText(ctx, sql.GetLabelsWithSearchTextParams{
			OwnerUserID: &userId,
			Name:        fmt.Sprintf("%%%s%%", parameters.SearchText),
			Limit:       parameters.Limit,
			Offset:      parameters.Offset,
		})
		if err != nil {
			return nil, 0, err
		}
		if len(response) == 0 {
			return nil, 0, nil
		}
		totalCount = response[0].TotalCount
		labels = slicesx.Select(response, func(model sql.GetLabelsWithSearchTextRow) label.Label {
			return createLabelFromSqlc(model.Label)
		})
	} else {
		response, err := r.dbContext.GetQueries(ctx).GetLabels(ctx, sql.GetLabelsParams{
			OwnerUserID: &userId,
			Limit:       parameters.Limit,
			Offset:      parameters.Offset,
		})
		if err != nil {
			return nil, 0, err
		}
		if len(response) == 0 {
			return nil, 0, nil
		}
		totalCount = response[0].TotalCount
		labels = slicesx.Select(response, func(model sql.GetLabelsRow) label.Label {
			return createLabelFromSqlc(model.Label)
		})
	}

	return labels, totalCount, nil
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
