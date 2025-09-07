package user

type Plan string

const (
	PlanFree    Plan = "FREE"
	PlanPremium Plan = "PREMIUM"
)

type PlanLimits interface {
	MaxActiveSubscriptions() *int
	MaxCustomLabels() *int
	MaxCustomProviders() *int
}

type planLimits struct {
	maxActiveSubscriptions *int
	maxCustomLabels        *int
	maxCustomProviders     *int
}

func (p planLimits) MaxActiveSubscriptions() *int {
	return p.maxActiveSubscriptions
}

func (p planLimits) MaxCustomLabels() *int {
	return p.maxCustomLabels
}

func (p planLimits) MaxCustomProviders() *int {
	return p.maxCustomProviders
}

func (p Plan) String() string {
	return string(p)
}

func (p Plan) Limits() PlanLimits {
	switch p {
	case PlanPremium:
		return planLimits{
			maxActiveSubscriptions: nil,
			maxCustomLabels:        nil,
			maxCustomProviders:     nil,
		}
	default:
		sub := 10
		lbl := 5
		prov := 5
		return planLimits{
			maxActiveSubscriptions: &sub,
			maxCustomLabels:        &lbl,
			maxCustomProviders:     &prov,
		}
	}
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
