package currency

type Amount interface {
	IsZero() bool
	IsNegative() bool
	IsPositive() bool
	IsEqual(other Amount) bool
	IsGreaterThan(other Amount) bool
	IsLessThan(other Amount) bool
	Value() float64
	Source() Amount
	Currency() Unit
	IsValid() bool
	Add(otherAmount Amount) Amount
}
type amount struct {
	value    float64
	currency Unit
	source   Amount
	isValid  bool
}

func (a amount) Add(otherAmount Amount) Amount {
	if !a.isValid || !otherAmount.IsValid() {
		return NewInvalidAmount()
	}

	if a.currency != otherAmount.Currency() {
		return NewInvalidAmount()
	}
	return amount{
		value:    a.value + otherAmount.Value(),
		currency: a.currency,
		source:   nil,
		isValid:  true,
	}
}

func (a amount) Source() Amount {
	return a.source
}

func (a amount) IsZero() bool {
	return a.value == 0
}

func (a amount) IsNegative() bool {
	return a.value < 0
}

func (a amount) IsPositive() bool {
	return a.value > 0
}

func (a amount) IsEqual(other Amount) bool {
	return a.value == other.Value() &&
		a.currency == other.Currency()
}

func (a amount) IsGreaterThan(other Amount) bool {
	return a.value > other.Value()
}

func (a amount) IsLessThan(other Amount) bool {
	return a.value < other.Value()
}

func (a amount) Value() float64 {
	return a.value
}

func (a amount) Currency() Unit {
	return a.currency
}

func (a amount) IsValid() bool {
	return a.isValid
}

func NewAmount(value float64, unit Unit) Amount {
	return amount{
		value:    value,
		currency: unit,
		isValid:  true,
	}
}

func NewAmountWithSource(value float64, unit Unit, source Amount) Amount {
	return amount{
		value:    value,
		currency: unit,
		source:   source,
		isValid:  true,
	}
}

func NewInvalidAmount() Amount {
	return amount{
		isValid: false,
	}
}
