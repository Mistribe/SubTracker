package subscription

type RecurrencyType string
type RecurrencyTypes []RecurrencyType

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
	result := ParseRecurrencyTypeOrDefault(input, UnknownRecurrency)
	if result == UnknownRecurrency {
		return UnknownRecurrency, ErrUnknownRecurrencyType
	}
	return result, nil
}

func ParseRecurrencyTypeOrDefault(input string, defaultValue RecurrencyType) RecurrencyType {
	switch input {
	case string(OneTimeRecurrency):
		return OneTimeRecurrency
	case string(MonthlyRecurrency):
		return MonthlyRecurrency
	case string(QuarterlyRecurrency):
		return QuarterlyRecurrency
	case string(HalfYearlyRecurrency):
		return HalfYearlyRecurrency
	case string(YearlyRecurrency):
		return YearlyRecurrency
	case string(CustomRecurrency):
		return CustomRecurrency
	default:
		return defaultValue
	}
}

func MustParseRecurrencyType(input string) RecurrencyType {
	t, err := ParseRecurrencyType(input)
	if err != nil {
		panic(err)
	}
	return t
}

func ParseRecurrencyTypes(inputs []string) (RecurrencyTypes, error) {
	result := make(RecurrencyTypes, 0, len(inputs))
	for _, input := range inputs {
		v, err := ParseRecurrencyType(input)
		if err != nil {
			return nil, err
		}
		result = append(result, v)
	}
	return result, nil
}

func ParseRecurrencyTypesOrDefault(inputs []string, defaultValue RecurrencyType) RecurrencyTypes {
	result := make(RecurrencyTypes, 0, len(inputs))
	for _, input := range inputs {
		result = append(result, ParseRecurrencyTypeOrDefault(input, defaultValue))
	}
	return result
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
