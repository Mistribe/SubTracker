package types

import (
	"github.com/google/uuid"
)

type LabelID uuid.UUID

func (l LabelID) String() string {
	u := uuid.UUID(l)
	return u.String()
}

func NewLabelID() LabelID {
	return LabelID(uuid.Must(uuid.NewV7()))
}

func LabelIDComparer(l1, l2 LabelID) bool {
	return l1 == l2
}

func ParseLabelID(s string) (LabelID, error) {
	u, err := uuid.Parse(s)
	if err != nil {
		return LabelID{}, err
	}
	return LabelID(u), nil
}

func ParseLabelIDOrNil(s *string) (*LabelID, error) {
	if s == nil {
		return nil, nil
	}

	u, err := ParseLabelID(*s)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
