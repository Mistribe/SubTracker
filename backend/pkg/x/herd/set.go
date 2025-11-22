package herd

type Set[T comparable] map[T]struct{}

func NewSet[T comparable]() Set[T] {
	return make(map[T]struct{})
}

func NewSetFromSlice[T comparable](values []T) Set[T] {
	result := NewSet[T]()
	for _, value := range values {
		result.Add(value)
	}
	return result
}

func (s Set[T]) Add(element T) {
	s[element] = struct{}{}
}

func (s Set[T]) Contains(element T) bool {
	_, exists := s[element]
	return exists
}

func (s Set[T]) Remove(element T) {
	delete(s, element)
}

func (s Set[T]) ToSlice() []T {
	result := make([]T, 0, len(s))
	for element := range s {
		result = append(result, element)
	}
	return result
}

func (s Set[T]) Len() int {
	return len(s)
}

func (s Set[T]) It() Enumerator[T] {
	return func(yield func(T) bool) {
		for element := range s {
			if !yield(element) {
				return
			}
		}
	}
}
