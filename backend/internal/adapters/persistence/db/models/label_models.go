package models

import (
	"github.com/mistribe/subtracker/internal/adapters/persistence/db/jet/app/public/model"
	"github.com/mistribe/subtracker/internal/domain/label"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/pkg/x"
)

func CreateLabelFromModel(source model.Labels) label.Label {
	ownerType := types.MustParseOwnerType(source.OwnerType)
	var ownerFamilyID *types.FamilyID
	var ownerUserID *types.UserID
	if source.OwnerFamilyID != nil {
		ownerFamilyID = x.P(types.FamilyID(*source.OwnerFamilyID))
	}
	if source.OwnerUserID != nil {
		ownerUserID = x.P(types.UserID(*source.OwnerUserID))
	}
	owner := types.NewOwner(ownerType, ownerFamilyID, ownerUserID)
	lbl := label.NewLabel(
		types.LabelID(source.ID),
		owner,
		source.Name,
		source.Key,
		source.Color,
		source.CreatedAt,
		source.UpdatedAt)
	lbl.Clean()
	return lbl
}
