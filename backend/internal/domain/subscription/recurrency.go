package subscription

type RecurrencyType string

const (
	UnknownRecurrency    RecurrencyType = "unknown"
	OneTimeRecurrency    RecurrencyType = "one_time"
	MonthlyRecurrency    RecurrencyType = "monthly"
	QuarterlyRecurrency  RecurrencyType = "quarterly"
	HalfYearlyRecurrency RecurrencyType = "half_yearly"
	YearlyRecurrency     RecurrencyType = "yearly"
	CustomRecurrency     RecurrencyType = "custom"
)

func (r RecurrencyType) String() string {
	return string(r)
}

func ParseRecurrencyType(input string) (RecurrencyType, error) {
	switch input {
	case string(OneTimeRecurrency):
		return OneTimeRecurrency, nil
	case string(MonthlyRecurrency):
		return MonthlyRecurrency, nil
	case string(QuarterlyRecurrency):
		return QuarterlyRecurrency, nil
	case string(HalfYearlyRecurrency):
		return HalfYearlyRecurrency, nil
	case string(YearlyRecurrency):
		return YearlyRecurrency, nil
	case string(CustomRecurrency):
		return CustomRecurrency, nil
	default:
		return UnknownRecurrency, ErrUnknownRecurrencyType
	}
}

func MustParseRecurrencyType(input string) RecurrencyType {
	t, err := ParseRecurrencyType(input)
	if err != nil {
		panic(err)
	}
	return t
}

func toRecurrencyPrice(
	amount float64,
	from RecurrencyType,
	to RecurrencyType) float64 {
	if from == to {
		return amount
	}
	return amount * 12
}
