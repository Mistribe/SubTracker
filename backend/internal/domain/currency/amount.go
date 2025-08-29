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
	Add(value float64) Amount
	Sub(value float64) Amount
	Div(value float64) Amount
	Mul(value float64) Amount
}
type amount struct {
	value    float64
	currency Unit
	source   Amount
	isValid  bool
}

func (a amount) Add(value float64) Amount {
	if !a.isValid {
		return NewInvalidAmount()
	}

	var source Amount
	if a.source != nil {
		source = source.Add(value)
	}

	return amount{
		value:    a.value + value,
		currency: a.currency,
		source:   source,
		isValid:  true,
	}
}

func (a amount) Sub(value float64) Amount {
	if !a.isValid {
		return NewInvalidAmount()
	}
	var source Amount
	if a.source != nil {
		source = source.Sub(value)
	}

	return amount{
		value:    a.value - value,
		currency: a.currency,
		source:   source,
		isValid:  true,
	}
}

func (a amount) Div(value float64) Amount {
	if !a.isValid {
		return NewInvalidAmount()
	}
	var source Amount
	if a.source != nil {
		source = source.Div(value)
	}

	return amount{
		value:    a.value / value,
		currency: a.currency,
		source:   source,
		isValid:  true,
	}
}

func (a amount) Mul(value float64) Amount {
	if !a.isValid {
		return NewInvalidAmount()
	}
	var source Amount
	if a.source != nil {
		source = source.Mul(value)
	}

	return amount{
		value:    a.value * value,
		currency: a.currency,
		source:   source,
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
