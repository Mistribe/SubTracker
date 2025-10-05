package types

import (
	"github.com/google/uuid"
)

type SubscriptionID uuid.UUID

func (id SubscriptionID) String() string {
	// Convert to real uuid.UUID to avoid infinite recursion
	u := uuid.UUID(id)
	return u.String()
}

func NewSubscriptionID() SubscriptionID {
	return SubscriptionID(uuid.Must(uuid.NewV7()))
}

func ParseSubscriptionID(id string) (SubscriptionID, error) {
	u, err := uuid.Parse(id)
	if err != nil {
		return SubscriptionID{}, err
	}
	return SubscriptionID(u), nil
}

func ParseSubscriptionIDOrNil(id *string) (*SubscriptionID, error) {
	if id == nil {
		return nil, nil
	}

	u, err := ParseSubscriptionID(*id)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func MustParseSubscriptionID(id string) SubscriptionID {
	u, err := ParseSubscriptionID(id)
	if err != nil {
		panic(err)
	}

	return u
}
