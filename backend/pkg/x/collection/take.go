package collection

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
