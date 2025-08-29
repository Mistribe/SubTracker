package collection_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mistribe/subtracker/pkg/x/collection"
)

func TestSet(t *testing.T) {
	t.Run("Add_AddsNewElement", func(t *testing.T) {
		set := collection.NewSet[int]()
		set.Add(1)
		assert.True(t, set.Contains(1))
	})

	t.Run("Add_DuplicateElement", func(t *testing.T) {
		set := collection.NewSet[int]()
		set.Add(1)
		set.Add(1)
		assert.Equal(t, 1, set.Len())
	})

	t.Run("Contains_ExistingElement", func(t *testing.T) {
		set := collection.NewSet[string]()
		set.Add("test")
		assert.True(t, set.Contains("test"))
	})

	t.Run("Contains_NonExistingElement", func(t *testing.T) {
		set := collection.NewSet[string]()
		assert.False(t, set.Contains("test"))
	})

	t.Run("Remove_ExistingElement", func(t *testing.T) {
		set := collection.NewSet[int]()
		set.Add(1)
		set.Remove(1)
		assert.False(t, set.Contains(1))
	})

	t.Run("Remove_NonExistingElement", func(t *testing.T) {
		set := collection.NewSet[int]()
		set.Remove(1)
		assert.Equal(t, 0, set.Len())
	})

	t.Run("ToSlice_EmptySet", func(t *testing.T) {
		set := collection.NewSet[int]()
		slice := set.ToSlice()
		assert.Empty(t, slice)
	})

	t.Run("ToSlice_NonEmptySet", func(t *testing.T) {
		set := collection.NewSet[int]()
		set.Add(1)
		set.Add(2)
		assert.Len(t, set.ToSlice(), 2)
	})

	t.Run("Len_EmptySet", func(t *testing.T) {
		set := collection.NewSet[int]()
		assert.Equal(t, 0, set.Len())
	})

	t.Run("Len_NonEmptySet", func(t *testing.T) {
		set := collection.NewSet[int]()
		set.Add(1)
		set.Add(2)
		assert.Equal(t, 2, set.Len())
	})
}
