package types

import (
	"github.com/google/uuid"
)

type FamilyMemberID uuid.UUID

func (f FamilyMemberID) String() string {
	return f.String()
}

func NewFamilyMemberID() FamilyMemberID {
	return FamilyMemberID(uuid.Must(uuid.NewV7()))
}

func FamilyMemberIdComparer(fm1, fm2 FamilyMemberID) bool {
	return fm1 == fm2
}

func ParseFamilyMemberID(s string) (FamilyMemberID, error) {
	u, err := uuid.Parse(s)
	if err != nil {
		return FamilyMemberID{}, err
	}
	return FamilyMemberID(u), nil
}

func ParseFamilyMemberIDOrNil(s *string) (*FamilyMemberID, error) {
	if s == nil {
		return nil, nil
	}

	u, err := ParseFamilyMemberID(*s)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func MustParseFamilyMemberID(s string) FamilyMemberID {
	u, err := ParseFamilyMemberID(s)
	if err != nil {
		panic(err)
	}
	return u
}
