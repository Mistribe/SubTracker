package herd

type Enumerator[T any] func(func(T) bool)

func (e Enumerator[T]) ToSlice() []T {
	var items []T
	for item := range e {
		items = append(items, item)
	}
	return items
}

func (e Enumerator[T]) ToList() List[T] {
	return e.ToSlice()
}
