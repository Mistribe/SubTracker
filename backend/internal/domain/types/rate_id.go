package types

import (
	"github.com/google/uuid"
)

type RateID uuid.UUID

func (r RateID) String() string {
	u := uuid.UUID(r)
	return u.String()
}

func NewRateID() RateID {
	return RateID(uuid.Must(uuid.NewV7()))
}
