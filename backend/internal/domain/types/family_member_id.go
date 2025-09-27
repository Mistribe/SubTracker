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
