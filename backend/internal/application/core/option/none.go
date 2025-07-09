package option

func None[TValue any]() Option[TValue] {
	return none[TValue]{}
}

type none[TValue any] struct {
}

func (n none[TValue]) getValue() TValue {
	panic("cannot retrive value in none case")
}

func (n none[TValue]) IsSome() bool {
	return false
}

func (n none[TValue]) IsNone() bool {
	return true
}

func (n none[TValue]) IfSome(_ func(TValue)) {
}

func (n none[TValue]) IfNone(action func()) {
	action()
}

func (n none[TValue]) Value() *TValue {
	return nil
}

func (n none[TValue]) ValueWithDefault(value TValue) TValue {
	return value
}

func (n none[TValue]) Transform(f func(TValue) TValue) Option[TValue] {
	return n
}
