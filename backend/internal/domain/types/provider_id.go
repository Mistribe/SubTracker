package types

import (
	"github.com/google/uuid"
)

type ProviderID uuid.UUID

func (id ProviderID) String() string {
	return id.String()
}

func NewProviderID() ProviderID {
	return ProviderID(uuid.Must(uuid.NewV7()))
}
