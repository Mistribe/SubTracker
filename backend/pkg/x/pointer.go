package x

// P returns a pointer to the given value of any type T.
func P[T any](value T) *T {
	return &value
}

// Pointer returns a pointer to the given value of any type T.
func Pointer[T any](value T) *T {
	return P(value)
}

func Int64Pointer(value int64) *int64 {
	return &value
}

func Int64P(value int64) *int64 {
	return Int64Pointer(value)
}
