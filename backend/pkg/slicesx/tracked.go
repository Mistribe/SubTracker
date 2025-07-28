package slicesx

import (
	"iter"
	"sort"
)

type Tracked[T comparable] struct {
	_added   []T
	_removed []T
	_updated []T
	_values  []T

	_uniqueCompareFunc func(T, T) bool
	_compareFunc       func(T, T) bool
}

func NewTracked[T comparable](
	values []T,
	uniqueCompareFunc func(T, T) bool,
	compareFunc func(T, T) bool) *Tracked[T] {
	return &Tracked[T]{
		_added:             nil,
		_removed:           nil,
		_updated:           nil,
		_values:            values,
		_uniqueCompareFunc: uniqueCompareFunc,
		_compareFunc:       compareFunc,
	}
}

func (d *Tracked[T]) HasChanges() bool {
	return len(d._added) > 0 || len(d._removed) > 0
}

func (d *Tracked[T]) Values() []T {
	return d._values
}

func (d *Tracked[T]) It() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range d._values {
			if !yield(v) {
				return
			}
		}
	}

}

func (d *Tracked[T]) Added() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range d._added {
			if !yield(v) {
				return
			}
		}
	}
}

func (d *Tracked[T]) Removed() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range d._removed {
			if !yield(v) {
				return
			}
		}
	}
}

func (d *Tracked[T]) Updated() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range d._updated {
			if !yield(v) {
				return
			}
		}
	}
}

func (d *Tracked[T]) Clear() {
	d._added = nil
	d._removed = nil
	d._updated = nil
	d._values = nil
}

func (d *Tracked[T]) ClearChanges() {
	d._added = nil
	d._removed = nil
	d._updated = nil
}

func (d *Tracked[T]) Len() int {
	return len(d._values)
}

func (d *Tracked[T]) GetAt(index int) T {
	return d._values[index]
}

func (d *Tracked[T]) Get(value T) (T, bool) {
	for _, v := range d._values {
		if d._uniqueCompareFunc(v, value) {
			return v, true
		}
	}
	return value, false
}

func (d *Tracked[T]) Set(values []T) {
	var newValues []T
	for _, vIn := range d._values {
		for _, vOut := range values {
			if d._uniqueCompareFunc(vIn, vOut) {
				if !d._compareFunc(vIn, vOut) {
					d._updated = append(d._updated, vOut)
				}
				newValues = append(newValues, vOut)
				return
			}
		}
	}

	addedElements := Except(values, d._values)
	d._added = append(d._added, addedElements...)
	newValues = append(newValues, addedElements...)
	deletedElements := Except(d._values, values)
	d._removed = append(d._removed, deletedElements...)
	d._values = newValues
}

func (d *Tracked[T]) Add(value T) bool {
	d._added = append(d._added, value)
	if d.Contains(value) {
		return false
	}
	d._values = append(d._values, value)
	return true
}

func (d *Tracked[T]) Contains(value T) bool {
	for _, v := range d._values {
		if d._uniqueCompareFunc(v, value) {
			return true
		}
	}
	return false
}

func (d *Tracked[T]) Remove(value T) bool {
	d._removed = append(d._removed, value)
	for idx, v := range d._values {
		if d._uniqueCompareFunc(v, value) {
			d._values = append(d._values[:idx], d._values[idx+1:]...)
			return true
		}
	}
	return false
}

func (d *Tracked[T]) Update(value T) bool {
	d._updated = append(d._updated, value)
	for idx, v := range d._values {
		if d._uniqueCompareFunc(v, value) {
			d._values[idx] = value
			return true
		}
	}
	return false
}

func (d *Tracked[T]) Equals(others []T) bool {
	if len(d._values) != len(others) {
		return false
	}
	comparedElementCount := 0
	for _, element := range d._values {
		for _, other := range others {
			if d._uniqueCompareFunc(element, other) {
				if !d._compareFunc(element, other) {
					return false
				} else {
					comparedElementCount++
				}
			}
		}
	}
	if comparedElementCount != len(d._values) {
		return false
	}
	return true
}

func (d *Tracked[T]) IsEmpty() bool {
	return len(d._values) == 0
}

func (d *Tracked[T]) IsNotEmpty() bool {
	return len(d._values) > 0
}

func (d *Tracked[T]) Sort(less func(T, T) bool) {
	sort.Slice(d._values, func(i, j int) bool {
		return less(d._values[i], d._values[j])
	})
}
