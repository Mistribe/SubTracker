package types

import (
	"github.com/google/uuid"
)

type SubscriptionID uuid.UUID

func (id SubscriptionID) String() string {
	return id.String()
}

func NewSubscriptionID() SubscriptionID {
	return SubscriptionID(uuid.Must(uuid.NewV7()))
}
