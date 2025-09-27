package herd_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mistribe/subtracker/pkg/x/herd"
)

func makeEnumerator[T any](items []T) herd.Enumerator[T] {
	return func(yield func(T) bool) {
		for _, v := range items {
			if !yield(v) {
				return
			}
		}
	}
}

func TestEnumerator_Where(t *testing.T) {
	t.Run("Where filters items by predicate", func(t *testing.T) {
		src := makeEnumerator([]int{1, 2, 3, 4, 5})
		even := src.Where(func(i int) bool { return i%2 == 0 })
		var got []int
		even(func(i int) bool {
			got = append(got, i)
			return true
		})
		assert.Equal(t, []int{2, 4}, got)
	})

	t.Run("Where returns empty if no items match", func(t *testing.T) {
		src := makeEnumerator([]int{1, 3, 5})
		even := src.Where(func(i int) bool { return i%2 == 0 })
		var got []int
		even(func(i int) bool {
			got = append(got, i)
			return true
		})
		assert.Empty(t, got)
	})

	t.Run("Where returns all if predicate always true", func(t *testing.T) {
		src := makeEnumerator([]int{1, 2, 3})
		all := src.Where(func(i int) bool { return true })
		var got []int
		all(func(i int) bool {
			got = append(got, i)
			return true
		})
		assert.Equal(t, []int{1, 2, 3}, got)
	})

	t.Run("Where returns empty for empty enumerator", func(t *testing.T) {
		src := makeEnumerator([]int{})
		filtered := src.Where(func(i int) bool { return true })
		var got []int
		filtered(func(i int) bool {
			got = append(got, i)
			return true
		})
		assert.Empty(t, got)
	})

	t.Run("Where supports early break in yield", func(t *testing.T) {
		src := makeEnumerator([]int{1, 2, 3, 4})
		even := src.Where(func(i int) bool { return i%2 == 0 })
		var got []int
		even(func(i int) bool {
			got = append(got, i)
			return false // break after first match
		})
		assert.Equal(t, []int{2}, got)
	})
}
