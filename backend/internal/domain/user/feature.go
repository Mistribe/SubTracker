package user

import (
	"github.com/mistribe/subtracker/pkg/x"
)

type Category string

const (
	CategoryFamily       Category = "family"
	CategoryUser         Category = "user"
	CategorySubscription Category = "subscription"
	CategoryProvider     Category = "provider"
	CategoryLabel        Category = "label"
)

func (c Category) String() string {
	return string(c)
}

type Feature string

func (f Feature) String() string {
	return string(f)
}

const (
	FeatureActiveSubscriptions Feature = "active_subscriptions"
	FeatureCustomLabels        Feature = "custom_labels"
	FeatureCustomProviders     Feature = "custom_providers"
	FeatureFamilyMembers       Feature = "family_members"
)

var (
	Features = map[Feature]map[Plan]FeatureDetail{
		FeatureActiveSubscriptions: {
			PlanFree:    newFeatureDetail(CategorySubscription, FeatureActiveSubscriptions, x.Int64P(10)),
			PlanPremium: newFeatureDetail(CategorySubscription, FeatureActiveSubscriptions, nil),
		},
		FeatureCustomLabels: {
			PlanPremium: newFeatureDetail(CategoryLabel, FeatureCustomLabels, nil),
		},
		FeatureCustomProviders: {
			PlanPremium: newFeatureDetail(CategoryProvider, FeatureCustomProviders, nil),
		},
		FeatureFamilyMembers: {
			PlanFree:    newFeatureDetail(CategoryFamily, FeatureFamilyMembers, x.Int64P(25)),
			PlanPremium: newFeatureDetail(CategoryFamily, FeatureFamilyMembers, x.Int64P(25)),
		},
	}
)

type FeatureDetail interface {
	Category() Category
	Name() Feature
	Limit() *int64
	HasLimit() bool
}

type featureDetail struct {
	category Category
	feature  Feature
	limit    *int64
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

func (f featureDetail) Category() Category {
	return f.category
}

func newFeatureDetail(category Category, feature Feature, limit *int64) FeatureDetail {
	return featureDetail{
		category: category,
		feature:  feature,
		limit:    limit,
	}
}
