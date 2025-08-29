package collection_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mistribe/subtracker/pkg/x/collection"
)

func TestTake(t *testing.T) {
	t.Run("should return n first elements", func(t *testing.T) {
		arr := []int{1, 2, 3, 4, 5}
		got := collection.Take(arr, 3)
		assert.Equal(t, []int{1, 2, 3}, got)
	})

	t.Run("should return empty slice when n is 0", func(t *testing.T) {
		arr := []int{1, 2, 3}
		got := collection.Take(arr, 0)
		assert.Empty(t, got)
	})

	t.Run("should return nil when source is nil", func(t *testing.T) {
		var arr []int
		got := collection.Take(arr, 3)
		assert.Nil(t, got)
	})

	t.Run("should return empty slice when source is empty", func(t *testing.T) {
		arr := []int{}
		got := collection.Take(arr, 3)
		assert.Empty(t, got)
	})

	t.Run("should return all elements when n is greater than slice length", func(t *testing.T) {
		arr := []int{1, 2, 3}
		got := collection.Take(arr, 5)
		assert.Equal(t, []int{1, 2, 3}, got)
	})
}
