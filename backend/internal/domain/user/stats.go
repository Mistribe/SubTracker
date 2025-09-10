package user

type Stats interface {
	Providers() int64
	Labels() int64
	ActiveSubscriptions() int64
	GetCountFromFeature(feature Feature) int64
}

func NewStats(providers, labels, activeSubscriptions int64) Stats {
	return &stats{
		providers:           providers,
		labels:              labels,
		activeSubscriptions: activeSubscriptions,
	}
}

func NewEmptyStats() Stats {
	return &stats{}
}

type stats struct {
	providers, labels, activeSubscriptions int64
}

func (s stats) GetCountFromFeature(feature Feature) int64 {
	switch feature {
	case FeatureCustomLabels:
		return s.labels
	case FeatureCustomProviders:
		return s.providers
	case FeatureActiveSubscriptions:
		return s.activeSubscriptions
	}

	panic("unknown feature: " + feature)
}

func (s stats) Providers() int64 {
	return s.providers
}

func (s stats) Labels() int64 {
	return s.labels
}

func (s stats) ActiveSubscriptions() int64 {
	return s.activeSubscriptions
}
