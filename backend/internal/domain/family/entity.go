package family

import (
	"time"

	"github.com/google/uuid"
)

type FamilyMember struct {
	id        uuid.UUID
	name      string
	isKid     bool
	createdAt time.Time
	updatedAt time.Time
}
