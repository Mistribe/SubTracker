package slicesx

// Intersect finds the intersection of two slices and returns a new slice containing elements present in both slices.
// It accepts two input slices of a comparable type and performs the operation in a generic manner.
// The function returns an empty slice if either input slice is empty or if there are no common elements.
func Intersect[T comparable](a, b []T) []T {
	result := make([]T, 0)
	if len(a) == 0 || len(b) == 0 {
		return result
	}

	lookup := make(map[T]struct{})
	for _, item := range a {
		lookup[item] = struct{}{}
	}

	processed := make(map[T]struct{})
	for _, item := range b {
		if _, exists := lookup[item]; exists {
			if _, alreadyProcessed := processed[item]; !alreadyProcessed {
				result = append(result, item)
				processed[item] = struct{}{}
			}
		}
	}

	return result
}
