package persistence

import (
	"database/sql"
	"time"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/entity"
)

type BaseSqlModel struct {
	Id        uuid.UUID `gorm:"primaryKey;type:uuid;not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
	Etag      string    `gorm:"not null;uniqueIndex:idx_etag"`
}

type BaseOwnerSqlModel struct {
	OwnerType     string          `gorm:"type:varchar(20);not null"`
	OwnerFamilyId *uuid.UUID      `gorm:"type:uuid"`
	OwnerFamily   *FamilySqlModel `gorm:"foreignKey:OwnerFamilyId;references:Id"`
	OwnerUserId   sql.NullString  `gorm:"type:varchar(50)"`
}

func newBaseSqlModel(entity entity.Entity) BaseSqlModel {
	return BaseSqlModel{
		Id:        entity.Id(),
		CreatedAt: entity.CreatedAt(),
		UpdatedAt: entity.UpdatedAt(),
		Etag:      entity.ETag(),
	}
}
