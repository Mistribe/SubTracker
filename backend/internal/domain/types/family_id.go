package types

import (
	"github.com/google/uuid"
)

type FamilyID uuid.UUID

func (f FamilyID) String() string {
	u := uuid.UUID(f)
	return u.String()
}

func NewFamilyID() FamilyID {
	return FamilyID(uuid.Must(uuid.NewV7()))
}

func ParseFamilyID(s string) (FamilyID, error) {
	u, err := uuid.Parse(s)
	if err != nil {
		return FamilyID{}, err
	}
	return FamilyID(u), nil
}

func ParseFamilyIDOrNil(s *string) (*FamilyID, error) {
	if s == nil {
		return nil, nil
	}

	u, err := ParseFamilyID(*s)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
