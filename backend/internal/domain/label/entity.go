package label

import (
	"time"

	"github.com/google/uuid"
)

type Label struct {
	id        uuid.UUID
	name      string
	isDefault bool
	color     string
	createdAt time.Time
	updatedAt time.Time
}
