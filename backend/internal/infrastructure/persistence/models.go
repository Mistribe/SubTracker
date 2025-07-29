package persistence

import (
	"time"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/entity"
)

type baseSqlModel struct {
	Id        uuid.UUID `gorm:"primaryKey;type:uuid"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Etag      string
}

type baseOwnerSqlModel struct {
	OwnerType     string          `gorm:"type:varchar(20);not null"`
	OwnerFamilyId *uuid.UUID      `gorm:"type:uuid"`
	OwnerFamily   *familySqlModel `gorm:"foreignKey:OwnerFamilyId;references:Id"`
	OwnerUserId   *string         `gorm:"type:varchar(20);not null"`
}

func newBaseSqlModel(entity entity.Entity) baseSqlModel {
	return baseSqlModel{
		Id:        entity.Id(),
		CreatedAt: entity.CreatedAt(),
		UpdatedAt: entity.UpdatedAt(),
		Etag:      entity.ETag(),
	}
}
