package types

import (
	"github.com/google/uuid"
)

type FamilyMemberID uuid.UUID

func (f FamilyMemberID) String() string {
	return f.String()
}

func FamilyMemberIdComparer(fm1, fm2 FamilyMemberID) bool {
	return fm1 == fm2
}
