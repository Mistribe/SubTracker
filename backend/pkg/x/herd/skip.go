package herd

func (e Enumerator[T]) Skip(c int) Enumerator[T] {
	return func(yield func(T) bool) {
		skipped := 0
		e(func(item T) bool {
			if skipped < c {
				skipped++
				return true
			}
			return yield(item)
		})
	}
}
