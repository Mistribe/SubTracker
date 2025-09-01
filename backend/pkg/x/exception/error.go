package exception

type Exception interface {
	error
	Code() Code
}

type exceptions []Exception

type exception struct {
	message    string
	code       Code
	innerError error
}

func (e exception) Error() string {
	return e.message
}

func (e exception) Code() Code {
	return e.code
}

type Options struct {
	Code       Code
	InnerError error
}

func WithCode(code Code) func(options *Options) {
	return func(options *Options) {
		options.Code = code
	}
}

func WithInnerError(err error) func(options *Options) {
	return func(options *Options) {
		options.InnerError = err
	}
}

func New(message string, functions ...func(options *Options)) Exception {
	var options Options
	for _, f := range functions {
		f(&options)
	}
	return exception{
		message:    message,
		code:       options.Code,
		innerError: options.InnerError,
	}
}

func NewNotFound(message string, functions ...func(options *Options)) Exception {
	return New(message, append(functions, WithCode(NotFound))...)
}

func NewAlreadyExists(message string, functions ...func(options *Options)) Exception {
	return New(message, append(functions, WithCode(AlreadyExists))...)
}

func NewUnauthorized(message string, functions ...func(options *Options)) Exception {
	return New(message, append(functions, WithCode(Unauthorized))...)
}

func NewInvalidValue(message string, functions ...func(options *Options)) Exception {
	return New(message, append(functions, WithCode(InvalidValue))...)
}

func NewUnauthenticated(message string, functions ...func(options *Options)) Exception {
	return New(message, append(functions, WithCode(Unauthenticated))...)
}
