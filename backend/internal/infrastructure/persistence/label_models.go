package persistence

import (
	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/label"
	"github.com/oleexo/subtracker/internal/domain/user"
)

type labelSqlModel struct {
	baseSqlModel
	OwnerId   *uuid.UUID     `gorm:"type:uuid"`
	Owner     *ownerSqlModel `gorm:"foreignKey:OwnerId;references:Id"`
	Name      string         `gorm:"type:varchar(100);not null"`
	Color     string         `gorm:"type:varchar(20);not null"`
	IsDefault bool           `gorm:"type:boolean;not null;default:false;index"`
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

	if source.Owner() != nil {
		ownerId := source.Owner().Id()
		model.OwnerId = &ownerId
	}

	return model
}

func newLabel(source labelSqlModel) label.Label {
	var owner user.Owner
	if source.Owner != nil {
		owner = newOwner(*source.Owner)
	}
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
