package langext

// Equatable represents types that can be compared to equality
type Equatable[T any] interface {
	Equal(other T) bool
}
