package x

// P creates and returns a pointer to the given value of any type T.
func P[T any](value T) *T {
	return &value
}
