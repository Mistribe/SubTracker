package types

import (
	"github.com/google/uuid"
)

type LabelID uuid.UUID

func (l LabelID) String() string {
	return l.String()
}

func NewLabelID() LabelID {
	return LabelID(uuid.Must(uuid.NewV7()))
}

func LabelIDComparer(l1, l2 LabelID) bool {
	return l1 == l2
}
