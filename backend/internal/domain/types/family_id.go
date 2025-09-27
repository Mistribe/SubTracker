package types

import (
	"github.com/google/uuid"
)

type FamilyID uuid.UUID

func (f FamilyID) String() string {
	return f.String()
}
