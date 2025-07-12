package result

import (
	"reflect"
)

func Success[V any](value V) Result[V] {
	return success[V]{
		value: value,
	}
}

type success[TValue any] struct {
	value TValue
}

func (s success[TValue]) Equal(otherResult Result[TValue]) bool {
	if otherResult.IsFaulted() {
		return false
	}
	otherValue := otherResult.getValue()
	if currentValueEq, ok := any(s.value).(interface{ Equal(otherValue TValue) bool }); ok {
		return currentValueEq.Equal(otherValue)
	}

	return reflect.DeepEqual(s.value, otherValue)
}

func (s success[TValue]) getValue() TValue {
	return s.value
}

func (s success[TValue]) getError() error {
	panic("cannot retrieve error in success case")
}

func (s success[TValue]) IsFaulted() bool {
	return false
}

func (s success[TValue]) IsSuccess() bool {
	return true
}

func (s success[TValue]) IfSuccess(action func(value TValue)) {
	action(s.value)
}

func (s success[TValue]) IfFailure(_ func(err error)) {
}
