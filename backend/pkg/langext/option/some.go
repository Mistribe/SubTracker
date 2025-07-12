package option

import (
	"reflect"
)

func Some[TValue any](someValue TValue) Option[TValue] {
	return some[TValue]{
		value: someValue,
	}
}

type some[TValue any] struct {
	value TValue
}

func (s some[TValue]) getValue() TValue {
	return s.value
}

func (s some[TValue]) Equal(otherOption Option[TValue]) bool {
	if otherOption.IsNone() {
		return false
	}

	otherValue := otherOption.getValue()
	if currentValueEq, ok := any(s.value).(interface{ Equal(otherValue TValue) bool }); ok {
		return currentValueEq.Equal(otherValue)
	}

	return reflect.DeepEqual(s.value, otherValue)
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
