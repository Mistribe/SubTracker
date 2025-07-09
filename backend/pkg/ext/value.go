package ext

func ValueOrDefault[T any](in *T, defaultValue T) T {
	if in == nil {
		return defaultValue
	}
	return *in
}
