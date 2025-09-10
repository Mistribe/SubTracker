package user

type Plan string

const (
	PlanFree    Plan = "FREE"
	PlanPremium Plan = "PREMIUM"
)

func (p Plan) String() string {
	return string(p)
}

func (p Plan) HasFeature(feature Feature) bool {
	switch feature {
	case FeatureCustomLabels:
		return p == PlanPremium
	case FeatureCustomProviders:
		return p == PlanPremium
	case FeatureActiveSubscriptions:
		return true
	default:
		return false
	}
}

func (p Plan) GetFeatureDetail(feature Feature) FeatureDetail {
	var list map[Plan]FeatureDetail
	var exists bool
	var detail FeatureDetail

	if list, exists = Features[feature]; !exists {
		return nil
	}
	if detail, exists = list[p]; !exists {
		return nil
	}

	return detail
}

func ParsePlan(in string) (Plan, error) {
	switch in {
	case string(PlanFree):
		return PlanFree, nil
	case string(PlanPremium):
		return PlanPremium, nil
	default:
		return "", ErrUnknownPlan
	}
}

func MustParsePlan(in string) Plan {
	p, err := ParsePlan(in)
	if err != nil {
		panic(err)
	}
	return p
}
