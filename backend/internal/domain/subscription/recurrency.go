package subscription

type RecurrencyType string

const (
	UnknownRecurrency RecurrencyType = "unknown"
	OneTime           RecurrencyType = "one_time"
	Monthly           RecurrencyType = "monthly"
	Quarterly         RecurrencyType = "quarterly"
	HalfYearly        RecurrencyType = "half_yearly"
	Yearly            RecurrencyType = "yearly"
	Custom            RecurrencyType = "custom"
)

func (r RecurrencyType) String() string {
	return string(r)
}

func ParseRecurrencyType(input string) (RecurrencyType, error) {
	switch input {
	case string(OneTime):
		return OneTime, nil
	case string(Monthly):
		return Monthly, nil
	case string(Quarterly):
		return Quarterly, nil
	case string(HalfYearly):
		return HalfYearly, nil
	case string(Yearly):
		return Yearly, nil
	case string(Custom):
		return Custom, nil
	default:
		return UnknownRecurrency, ErrUnknownRecurrencyType
	}
}
