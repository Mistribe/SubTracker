package user

import (
	"github.com/mistribe/subtracker/pkg/x"
)

type Feature string

const (
	FeatureActiveSubscriptions Feature = "active_subscriptions"
	FeatureCustomLabels        Feature = "custom_labels"
	FeatureCustomProviders     Feature = "custom_providers"
)

var (
	Features = map[Feature]map[Plan]FeatureDetail{
		FeatureActiveSubscriptions: {
			PlanFree:    newFeatureDetail(FeatureActiveSubscriptions, x.P(int64(10))),
			PlanPremium: newFeatureDetail(FeatureActiveSubscriptions, nil),
		},
		FeatureCustomLabels: {
			PlanPremium: newFeatureDetail(FeatureCustomLabels, nil),
		},
		FeatureCustomProviders: {
			PlanPremium: newFeatureDetail(FeatureCustomProviders, nil),
		},
	}
)

type FeatureDetail interface {
	Name() Feature
	Limit() *int64
	HasLimit() bool
}

type featureDetail struct {
	feature Feature
	limit   *int64
}

func (f featureDetail) Name() Feature {
	return f.feature
}

func (f featureDetail) Limit() *int64 {
	return f.limit
}

func (f featureDetail) HasLimit() bool {
	return f.limit != nil
}

func newFeatureDetail(feature Feature, limit *int64) FeatureDetail {
	return featureDetail{
		feature: feature,
		limit:   limit,
	}
}
