package models

import (
	"github.com/mistribe/subtracker/internal/adapters/persistence/db/jet/app/public/model"
	"github.com/mistribe/subtracker/pkg/x/collection"

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

func CreateFamilyWithMembersFromModel(source FamilyWithMembers) (family.Family, error) {
	members, err := collection.SelectErr(source.FamilyMembers, func(in model.FamilyMembers) (family.Member, error) {
		return createFamilyMemberFromJet(in)
	})
	if err != nil {
		return nil, err
	}
	newFamily := family.NewFamily(
		source.ID,
		source.OwnerID,
		source.Name,
		members,
		source.CreatedAt,
		source.UpdatedAt,
	)
	newFamily.Clean()

	return newFamily, nil
}
