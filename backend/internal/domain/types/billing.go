package types

import "errors"

type (
	FeatureType uint8
	PlanID      uint8
	FeatureID   uint8
)

const (
	PlanUnknown PlanID = iota
	PlanFree
	PlanPremium
)

const (
	PlanUnknownString = "unknown"
	PlanFreeString    = "free"
	PlanPremiumString = "premium"
)

func ParsePlan(input string) (PlanID, error) {
	switch input {
	case PlanUnknownString:
		return PlanUnknown, nil
	case PlanFreeString:
		return PlanFree, nil
	case PlanPremiumString:
		return PlanPremium, nil
	}

	return PlanUnknown, errors.New("invalid plan")
}

func ParsePlanOrDefault(input string, defaultValue PlanID) PlanID {
	p, err := ParsePlan(input)
	if err != nil {
		return defaultValue
	}
	return p
}

func (p PlanID) String() string {
	switch p {
	case PlanFree:
		return PlanFreeString
	case PlanPremium:
		return PlanPremiumString
	default:
		return PlanUnknownString
	}
}
