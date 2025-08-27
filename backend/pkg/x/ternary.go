package x

func Ternary[T any](condition bool, ifTrue T, ifFalse T) T {
	if condition {
		return ifTrue
	}
	return ifFalse
}

func TernaryFunc[T any](condition bool, ifTrue func() T, ifFalse func() T) T {
	if condition {
		return ifTrue()
	}

	return ifFalse()
}

func TernaryFunc2[T1 any, T2 any](condition bool, ifTrue func() (T1, T2), ifFalse func() (T1, T2)) (T1, T2) {
	if condition {
		return ifTrue()
	}

	return ifFalse()
}
