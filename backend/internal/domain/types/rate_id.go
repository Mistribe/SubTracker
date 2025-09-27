package types

import (
	"github.com/google/uuid"
)

type RateID uuid.UUID

func (r RateID) String() string {
	return r.String()
}

func NewRateID() RateID {
	return RateID(uuid.Must(uuid.NewV7()))
}
