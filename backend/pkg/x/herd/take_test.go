package herd_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mistribe/subtracker/pkg/x/herd"
)

func TestTake(t *testing.T) {
	t.Run("should return n first elements", func(t *testing.T) {
		arr := []int{1, 2, 3, 4, 5}
		got := herd.Take(arr, 3)
		assert.Equal(t, []int{1, 2, 3}, got)
	})

	t.Run("should return empty slice when n is 0", func(t *testing.T) {
		arr := []int{1, 2, 3}
		got := herd.Take(arr, 0)
		assert.Empty(t, got)
	})

	t.Run("should return nil when source is nil", func(t *testing.T) {
		var arr []int
		got := herd.Take(arr, 3)
		assert.Nil(t, got)
	})

	t.Run("should return empty slice when source is empty", func(t *testing.T) {
		arr := []int{}
		got := herd.Take(arr, 3)
		assert.Empty(t, got)
	})

	t.Run("should return all elements when n is greater than slice length", func(t *testing.T) {
		arr := []int{1, 2, 3}
		got := herd.Take(arr, 5)
		assert.Equal(t, []int{1, 2, 3}, got)
	})
}

func TestEnumerator_Take(t *testing.T) {
	t.Run("Take returns n first elements", func(t *testing.T) {
		src := makeEnumerator([]int{1, 2, 3, 4, 5})
		taken := src.Take(3)
		var got []int
		taken(func(i int) bool {
			got = append(got, i)
			return true
		})
		assert.Equal(t, []int{1, 2, 3}, got)
	})

	t.Run("Take with n=0 returns empty", func(t *testing.T) {
		src := makeEnumerator([]int{1, 2, 3})
		taken := src.Take(0)
		var got []int
		taken(func(i int) bool {
			got = append(got, i)
			return true
		})
		assert.Empty(t, got)
	})

	t.Run("Take with n greater than length returns all", func(t *testing.T) {
		src := makeEnumerator([]int{1, 2, 3})
		taken := src.Take(5)
		var got []int
		taken(func(i int) bool {
			got = append(got, i)
			return true
		})
		assert.Equal(t, []int{1, 2, 3}, got)
	})

	t.Run("Take on empty enumerator returns empty", func(t *testing.T) {
		src := makeEnumerator([]int{})
		taken := src.Take(3)
		var got []int
		taken(func(i int) bool {
			got = append(got, i)
			return true
		})
		assert.Empty(t, got)
	})

	t.Run("Take supports early break in yield", func(t *testing.T) {
		src := makeEnumerator([]int{1, 2, 3, 4})
		taken := src.Take(3)
		var got []int
		taken(func(i int) bool {
			got = append(got, i)
			return false // break after first yielded
		})
		assert.Equal(t, []int{1}, got)
	})
}
