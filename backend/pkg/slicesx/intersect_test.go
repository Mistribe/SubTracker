package slicesx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntersect(t *testing.T) {
	t.Run("both_empty", func(t *testing.T) {
		a := []int{}
		b := []int{}
		expected := []int{}
		result := Intersect(a, b)
		assert.Equal(t, expected, result)
	})

	t.Run("first_empty", func(t *testing.T) {
		a := []int{}
		b := []int{1, 2, 3}
		expected := []int{}
		result := Intersect(a, b)
		assert.Equal(t, expected, result)
	})

	t.Run("second_empty", func(t *testing.T) {
		a := []int{1, 2, 3}
		b := []int{}
		expected := []int{}
		result := Intersect(a, b)
		assert.Equal(t, expected, result)
	})

	t.Run("no_common_elements", func(t *testing.T) {
		a := []int{1, 2, 3}
		b := []int{4, 5, 6}
		expected := []int{}
		result := Intersect(a, b)
		assert.Equal(t, expected, result)
	})

	t.Run("some_common_elements", func(t *testing.T) {
		a := []int{1, 2, 3}
		b := []int{3, 4, 5}
		expected := []int{3}
		result := Intersect(a, b)
		assert.Equal(t, expected, result)
	})

	t.Run("all_common_elements", func(t *testing.T) {
		a := []int{1, 2, 3}
		b := []int{1, 2, 3}
		expected := []int{1, 2, 3}
		result := Intersect(a, b)
		assert.Equal(t, expected, result)
	})

	t.Run("duplicates_in_a", func(t *testing.T) {
		a := []int{1, 2, 2, 3}
		b := []int{2, 3}
		expected := []int{2, 3}
		result := Intersect(a, b)
		assert.Equal(t, expected, result)
	})

	t.Run("duplicates_in_b", func(t *testing.T) {
		a := []int{1, 2, 3}
		b := []int{2, 2, 3}
		expected := []int{2, 3}
		result := Intersect(a, b)
		assert.Equal(t, expected, result)
	})

	t.Run("duplicates_in_both", func(t *testing.T) {
		a := []int{1, 2, 2, 3}
		b := []int{2, 3, 3}
		expected := []int{2, 3}
		result := Intersect(a, b)
		assert.Equal(t, expected, result)
	})

	t.Run("longer_lists_with_some_common", func(t *testing.T) {
		a := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
		b := []int{5, 6, 10, 11}
		expected := []int{5, 6}
		result := Intersect(a, b)
		assert.Equal(t, expected, result)
	})
}

func TestIntersectStrings(t *testing.T) {
	t.Run("string_intersection", func(t *testing.T) {
		a := []string{"a", "b", "c"}
		b := []string{"b", "d", "c"}
		expected := []string{"b", "c"}
		result := Intersect(a, b)
		assert.Equal(t, expected, result)
	})
}
