package persistence

import (
	"database/sql"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/family"
)

type familySqlModel struct {
	baseSqlModel
	Name    string                 `gorm:"type:varchar(100);not null"`
	OwnerId string                 `gorm:"type:varchar(100);not null"`
	Members []familyMemberSqlModel `gorm:"foreignKey:FamilyId;references:Id"`
}

func (f familySqlModel) TableName() string {
	return "families"
}

func newFamilySqlModel(source family.Family) familySqlModel {
	return familySqlModel{
		baseSqlModel: newBaseSqlModel(source),
		Name:         source.Name(),
		OwnerId:      source.OwnerId(),
	}
}

func newFamily(source familySqlModel) family.Family {
	members := make([]family.Member, 0, len(source.Members))
	for _, member := range source.Members {
		members = append(members, newFamilyMember(member))
	}
	fam := family.NewFamily(
		source.Id,
		source.OwnerId,
		source.Name,
		members,
		source.CreatedAt,
		source.UpdatedAt,
	)

	fam.Clean()
	return fam
}

type familyMemberSqlModel struct {
	baseSqlModel
	Name     string         `gorm:"type:varchar(100);not null"`
	FamilyId uuid.UUID      `gorm:"type:uuid;not null"`
	Family   familySqlModel `gorm:"foreignKey:FamilyId;references:Id"`
	Email    sql.NullString `gorm:"type:varchar(100)"`
	UserId   sql.NullString `gorm:"type:varchar(100)"`
	IsKid    bool           `gorm:"type:boolean;not null;default:false"`
}

func (f familyMemberSqlModel) TableName() string {
	return "family_members"
}

func newFamilyMember(source familyMemberSqlModel) family.Member {
	var email *string
	if source.Email.Valid {
		email = &source.Email.String
	} else {
		email = nil
	}
	mbr := family.NewMember(
		source.Id,
		source.FamilyId,
		source.Name,
		email,
		source.IsKid,
		source.CreatedAt,
		source.UpdatedAt,
	)
	mbr.Clean()
	return mbr
}

func newFamilyMemberSqlModel(source family.Member) familyMemberSqlModel {
	model := familyMemberSqlModel{
		baseSqlModel: baseSqlModel{
			Id:        source.Id(),
			CreatedAt: source.CreatedAt(),
			UpdatedAt: source.UpdatedAt(),
			Etag:      source.ETag(),
		},
		Name:     source.Name(),
		IsKid:    source.IsKid(),
		FamilyId: source.FamilyId(),
	}

	if source.Email() != nil {
		model.Email = sql.NullString{
			String: *source.Email(),
			Valid:  true,
		}
	} else {
		model.Email = sql.NullString{
			Valid: false,
		}
	}

	return model
}
