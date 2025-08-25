package slicesx

type Set[T comparable] struct {
	elements map[T]struct{}
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		elements: make(map[T]struct{}),
	}
}

func (s *Set[T]) Add(element T) {
	s.elements[element] = struct{}{}
}

func (s *Set[T]) Contains(element T) bool {
	_, exists := s.elements[element]
	return exists
}

func (s *Set[T]) Remove(element T) {
	delete(s.elements, element)
}

func (s *Set[T]) ToSlice() []T {
	result := make([]T, 0, len(s.elements))
	for element := range s.elements {
		result = append(result, element)
	}
	return result
}

func (s *Set[T]) Len() int {
	return len(s.elements)
}
