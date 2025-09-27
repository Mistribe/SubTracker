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

func ParseProviderID(id string) (ProviderID, error) {
	u, err := uuid.Parse(id)
	if err != nil {
		return ProviderID{}, err
	}
	return ProviderID(u), nil
}

func ParseProviderIDOrNil(id *string) (*ProviderID, error) {
	if id == nil {
		return nil, nil
	}

	u, err := ParseProviderID(*id)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
