package persistence

import (
	"github.com/mistribe/subtracker/internal/infrastructure/persistence/jet/app/public/model"
	"github.com/mistribe/subtracker/pkg/slicesx"

	"github.com/mistribe/subtracker/internal/domain/family"
)

type FamilyWithMembers struct {
	model.Families

	FamilyMembers []model.FamilyMembers
}

func createFamilyMemberFromJet(jetMember model.FamilyMembers) (family.Member, error) {
	memberType, err := family.ParseMemberType(jetMember.Type)
	if err != nil {
		return nil, err
	}
	member := family.NewMember(
		jetMember.ID,
		jetMember.FamilyID,
		jetMember.Name,
		memberType,
		jetMember.InvitationCode,
		jetMember.CreatedAt,
		jetMember.UpdatedAt,
	)

	if jetMember.UserID != nil {
		member.SetUserId(jetMember.UserID)
	}

	member.Clean()
	return member, nil
}

func createFamilyWithMembersFromJet(jetFamily FamilyWithMembers) (family.Family, error) {
	members, err := slicesx.SelectErr(jetFamily.FamilyMembers, func(in model.FamilyMembers) (family.Member, error) {
		return createFamilyMemberFromJet(in)
	})
	if err != nil {
		return nil, err
	}
	newFamily := family.NewFamily(
		jetFamily.ID,
		jetFamily.OwnerID,
		jetFamily.Name,
		members,
		jetFamily.CreatedAt,
		jetFamily.UpdatedAt,
	)
	newFamily.Clean()

	return newFamily, nil
}
