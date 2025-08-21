package persistence

import (
	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/internal/domain/label"
	"github.com/oleexo/subtracker/internal/infrastructure/persistence/jet/app/public/model"
)

func createLabelFromJet(jetLabel model.Labels) label.Label {
	ownerType := auth.MustParseOwnerType(jetLabel.OwnerType)
	owner := auth.NewOwner(ownerType, jetLabel.OwnerFamilyID, jetLabel.OwnerUserID)
	lbl := label.NewLabel(jetLabel.ID,
		owner,
		jetLabel.Name,
		jetLabel.Key,
		jetLabel.Color,
		jetLabel.CreatedAt,
		jetLabel.UpdatedAt)
	lbl.Clean()
	return lbl
}
