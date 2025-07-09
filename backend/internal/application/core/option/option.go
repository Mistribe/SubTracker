package option

type Option[TValue any] interface {
	getValue() TValue

	IsSome() bool
	IsNone() bool

	IfSome(action func(TValue))
	IfNone(action func())

	Value() *TValue
}

func Bind[TIn any, TOut any](option Option[TIn],
	some func(TIn) Option[TOut]) Option[TOut] {
	if option.IsSome() {
		return some(option.getValue())
	}
	return None[TOut]()
}

func Map[TIn any, TOut any](option Option[TIn],
	some func(TIn) TOut) Option[TOut] {
	return Bind(option, func(in TIn) Option[TOut] {
		return Some(some(in))
	})
}

func Match[TIn any, TOut any](option Option[TIn],
	some func(TIn) TOut,
	none func() TOut) TOut {
	if option.IsSome() {
		return some(option.getValue())
	}
	return none()
}
