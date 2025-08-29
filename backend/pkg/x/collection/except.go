package collection

// Except returns a slice containing elements in 'a' that are not present in 'b'.
// Elements are compared using their natural equality defined by the comparable interface.
func Except[T comparable](a, b []T) []T {
	result := make([]T, 0)
	if len(a) == 0 {
		return result
	}
	if len(b) == 0 {
		return append(result, a...)
	}

	lookup := make(map[T]struct{})
	for _, item := range b {
		lookup[item] = struct{}{}
	}

	for _, item := range a {
		if _, exists := lookup[item]; !exists {
			result = append(result, item)
		}
	}

	return result
}
