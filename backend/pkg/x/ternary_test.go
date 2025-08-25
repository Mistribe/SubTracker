package x_test

import (
	"testing"

	"github.com/oleexo/subtracker/pkg/x"
)

func TestTernary(t *testing.T) {
	t.Run("condition_true_returns_ifTrue_string", func(t *testing.T) {
		result := x.Ternary(true, "hello", "world")
		expected := "hello"
		if result != expected {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("condition_false_returns_ifFalse_string", func(t *testing.T) {
		result := x.Ternary(false, "hello", "world")
		expected := "world"
		if result != expected {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("condition_true_returns_ifTrue_int", func(t *testing.T) {
		result := x.Ternary(true, 42, 0)
		expected := 42
		if result != expected {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("condition_false_returns_ifFalse_int", func(t *testing.T) {
		result := x.Ternary(false, 42, 0)
		expected := 0
		if result != expected {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("condition_true_returns_ifTrue_bool", func(t *testing.T) {
		result := x.Ternary(true, true, false)
		expected := true
		if result != expected {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("condition_false_returns_ifFalse_bool", func(t *testing.T) {
		result := x.Ternary(false, true, false)
		expected := false
		if result != expected {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("condition_true_with_nil_pointers", func(t *testing.T) {
		var ptr1 *int = nil
		var ptr2 = new(int)
		*ptr2 = 10

		result := x.Ternary(true, ptr1, ptr2)
		if result != ptr1 {
			t.Errorf("expected ptr1 (nil), got %v", result)
		}
	})

	t.Run("condition_false_with_nil_pointers", func(t *testing.T) {
		var ptr1 *int = nil
		var ptr2 = new(int)
		*ptr2 = 10

		result := x.Ternary(false, ptr1, ptr2)
		if result != ptr2 {
			t.Errorf("expected ptr2, got %v", result)
		}
	})
}

func TestTernaryFunc(t *testing.T) {
	t.Run("condition_true_returns_ifTrue_string", func(t *testing.T) {
		result := x.TernaryFunc(true, func() string { return "hello" }, func() string { return "world" })
		expected := "hello"
		if result != expected {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("condition_false_returns_ifFalse_string", func(t *testing.T) {
		result := x.TernaryFunc(false, func() string { return "hello" }, func() string { return "world" })
		expected := "world"
		if result != expected {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("condition_true_returns_ifTrue_int", func(t *testing.T) {
		result := x.TernaryFunc(true, func() int { return 42 }, func() int { return 0 })
		expected := 42
		if result != expected {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("condition_false_returns_ifFalse_int", func(t *testing.T) {
		result := x.TernaryFunc(false, func() int { return 42 }, func() int { return 0 })
		expected := 0
		if result != expected {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("condition_true_returns_ifTrue_pointer", func(t *testing.T) {
		var ptr1 *int = nil
		var ptr2 = new(int)
		*ptr2 = 10

		result := x.TernaryFunc(true, func() *int { return ptr1 }, func() *int { return ptr2 })
		if result != ptr1 {
			t.Errorf("expected ptr1 (nil), got %v", result)
		}
	})

	t.Run("condition_false_returns_ifFalse_pointer", func(t *testing.T) {
		var ptr1 *int = nil
		var ptr2 = new(int)
		*ptr2 = 10

		result := x.TernaryFunc(false, func() *int { return ptr1 }, func() *int { return ptr2 })
		if result != ptr2 {
			t.Errorf("expected ptr2, got %v", result)
		}
	})
}
