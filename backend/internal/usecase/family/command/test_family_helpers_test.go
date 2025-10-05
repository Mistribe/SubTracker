package command_test

import (
	"time"

	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/domain/types"
)

func helperNewMember(famID types.FamilyID, name string, mt family.MemberType) family.Member {
	return family.NewMember(types.NewFamilyMemberID(), famID, name, mt, nil, time.Now(), time.Now())
}

func helperNewMemberWithID(
	famID types.FamilyID,
	memberID types.FamilyMemberID,
	name string,
	mt family.MemberType) family.Member {
	return family.NewMember(memberID, famID, name, mt, nil, time.Now(), time.Now())
}

func helperNewFamily(owner types.UserID, name string, members []family.Member) family.Family {
	return family.NewFamily(types.NewFamilyID(), owner, name, members, time.Now(), time.Now())
}

func helperNewFamilyWithID(
	famID types.FamilyID,
	owner types.UserID,
	name string,
	members []family.Member) family.Family {
	return family.NewFamily(famID, owner, name, members, time.Now(), time.Now())
}
