package collection

func Select[TIn any, TOut any](source []TIn, f func(TIn) TOut) []TOut {
	if source == nil {
		return nil
	}
	result := make([]TOut, len(source))
	for i, v := range source {
		result[i] = f(v)
	}
	return result
}

func SelectErr[TIn any, TOut any](source []TIn, f func(TIn) (TOut, error)) ([]TOut, error) {
	if source == nil {
		return nil, nil
	}
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
	if source == nil {
		return nil
	}
	result := make([]TOut, 0)
	for _, v := range source {
		result = append(result, f(v)...)
	}
	return result
}

func SelectManyErr[TIn any, TOut any](source []TIn, f func(TIn) ([]TOut, error)) ([]TOut, error) {
	if source == nil {
		return nil, nil
	}
	result := make([]TOut, 0)
	for _, v := range source {
		r, err := f(v)
		if err != nil {
			return nil, err
		}
		result = append(result, r...)
	}

	return result, nil
}
