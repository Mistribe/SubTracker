package herd_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mistribe/subtracker/pkg/x/herd"
)

func intComparer(a, b int) bool { return a == b }

func TestList_Add(t *testing.T) {
	t.Run("Add appends item to list", func(t *testing.T) {
		l := herd.NewList[int]()
		l.Add(1)
		l.Add(2)
		assert.Equal(t, []int{1, 2}, l.ToSlice())
	})
}

func TestList_Remove(t *testing.T) {
	t.Run("Remove deletes item and returns true", func(t *testing.T) {
		l := herd.NewList[int]()
		l.Add(1)
		l.Add(2)
		ok := l.Remove(1, intComparer)
		assert.True(t, ok)
		assert.Equal(t, []int{2}, l.ToSlice())
	})
	t.Run("Remove returns false if item not found", func(t *testing.T) {
		l := herd.NewList[int]()
		l.Add(1)
		ok := l.Remove(2, intComparer)
		assert.False(t, ok)
		assert.Equal(t, []int{1}, l.ToSlice())
	})
}

func TestList_Contains(t *testing.T) {
	t.Run("Contains returns true if item present", func(t *testing.T) {
		l := herd.NewList[int]()
		l.Add(1)
		assert.True(t, l.Contains(1, intComparer))
	})
	t.Run("Contains returns false if item not present", func(t *testing.T) {
		l := herd.NewList[int]()
		l.Add(1)
		assert.False(t, l.Contains(2, intComparer))
	})
}

func TestList_Clear(t *testing.T) {
	t.Run("Clear removes all items", func(t *testing.T) {
		l := herd.NewList[int]()
		l.Add(1)
		l.Add(2)
		l.Clear()
		assert.Equal(t, 0, l.Len())
		assert.Equal(t, []int{}, l.ToSlice())
	})
}

func TestList_ToSlice(t *testing.T) {
	t.Run("ToSlice returns all items", func(t *testing.T) {
		l := herd.NewList[int]()
		l.Add(1)
		l.Add(2)
		s := l.ToSlice()
		assert.Equal(t, []int{1, 2}, s)
	})
}

func TestList_Len(t *testing.T) {
	t.Run("Len returns correct count", func(t *testing.T) {
		l := herd.NewList[int]()
		assert.Equal(t, 0, l.Len())
		l.Add(1)
		l.Add(2)
		assert.Equal(t, 2, l.Len())
	})
}

func TestList_It(t *testing.T) {
	t.Run("It iterates all items", func(t *testing.T) {
		l := herd.NewList[int]()
		l.Add(1)
		l.Add(2)
		got := []int{}
		l.It()(func(v int) bool {
			got = append(got, v)
			return true
		})
		assert.ElementsMatch(t, []int{1, 2}, got)
	})
	t.Run("It stops iteration when false is returned", func(t *testing.T) {
		l := herd.NewList[int]()
		l.Add(1)
		l.Add(2)
		count := 0
		l.It()(func(v int) bool {
			count++
			return false
		})
		assert.Equal(t, 1, count)
	})
	t.Run("It does nothing on empty list", func(t *testing.T) {
		l := herd.NewList[int]()
		called := false
		for range l.It() {
			called = true
		}
		assert.False(t, called)
	})
}
