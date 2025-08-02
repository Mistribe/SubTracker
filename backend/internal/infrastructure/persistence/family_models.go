package persistence

import (
	"github.com/oleexo/subtracker/internal/infrastructure/persistence/sql"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/family"
)

func createFamilyMemberFromSqlc(sqlcMember sql.FamilyMember) family.Member {
	memberType := family.MustParseMemberType(sqlcMember.Type)
	member := family.NewMember(
		sqlcMember.ID,
		sqlcMember.FamilyID,
		sqlcMember.Name,
		memberType,
		sqlcMember.CreatedAt,
		sqlcMember.UpdatedAt,
	)

	if sqlcMember.UserID != nil {
		member.SetUserId(sqlcMember.UserID)
	}

	member.Clean()
	return member
}

func createFamilyFromSqlcRows[T any](rows []T, getFamilyFunc func(T) sql.Family,
	getMemberFunc func(T) sql.FamilyMember) []family.Family {
	if len(rows) == 0 {
		return nil
	}

	familyMap := make(map[uuid.UUID][]family.Member)
	var lastFamily sql.Family

	for _, row := range rows {
		sqlcFamily := getFamilyFunc(row)
		sqlcMember := getMemberFunc(row)

		lastFamily = sqlcFamily // Keep reference for family data

		// Add member if valid (not null from LEFT JOIN)
		if sqlcMember.ID != uuid.Nil {
			member := createFamilyMemberFromSqlc(sqlcMember)
			familyMap[sqlcFamily.ID] = append(familyMap[sqlcFamily.ID], member)
		}
	}

	// Convert to domain families
	families := make([]family.Family, 0, len(familyMap))
	for _, members := range familyMap {
		domainFamily := family.NewFamily(
			lastFamily.ID,
			lastFamily.OwnerID,
			lastFamily.Name,
			members,
			lastFamily.CreatedAt,
			lastFamily.UpdatedAt,
		)
		domainFamily.Clean()
		families = append(families, domainFamily)
	}

	return families
}

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
	Type     string         `gorm:"type:varchar(10);not null"`
}

func (f FamilyMemberSqlModel) TableName() string {
	return "family_members"
}
func newFamilySqlModel(source family.Family) FamilySqlModel {
	return FamilySqlModel{
		BaseSqlModel: newBaseSqlModel(source, source.ETag()),
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
	memberType := family.MustParseMemberType(source.Type)
	mbr := family.NewMember(
		source.Id,
		source.FamilyId,
		source.Name,
		memberType,
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
		Type:     source.Type().String(),
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
