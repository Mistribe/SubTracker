package subscription

import (
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/pkg/x/herd"
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
	LabelId types.LabelID
	Source  LabelSource
}

func LabelRefUniqueComparer(a LabelRef, b LabelRef) bool {
	return a.LabelId == b.LabelId && a.Source == b.Source
}

func LabelRefComparer(a LabelRef, b LabelRef) bool {
	return a.LabelId == b.LabelId && a.Source == b.Source
}

func NewSubscriptionLabelRefs(labelIDs []types.LabelID) []LabelRef {
	return herd.Select(labelIDs, func(ID types.LabelID) LabelRef {
		return LabelRef{
			LabelId: ID,
			Source:  LabelSourceSubscription,
		}
	})
}
