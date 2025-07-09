package result

type Result[V any] interface {
	getValue() V
	getError() error
	IsFaulted() bool
	IsSuccess() bool

	IfSuccess(action func(value V))
	IfFailure(action func(err error))
}

func Bind[TIn any, TOut any](result Result[TIn],
	success func(value TIn) Result[TOut]) Result[TOut] {
	if result.IsSuccess() {
		return success(result.getValue())
	}
	return Fail[TOut](result.getError())
}

func BindReduce[TIn any, TOut any](results []Result[TIn],
	success func(values []TIn) Result[TOut]) Result[TOut] {
	values := make([]TIn, len(results))

	for idx, r := range results {
		if r.IsFaulted() {
			return Fail[TOut](r.getError())
		}
		values[idx] = r.getValue()
	}

	return success(values)
}

func Map[TIn any, TOut any](result Result[TIn],
	success func(value TIn) TOut) Result[TOut] {
	return Bind(result, func(value TIn) Result[TOut] {
		return Success(success(value))
	})
}

func Match[TIn any, TOut any](result Result[TIn],
	success func(value TIn) TOut,
	failure func(err error) TOut) TOut {
	if result.IsSuccess() {
		return success(result.getValue())
	}
	return failure(result.getError())
}

type Unit struct {
}

func Void() Result[Unit] {
	return Success(Unit{})
}
