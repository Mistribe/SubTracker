package herd

type Dictionary[TKey comparable, TValue any] map[TKey]TValue

type DictionaryIterator[TKey comparable, TValue any] func(func(TKey, TValue) bool)

func NewDictionary[TKey comparable, TValue any]() Dictionary[TKey, TValue] {
	return make(map[TKey]TValue)
}

func NewDictionaryFromSlice[TSliceValue any, TKey comparable, TValue any](
	input []TSliceValue, keySelector func(TSliceValue) TKey,
	valueSelector func(TSliceValue) TValue) Dictionary[TKey, TValue] {
	result := make(map[TKey]TValue)
	for _, v := range input {
		key := keySelector(v)
		value := valueSelector(v)
		result[key] = value
	}
	return result
}

func (d Dictionary[TKey, TValue]) Add(key TKey, value TValue) {
	d[key] = value
}

func (d Dictionary[TKey, TValue]) Get(key TKey) (TValue, bool) {
	value, exists := d[key]
	return value, exists
}

func (d Dictionary[TKey, TValue]) Remove(key TKey) {
	delete(d, key)
}

func (d Dictionary[TKey, TValue]) Contains(key TKey) bool {
	_, exists := d[key]
	return exists
}

func (d Dictionary[TKey, TValue]) Len() int {
	return len(d)
}

func (d Dictionary[TKey, TValue]) ToMap() map[TKey]TValue {
	return d
}

func (d Dictionary[TKey, TValue]) Clear() {
	for k := range d {
		delete(d, k)
	}
}

func (d Dictionary[TKey, TValue]) Keys() []TKey {
	keys := make([]TKey, 0, len(d))
	for k := range d {
		keys = append(keys, k)
	}
	return keys
}

func (d Dictionary[TKey, TValue]) Values() []TValue {
	values := make([]TValue, 0, len(d))
	for _, v := range d {
		values = append(values, v)
	}
	return values
}

func (d Dictionary[TKey, TValue]) TryGet(key TKey, value *TValue) bool {
	v, exists := d[key]
	if exists {
		*value = v
	}
	return exists
}

func (d Dictionary[TKey, TValue]) It() Enumerator[KeyValuePair[TKey, TValue]] {
	return func(yield func(KeyValuePair[TKey, TValue]) bool) {
		for k, v := range d {
			if !yield(KeyValuePair[TKey, TValue]{Key: k, Value: v}) {
				return
			}
		}
	}
}
