package slicesx

func MapToArr[TKey comparable, TValue any, TOut any](source map[TKey]TValue,
	f func(key TKey, value TValue) TOut) []TOut {
	result := make([]TOut, len(source))

	c := 0
	for k, v := range source {
		result[c] = f(k, v)
		c++
	}

	return result
}
