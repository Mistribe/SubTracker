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

func Match[TIn any, TOut any](result Result[TIn],
    success func(value TIn) TOut,
    failure func(err error) TOut) TOut {
    if result.IsSuccess() {
        return success(result.getValue())
    }
    return failure(result.getError())
}
