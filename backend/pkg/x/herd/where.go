package herd

func (e Enumerator[T]) Where(predicate func(T) bool) Enumerator[T] {
	return func(yield func(T) bool) {
		for item := range e {
			if predicate(item) {
				if !yield(item) {
					return
				}
			}
		}
	}
}
