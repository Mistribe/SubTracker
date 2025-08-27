package x

import (
	"github.com/google/uuid"
)

func ParseOrNewUUID(input *string) (uuid.UUID, error) {
	return TernaryFunc2(input != nil, func() (uuid.UUID, error) {
		return uuid.Parse(*input)
	}, func() (uuid.UUID, error) {
		return uuid.NewV7()
	})
}

func ParseOrNilUUID(input *string) (*uuid.UUID, error) {
	return TernaryFunc2(input != nil, func() (*uuid.UUID, error) {
		id, err := uuid.Parse(*input)
		if err != nil {
			return nil, err
		}
		return &id, nil
	}, func() (*uuid.UUID, error) {
		return nil, nil
	})
}
