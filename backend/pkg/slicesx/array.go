package slicesx

func Select[TIn any, TOut any](source []TIn, f func(TIn) TOut) []TOut {
	result := make([]TOut, len(source))
	for i, v := range source {
		result[i] = f(v)
	}
	return result
}

func SelectErr[TIn any, TOut any](source []TIn, f func(TIn) (TOut, error)) ([]TOut, error) {
	result := make([]TOut, len(source))
	for i, v := range source {
		r, err := f(v)
		if err != nil {
			return nil, err
		}
		result[i] = r
	}
	return result, nil
}

func SelectMany[TIn any, TOut any](source []TIn, f func(TIn) []TOut) []TOut {
	result := make([]TOut, 0)
	for _, v := range source {
		result = append(result, f(v)...)
	}
	return result
}

func ToMap[TKey comparable, TOut any, TValue any](source []TValue,
	keySelector func(TValue) TKey,
	valueSelector func(TValue) TOut) map[TKey]TOut {
	result := make(map[TKey]TOut)
	for _, v := range source {
		result[keySelector(v)] = valueSelector(v)
	}

	return result
}
