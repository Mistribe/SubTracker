package family

import (
	"github.com/google/uuid"
	"time"
)

type FamilyMember struct {
	id        uuid.UUID
	name      string
	isKid     bool
	createdAt time.Time
	updatedAt time.Time
}
