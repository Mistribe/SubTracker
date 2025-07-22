package slicesx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExcept(t *testing.T) {
	t.Run("No exclusion when b is empty", func(t *testing.T) {
		a := []int{1, 2, 3}
		b := []int{}
		expected := []int{1, 2, 3}
		assert.Equal(t, expected, Except(a, b))
	})

	t.Run("All elements excluded when b contains all elements of a", func(t *testing.T) {
		a := []int{1, 2, 3}
		b := []int{1, 2, 3}
		expected := []int{}
		assert.Equal(t, expected, Except(a, b))
	})

	t.Run("Partial exclusion when b contains some elements of a", func(t *testing.T) {
		a := []int{1, 2, 3}
		b := []int{2}
		expected := []int{1, 3}
		assert.Equal(t, expected, Except(a, b))
	})

	t.Run("Empty result when a is empty", func(t *testing.T) {
		a := []int{}
		b := []int{1, 2, 3}
		expected := []int{}
		assert.Equal(t, expected, Except(a, b))
	})

	t.Run("Duplicates in a are kept if not excluded", func(t *testing.T) {
		a := []int{1, 2, 2, 3, 3, 3}
		b := []int{3}
		expected := []int{1, 2, 2}
		assert.Equal(t, expected, Except(a, b))
	})

	t.Run("Duplicates in b are ignored during exclusion", func(t *testing.T) {
		a := []int{1, 2, 3}
		b := []int{3, 3, 3}
		expected := []int{1, 2}
		assert.Equal(t, expected, Except(a, b))
	})

	t.Run("No exclusion when a and b have no common elements", func(t *testing.T) {
		a := []int{1, 2, 3}
		b := []int{4, 5, 6}
		expected := []int{1, 2, 3}
		assert.Equal(t, expected, Except(a, b))
	})

	t.Run("Empty result when both a and b are empty", func(t *testing.T) {
		a := []int{}
		b := []int{}
		expected := []int{}
		assert.Equal(t, expected, Except(a, b))
	})
}
