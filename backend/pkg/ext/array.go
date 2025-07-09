package ext

func Map[TIn any, TOut any](source []TIn, f func(TIn) TOut) []TOut {
	result := make([]TOut, len(source))
	for i, v := range source {
		result[i] = f(v)
	}
	return result
}
