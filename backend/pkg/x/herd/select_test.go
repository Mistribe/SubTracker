package herd_test

import (
	"errors"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mistribe/subtracker/pkg/x/herd"
)

func TestSelect(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		arr := []int{1, 2, 3}
		want := []int{2, 3, 4}

		got := herd.Select(arr, func(i int) int {
			return i + 1
		})

		if !slices.Equal(got, want) {
			t.Fatalf("Select() = %v, want %v", got, want)
		}
	})

	t.Run("empty", func(t *testing.T) {
		arr := []int{}
		want := []int{}

		got := herd.Select(arr, func(i int) int {
			return i + 1
		})

		assert.Empty(t, got)
		assert.Equal(t, want, got)
	})

	t.Run("nil", func(t *testing.T) {
		var arr []int
		var want []int

		got := herd.Select(arr, func(i int) int {
			return i + 1
		})

		assert.Empty(t, got)
		assert.Equal(t, want, got)
	})
}

func TestSelectErr(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		arr := []int{1, 2, 3}
		want := []int{2, 3, 4}

		got, err := herd.SelectErr(arr, func(i int) (int, error) {
			return i + 1, nil
		})

		assert.NoError(t, err)
		if !slices.Equal(got, want) {
			t.Fatalf("Select() = %v, want %v", got, want)
		}
	})

	t.Run("err", func(t *testing.T) {
		arr := []int{1, 2, 3}

		_, err := herd.SelectErr(arr, func(i int) (int, error) {
			return i + 1, errors.New("test error")
		})

		assert.Error(t, err)
	})

	t.Run("nil", func(t *testing.T) {
		var arr []int
		var want []int

		got, err := herd.SelectErr(arr, func(i int) (int, error) {
			return i + 1, nil
		})

		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("empty", func(t *testing.T) {
		arr := []int{}
		want := []int{}

		got, err := herd.SelectErr(arr, func(i int) (int, error) {
			return i + 1, nil
		})

		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
}

func TestSelectMany(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		arr := [][]int{
			{1, 2, 3},
			{4, 5, 6},
		}
		want := []int{1, 2, 3, 4, 5, 6}

		got := herd.SelectMany(arr, func(values []int) []int {
			return values
		})

		if !slices.Equal(got, want) {
			t.Fatalf("SelectMany() = %v, want %v", got, want)
		}
	})
}

func TestSelectManyErr(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		arr := [][]int{
			{1, 2, 3},
			{4, 5, 6},
		}
		want := []int{1, 2, 3, 4, 5, 6}

		got, err := herd.SelectManyErr(arr, func(values []int) ([]int, error) {
			return values, nil
		})

		assert.NoError(t, err)
		if !slices.Equal(got, want) {
			t.Fatalf("SelectMany() = %v, want %v", got, want)
		}
	})
}
