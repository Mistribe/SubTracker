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
