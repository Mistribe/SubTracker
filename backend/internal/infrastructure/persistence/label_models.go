package persistence

import (
	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/internal/domain/label"
	"github.com/oleexo/subtracker/internal/infrastructure/persistence/sql"
)

func createLabelFromSqlc(sqlcLabel sql.Label) label.Label {
	ownerType := auth.MustParseOwnerType(sqlcLabel.OwnerType)
	owner := auth.NewOwner(ownerType, sqlcLabel.OwnerFamilyID, sqlcLabel.OwnerUserID)
	lbl := label.NewLabel(sqlcLabel.ID,
		owner,
		sqlcLabel.Name,
		sqlcLabel.Key,
		sqlcLabel.Color,
		sqlcLabel.CreatedAt,
		sqlcLabel.UpdatedAt)
	lbl.Clean()
	return lbl
}
