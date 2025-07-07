package result

func Success[V any](value V) Result[V] {
	return success[V]{
		value: value,
	}
}

type success[V any] struct {
	value V
}

func (s success[V]) getValue() V {
	return s.value
}

func (s success[V]) getError() error {
	panic("cannot retrieve error in success case")
}

func (s success[V]) IsFaulted() bool {
	return false
}

func (s success[V]) IsSuccess() bool {
	return false
}

func (s success[V]) IfSuccess(action func(value V)) {
	action(s.value)
}

func (s success[V]) IfFailure(_ func(err error)) {
}
