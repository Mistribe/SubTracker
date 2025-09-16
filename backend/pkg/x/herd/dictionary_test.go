package herd_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mistribe/subtracker/pkg/x/herd"
)

func TestDictionary_AddAndGet(t *testing.T) {
	t.Run("Add and Get returns correct value and exists", func(t *testing.T) {
		dict := herd.NewDictionary[string, int]()
		dict.Add("foo", 42)
		value, ok := dict.Get("foo")
		assert.True(t, ok)
		assert.Equal(t, 42, value)
	})
}

func TestDictionary_GetNotExists(t *testing.T) {
	t.Run("Get returns false for missing key", func(t *testing.T) {
		dict := herd.NewDictionary[string, int]()
		_, ok := dict.Get("bar")
		assert.False(t, ok)
	})
}

func TestDictionary_Remove(t *testing.T) {
	t.Run("Remove deletes the key", func(t *testing.T) {
		dict := herd.NewDictionary[string, int]()
		dict.Add("foo", 1)
		dict.Remove("foo")
		_, ok := dict.Get("foo")
		assert.False(t, ok)
	})
}

func TestDictionary_Contains(t *testing.T) {
	t.Run("Contains returns true for present key", func(t *testing.T) {
		dict := herd.NewDictionary[string, int]()
		dict.Add("foo", 1)
		assert.True(t, dict.Contains("foo"))
	})
	t.Run("Contains returns false for absent key", func(t *testing.T) {
		dict := herd.NewDictionary[string, int]()
		assert.False(t, dict.Contains("bar"))
	})
}

func TestDictionary_Len(t *testing.T) {
	t.Run("Len returns correct count", func(t *testing.T) {
		dict := herd.NewDictionary[string, int]()
		assert.Equal(t, 0, dict.Len())
		dict.Add("a", 1)
		dict.Add("b", 2)
		assert.Equal(t, 2, dict.Len())
	})
}

func TestDictionary_ToMap(t *testing.T) {
	t.Run("ToMap returns correct map", func(t *testing.T) {
		dict := herd.NewDictionary[string, int]()
		dict.Add("a", 1)
		dict.Add("b", 2)
		m := dict.ToMap()
		assert.Equal(t, map[string]int{"a": 1, "b": 2}, m)
	})
}

func TestDictionary_Clear(t *testing.T) {
	t.Run("Clear removes all entries", func(t *testing.T) {
		dict := herd.NewDictionary[string, int]()
		dict.Add("a", 1)
		dict.Add("b", 2)
		dict.Clear()
		assert.Equal(t, 0, dict.Len())
		assert.False(t, dict.Contains("a"))
	})
}

func TestDictionary_Keys(t *testing.T) {
	t.Run("Keys returns all keys", func(t *testing.T) {
		dict := herd.NewDictionary[string, int]()
		dict.Add("a", 1)
		dict.Add("b", 2)
		keys := dict.Keys()
		assert.ElementsMatch(t, []string{"a", "b"}, keys)
	})
}

func TestDictionary_Values(t *testing.T) {
	t.Run("Values returns all values", func(t *testing.T) {
		dict := herd.NewDictionary[string, int]()
		dict.Add("a", 1)
		dict.Add("b", 2)
		values := dict.Values()
		assert.ElementsMatch(t, []int{1, 2}, values)
	})
}

func TestDictionary_TryGet(t *testing.T) {
	t.Run("TryGet returns true and sets value if present", func(t *testing.T) {
		dict := herd.NewDictionary[string, int]()
		dict.Add("foo", 99)
		var v int
		ok := dict.TryGet("foo", &v)
		assert.True(t, ok)
		assert.Equal(t, 99, v)
	})
	t.Run("TryGet returns false and does not set value if absent", func(t *testing.T) {
		dict := herd.NewDictionary[string, int]()
		var v int = 123
		ok := dict.TryGet("bar", &v)
		assert.False(t, ok)
		assert.Equal(t, 123, v)
	})
}

func TestNewDictionaryFromSlice(t *testing.T) {
	t.Run("Creates dictionary from slice with key and value selectors", func(t *testing.T) {
		type pair struct {
			K string
			V int
		}
		input := []pair{{"a", 1}, {"b", 2}}
		dict := herd.NewDictionaryFromSlice(input, func(p pair) string { return p.K }, func(p pair) int { return p.V })
		assert.Equal(t, 2, dict.Len())
		v, ok := dict.Get("a")
		assert.True(t, ok)
		assert.Equal(t, 1, v)
		v, ok = dict.Get("b")
		assert.True(t, ok)
		assert.Equal(t, 2, v)
	})

	t.Run("Returns empty dictionary for empty slice", func(t *testing.T) {
		input := []int{}
		dict := herd.NewDictionaryFromSlice(input, func(i int) int { return i }, func(i int) int { return i })
		assert.Equal(t, 0, dict.Len())
	})

	t.Run("Handles duplicate keys (last wins)", func(t *testing.T) {
		type pair struct {
			K string
			V int
		}
		input := []pair{{"a", 1}, {"a", 2}, {"b", 3}}
		dict := herd.NewDictionaryFromSlice(input, func(p pair) string { return p.K }, func(p pair) int { return p.V })
		assert.Equal(t, 2, dict.Len())
		v, ok := dict.Get("a")
		assert.True(t, ok)
		assert.Equal(t, 2, v)
		v, ok = dict.Get("b")
		assert.True(t, ok)
		assert.Equal(t, 3, v)
	})

	t.Run("Value selector can transform data", func(t *testing.T) {
		input := []string{"foo", "bar"}
		dict := herd.NewDictionaryFromSlice(input, func(s string) string { return s },
			func(s string) int { return len(s) })
		assert.Equal(t, 3, dict["foo"])
		assert.Equal(t, 3, dict["bar"])
	})
}
