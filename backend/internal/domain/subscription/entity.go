package subscription

import (
    "github.com/google/uuid"
    "github.com/oleexo/subtracker/internal/domain/family"
    "github.com/oleexo/subtracker/internal/domain/label"
    "time"
)

type Subscription struct {
    id              uuid.UUID
    name            string
    payments        []Payment
    labels          []label.Label
    familyMembers   []family.FamilyMember
    familyMemberIds []uuid.UUID
    payer           family.FamilyMember
    payerId         uuid.UUID
    createdAt       time.Time
    updatedAt       time.Time
}

type Payment struct {
    id        uuid.UUID
    price     float64
    startDate time.Time
    endDate   *time.Time
    months    int
    currency  string
    createdAt time.Time
    updatedAt time.Time
}
