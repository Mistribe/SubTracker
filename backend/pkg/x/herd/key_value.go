package herd

type KeyValuePair[TKey comparable, TValue any] struct {
	Key   TKey
	Value TValue
}
