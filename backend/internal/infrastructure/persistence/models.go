package persistence

import (
	"time"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/entity"
	"github.com/oleexo/subtracker/internal/domain/user"
)

type baseSqlModel struct {
	Id        uuid.UUID `gorm:"primaryKey;type:uuid"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Etag      string
}

func newBaseSqlModel(entity entity.Entity) baseSqlModel {
	return baseSqlModel{
		Id:        entity.Id(),
		CreatedAt: entity.CreatedAt(),
		UpdatedAt: entity.UpdatedAt(),
		Etag:      entity.ETag(),
	}
}

type ownerSqlModel struct {
	baseSqlModel

	Type     string          `gorm:"type:varchar(100);not null"`
	FamilyId *uuid.UUID      `gorm:"type:uuid"`
	Family   *familySqlModel `gorm:"foreignKey:FamilyId;references:Id"`
	UserId   *string         `gorm:"type:varchar(100)"`
}

func (o ownerSqlModel) TableName() string {
	return "owners"
}

func newOwnerSqlModel(source user.Owner) ownerSqlModel {
	model := ownerSqlModel{
		baseSqlModel: newBaseSqlModel(source),
		Type:         source.Type().String(),
	}

	switch source.Type() {
	case user.FamilyOwner:
		familyId := source.FamilyId()
		model.FamilyId = &familyId
		break
	case user.PersonalOwner:
		ownerId := source.OwnerId()
		model.UserId = &ownerId
		break
	default:
		panic("unknown owner type")
	}

	return model
}

func newOwner(model ownerSqlModel) user.Owner {
	ownerType, err := user.ParseOwnerType(model.Type)
	if err != nil {
		panic(err)
	}

	switch ownerType {
	case user.FamilyOwner:
		if model.FamilyId == nil {
			panic("family id is nil for family owner")
		}
		return user.NewFamilyOwner(model.Id,
			*model.FamilyId,
			model.CreatedAt,
			model.UpdatedAt)
	case user.PersonalOwner:
		if model.UserId == nil {
			panic("user id is nil for personal owner")
		}
		return user.NewPersonalOwner(model.Id,
			*model.UserId,
			model.CreatedAt,
			model.UpdatedAt)
	default:
		panic("unknown owner type")
	}
}
