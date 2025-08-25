package x

// P returns a pointer to the given value of any type T.
func P[T any](value T) *T {
	return &value
}

// Pointer returns a pointer to the given value of any type T.
func Pointer[T any](value T) *T {
	return P(value)
}
