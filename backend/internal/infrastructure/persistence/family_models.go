package persistence

import (
	"database/sql"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/family"
)

type FamilySqlModel struct {
	BaseSqlModel `gorm:"embedded"`

	Name    string                 `gorm:"type:varchar(100);not null"`
	OwnerId string                 `gorm:"type:varchar(100);not null"`
	Members []FamilyMemberSqlModel `gorm:"foreignKey:FamilyId;references:Id"`
}

func (f FamilySqlModel) TableName() string {
	return "families"
}

type FamilyMemberSqlModel struct {
	BaseSqlModel `gorm:"embedded"`

	Name     string         `gorm:"type:varchar(100);not null"`
	FamilyId uuid.UUID      `gorm:"type:uuid;not null"`
	UserId   sql.NullString `gorm:"type:varchar(100)"`
	IsKid    bool           `gorm:"type:boolean;not null;default:false"`
}

func (f FamilyMemberSqlModel) TableName() string {
	return "family_members"
}
func newFamilySqlModel(source family.Family) FamilySqlModel {
	return FamilySqlModel{
		BaseSqlModel: newBaseSqlModel(source),
		Name:         source.Name(),
		OwnerId:      source.OwnerId(),
	}
}

func newFamily(source FamilySqlModel) family.Family {
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

func newFamilyMember(source FamilyMemberSqlModel) family.Member {
	mbr := family.NewMember(
		source.Id,
		source.FamilyId,
		source.Name,
		source.IsKid,
		source.CreatedAt,
		source.UpdatedAt,
	)
	if source.UserId.Valid {
		mbr.SetUserId(&source.UserId.String)
	}
	mbr.Clean()
	return mbr
}

func newFamilyMemberSqlModel(source family.Member) FamilyMemberSqlModel {
	model := FamilyMemberSqlModel{
		BaseSqlModel: BaseSqlModel{
			Id:        source.Id(),
			CreatedAt: source.CreatedAt(),
			UpdatedAt: source.UpdatedAt(),
			Etag:      source.ETag(),
		},
		Name:     source.Name(),
		IsKid:    source.IsKid(),
		FamilyId: source.FamilyId(),
	}

	if source.UserId() != nil {
		model.UserId = sql.NullString{
			String: *source.UserId(),
			Valid:  true,
		}
	} else {
		model.UserId = sql.NullString{
			Valid: false,
		}
	}

	return model
}
