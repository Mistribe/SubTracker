package herd

func (e Enumerator[T]) Count() int {
	var count int
	for range e {
		count++
	}
	return count
}
