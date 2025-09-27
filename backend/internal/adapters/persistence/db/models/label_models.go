package models

import (
	"github.com/mistribe/subtracker/internal/adapters/persistence/db/jet/app/public/model"
	"github.com/mistribe/subtracker/internal/domain/label"
	"github.com/mistribe/subtracker/internal/domain/types"
)

func CreateLabelFromModel(source model.Labels) label.Label {
	ownerType := types.MustParseOwnerType(source.OwnerType)
	owner := types.NewOwner(ownerType, source.OwnerFamilyID, source.OwnerUserID)
	lbl := label.NewLabel(source.ID,
		owner,
		source.Name,
		source.Key,
		source.Color,
		source.CreatedAt,
		source.UpdatedAt)
	lbl.Clean()
	return lbl
}
