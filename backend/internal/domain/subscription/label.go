package subscription

import (
	"github.com/google/uuid"
)

type LabelSource string

const (
	LabelSourceProvider     LabelSource = "provider"
	LabelSourceSubscription LabelSource = "subscription"
)

func (l LabelSource) String() string {
	return string(l)
}

type LabelRef struct {
	LabelId uuid.UUID
	Source  LabelSource
}

func LabelRefUniqueComparer(a LabelRef, b LabelRef) bool {
	return a.LabelId == b.LabelId && a.Source == b.Source
}

func LabelRefComparer(a LabelRef, b LabelRef) bool {
	return a.LabelId == b.LabelId && a.Source == b.Source
}
