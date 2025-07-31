package persistence

import (
	"database/sql"

	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/internal/domain/label"
)

type LabelSqlModel struct {
	BaseSqlModel      `gorm:"embedded"`
	BaseOwnerSqlModel `gorm:"embedded"`

	Name  string `gorm:"type:varchar(100);not null"`
	Key   string `gorm:"type:varchar(100);not null"`
	Color string `gorm:"type:varchar(50);not null"`
}

func (l LabelSqlModel) TableName() string {
	return "labels"
}

func newLabelSqlModel(source label.Label) LabelSqlModel {
	model := LabelSqlModel{
		BaseSqlModel: BaseSqlModel{
			Id:        source.Id(),
			CreatedAt: source.CreatedAt(),
			UpdatedAt: source.UpdatedAt(),
			Etag:      source.ETag(),
		},
		Name:  source.Name(),
		Key:   source.Key(),
		Color: source.Color(),
	}

	model.OwnerType = source.Owner().Type().String()
	switch source.Owner().Type() {
	case auth.FamilyOwner:
		familyId := source.Owner().FamilyId()
		model.OwnerFamilyId = &familyId
	case auth.PersonalOwner:
		userId := source.Owner().UserId()
		model.OwnerUserId = sql.NullString{
			String: userId,
			Valid:  true,
		}
	}

	return model
}

func newLabel(source LabelSqlModel) label.Label {
	ownerType, err := auth.ParseOwnerType(source.OwnerType)
	if err != nil {
		panic(err)
	}
	var ownerFamilyId *string
	if source.OwnerUserId.Valid {
		ownerFamilyId = &source.OwnerUserId.String
	}
	owner := auth.NewOwner(ownerType, source.OwnerFamilyId, ownerFamilyId)
	lbl := label.NewLabel(
		source.Id,
		owner,
		source.Name,
		source.Key,
		source.Color,
		source.CreatedAt,
		source.UpdatedAt,
	)
	lbl.Clean()
	return lbl
}
