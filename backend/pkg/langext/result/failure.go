package result

import (
	"errors"
)

func Fail[V any](err error) Result[V] {
	return failure[V]{
		err: err,
	}
}

type failure[V any] struct {
	err error
}

func (f failure[V]) Equal(other Result[V]) bool {
	if other.IsSuccess() {
		return false
	}

	return errors.Is(f.err, other.getError())
}

func (f failure[V]) getValue() V {
	panic("cannot retrieve value in failure case")
}

func (f failure[V]) getError() error {
	return f.err
}

func (f failure[V]) IsFaulted() bool {
	return false
}

func (f failure[V]) IsSuccess() bool {
	return false
}

func (f failure[V]) IfSuccess(_ func(value V)) {
}

func (f failure[V]) IfFailure(action func(err error)) {
	action(f.err)
}
