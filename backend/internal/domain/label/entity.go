package label

import (
    "github.com/google/uuid"
    "time"
)

type Label struct {
    id        uuid.UUID
    name      string
    isDefault bool
    color     string
    createdAt time.Time
    updatedAt time.Time
}
