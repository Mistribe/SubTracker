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

func (r RecurrencyType) TotalMonths() int {
	switch r {
	case MonthlyRecurrency:
		return 1
	case QuarterlyRecurrency:
		return 3
	case HalfYearlyRecurrency:
		return 6
	case YearlyRecurrency:
		return 12
	default:
		return 0
	}
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
	customRecurrency *int32,
	to RecurrencyType) float64 {
	if from == to {
		return amount
	}

	monthlyAmount := amount
	switch from {
	case MonthlyRecurrency:
		monthlyAmount = amount
	case QuarterlyRecurrency:
		monthlyAmount = amount / 3
	case HalfYearlyRecurrency:
		monthlyAmount = amount / 6
	case YearlyRecurrency:
		monthlyAmount = amount / 12
	case CustomRecurrency:
		if customRecurrency != nil {
			monthlyAmount = amount / float64(*customRecurrency)
		}
		return 0
	}

	switch to {
	case MonthlyRecurrency:
		return monthlyAmount
	case QuarterlyRecurrency:
		return monthlyAmount * 3
	case HalfYearlyRecurrency:
		return monthlyAmount * 6
	case YearlyRecurrency:
		return monthlyAmount * 12
	default:
		return amount
	}
}
