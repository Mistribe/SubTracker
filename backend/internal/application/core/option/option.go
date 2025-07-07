package option

type Option[TValue any] interface {
	getValue() TValue

	IsSome() bool
	IsNone() bool

	IfSome(action func(TValue))
	IfNone(action func())
}

func Bind[TIn any, TOut any](option Option[TIn],
	some func(TIn) Option[TOut]) Option[TOut] {
	if option.IsSome() {
		return some(option.getValue())
	}
	return None[TOut]()
}

func Match[TIn any, TOut any](option Option[TIn],
	some func(TIn) TOut,
	none func() TOut) TOut{
	if option.IsSome() {
		return some(option.getValue())
	}
	return none()
}