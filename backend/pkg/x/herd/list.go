package herd

type List[T any] []T

type ListComparer[T any] func(T, T) bool

func NewList[T any]() *List[T] {
	return &List[T]{}
}

func (l *List[T]) Add(item T) {
	*l = append(*l, item)
}

func (l *List[T]) Remove(item T, comparer ListComparer[T]) bool {
	for i, v := range *l {
		if comparer(v, item) {
			*l = append((*l)[:i], (*l)[i+1:]...)
			return true
		}
	}
	return false
}

func (l *List[T]) Contains(item T, comparer ListComparer[T]) bool {
	for _, v := range *l {
		if comparer(v, item) {
			return true
		}
	}
	return false
}

func (l *List[T]) Clear() {
	*l = (*l)[:0]
}

func (l *List[T]) ToSlice() []T {
	return *l
}

func (l *List[T]) Len() int {
	return len(*l)
}

func (l *List[T]) It() Enumerator[T] {
	return func(yield func(T) bool) {
		for _, v := range *l {
			if !yield(v) {
				return
			}
		}
	}
}
