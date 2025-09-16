package herd_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnumerator_Skip(t *testing.T) {
	t.Run("Skip skips the first n elements", func(t *testing.T) {
		src := makeEnumerator([]int{1, 2, 3, 4, 5})
		skipped := src.Skip(2)
		var got []int
		skipped(func(i int) bool {
			got = append(got, i)
			return true
		})
		assert.Equal(t, []int{3, 4, 5}, got)
	})

	t.Run("Skip with n=0 returns all elements", func(t *testing.T) {
		src := makeEnumerator([]int{1, 2, 3})
		skipped := src.Skip(0)
		var got []int
		skipped(func(i int) bool {
			got = append(got, i)
			return true
		})
		assert.Equal(t, []int{1, 2, 3}, got)
	})

	t.Run("Skip with n greater than length returns empty", func(t *testing.T) {
		src := makeEnumerator([]int{1, 2})
		skipped := src.Skip(5)
		var got []int
		skipped(func(i int) bool {
			got = append(got, i)
			return true
		})
		assert.Empty(t, got)
	})

	t.Run("Skip on empty enumerator returns empty", func(t *testing.T) {
		src := makeEnumerator([]int{})
		skipped := src.Skip(2)
		var got []int
		skipped(func(i int) bool {
			got = append(got, i)
			return true
		})
		assert.Empty(t, got)
	})

	t.Run("Skip supports early break in yield", func(t *testing.T) {
		src := makeEnumerator([]int{1, 2, 3, 4})
		skipped := src.Skip(1)
		var got []int
		skipped(func(i int) bool {
			got = append(got, i)
			return false // break after first yielded
		})
		assert.Equal(t, []int{2}, got)
	})
}
