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

func createFamilyFromSqlcRows[T any](rows []T,
	getFamilyFunc func(T) sql.Family,
	getMemberFunc func(T) *sql.FamilyMember) []family.Family {
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
