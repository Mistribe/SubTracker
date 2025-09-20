package authorization

import (
	"github.com/mistribe/subtracker/internal/domain/user"
	"github.com/mistribe/subtracker/internal/shared"
)

type limitInfo struct {
	category     user.Category
	currentCount int64
	feature      user.FeatureDetail
}

func (l limitInfo) Category() string {
	return l.category.String()
}

func (l limitInfo) Feature() string {
	return l.feature.Name().String()
}

func (l limitInfo) Remaining() int64 {
	if l.feature.Limit() == nil {
		return 0
	}

	return l.currentCount - *l.feature.Limit()

}

func (l limitInfo) Limit() int64 {
	if l.feature.Limit() == nil {
		return 0
	}
	return *l.feature.Limit()
}

func (l limitInfo) HasLimit() bool {
	return l.feature.HasLimit()
}

func (l limitInfo) CurrentCount() int64 {
	return l.currentCount
}

func newLimitInfo(currentCount int64, info user.FeatureDetail) shared.Limit {
	return &limitInfo{
		feature:      info,
		currentCount: currentCount,
	}
}
