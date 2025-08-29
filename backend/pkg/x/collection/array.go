package collection

func ToMap[TKey comparable, TOut any, TValue any](
	source []TValue,
	keySelector func(TValue) TKey,
	valueSelector func(TValue) TOut) map[TKey]TOut {
	result := make(map[TKey]TOut)
	for _, v := range source {
		key := keySelector(v)
		value := valueSelector(v)
		result[key] = value
	}

	return result
}
