package option

func Some[TValue any](value TValue) Option[TValue] {
	return some[TValue]{
		value: value,
	}
}

type some[TValue any] struct {
	value TValue
}

func (s some[TValue]) getValue() TValue {
	return s.value
}

func (s some[TValue]) IsSome() bool {
	return true
}

func (s some[TValue]) IsNone() bool {
	return false
}

func (s some[TValue]) IfSome(action func(TValue)) {
	action(s.value)
}

func (s some[TValue]) IfNone(_ func()) {
}

func (s some[TValue]) Value() *TValue {
	return &s.value
}

func (s some[TValue]) ValueWithDefault(_ TValue) TValue {
	return s.value
}

func (s some[TValue]) Transform(f func(TValue) TValue) Option[TValue] {
	return Some(f(s.value))
}
