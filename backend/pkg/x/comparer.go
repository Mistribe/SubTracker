package x

import (
	"github.com/google/uuid"
)

func UuidUniqueComparer(a uuid.UUID, b uuid.UUID) bool {
	return a == b
}

func UuidComparer(a uuid.UUID, b uuid.UUID) bool {
	return a == b
}
