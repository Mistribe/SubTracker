package exception

type Code int

const (
	Unknown Code = iota
	NotFound
	AlreadyExists
	InvalidValue
	InvalidOperation
	InvalidState
	NotImplemented
	InternalError
	Unauthorized
	Unauthenticated
	DeadlineExceeded
	Conflict
)
