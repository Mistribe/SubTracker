package persistence

import (
	"github.com/oleexo/subtracker/internal/domain/label"
	"github.com/oleexo/subtracker/internal/domain/user"
)

type labelSqlModel struct {
	baseSqlModel
	baseOwnerSqlModel

	Name      string `gorm:"type:varchar(100);not null"`
	Color     string `gorm:"type:varchar(20);not null"`
	IsDefault bool   `gorm:"type:boolean;not null;default:false;index"`
}

func (l labelSqlModel) TableName() string {
	return "labels"
}

func newLabelSqlModel(source label.Label) labelSqlModel {
	model := labelSqlModel{
		baseSqlModel: baseSqlModel{
			Id:        source.Id(),
			CreatedAt: source.CreatedAt(),
			UpdatedAt: source.UpdatedAt(),
			Etag:      source.ETag(),
		},
		Name:      source.Name(),
		Color:     source.Color(),
		IsDefault: source.IsDefault(),
	}

	model.OwnerType = source.Owner().Type().String()
	switch source.Owner().Type() {
	case user.FamilyOwner:
		familyId := source.Owner().FamilyId()
		model.OwnerFamilyId = &familyId
		break
	case user.PersonalOwner:
		userId := source.Owner().UserId()
		model.OwnerUserId = &userId
		break
	}

	return model
}

func newLabel(source labelSqlModel) label.Label {
	ownerType, err := user.ParseOwnerType(source.OwnerType)
	if err != nil {
		panic(err)
	}
	owner := user.NewOwner(ownerType, source.OwnerFamilyId, source.OwnerUserId)
	lbl := label.NewLabel(
		source.Id,
		owner,
		source.Name,
		source.IsDefault,
		source.Color,
		source.CreatedAt,
		source.UpdatedAt,
	)
	lbl.Clean()
	return lbl
}
