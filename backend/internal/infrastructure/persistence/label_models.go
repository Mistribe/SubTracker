package persistence

import (
	"github.com/mistribe/subtracker/internal/domain/auth"
	"github.com/mistribe/subtracker/internal/domain/label"
	"github.com/mistribe/subtracker/internal/infrastructure/persistence/jet/app/public/model"
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
