package herd

func Take[T any](source []T, count int) []T {
	if count >= len(source) {
		return source
	}

	if count <= 0 {
		return []T{}
	}
	result := make([]T, count)
	copy(result, source[:count])
	return result
}

func (e Enumerator[T]) Take(c int) Enumerator[T] {
	return func(yield func(T) bool) {
		taken := 0
		e(func(item T) bool {
			if taken >= c {
				return false
			}
			taken++
			return yield(item)
		})
	}
}
