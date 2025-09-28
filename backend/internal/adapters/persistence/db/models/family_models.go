package models

import (
	"github.com/mistribe/subtracker/internal/adapters/persistence/db/jet/app/public/model"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/pkg/x"
	"github.com/mistribe/subtracker/pkg/x/herd"

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
		types.FamilyMemberID(jetMember.ID),
		types.FamilyID(jetMember.FamilyID),
		jetMember.Name,
		memberType,
		jetMember.InvitationCode,
		jetMember.CreatedAt,
		jetMember.UpdatedAt,
	)

	if jetMember.UserID != nil {
		member.SetUserId(x.P(types.UserID(*jetMember.UserID)))
	}

	member.Clean()
	return member, nil
}

func CreateFamilyWithMembersFromModel(source FamilyWithMembers) (family.Family, error) {
	members, err := herd.SelectErr(source.FamilyMembers, func(in model.FamilyMembers) (family.Member, error) {
		return createFamilyMemberFromJet(in)
	})
	if err != nil {
		return nil, err
	}
	newFamily := family.NewFamily(
		types.FamilyID(source.ID),
		types.UserID(source.OwnerID),
		source.Name,
		members,
		source.CreatedAt,
		source.UpdatedAt,
	)
	newFamily.Clean()

	return newFamily, nil
}
