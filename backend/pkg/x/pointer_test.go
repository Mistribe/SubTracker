package x_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mistribe/subtracker/pkg/x"
)

func TestPointerFunc(t *testing.T) {
	t.Run("string pointer", func(t *testing.T) {
		value := "test"
		ptr := x.Pointer(value)
		assert.NotNil(t, ptr)
		assert.Equal(t, value, *ptr)
	})

	t.Run("integer pointer", func(t *testing.T) {
		value := 42
		ptr := x.Pointer(value)
		assert.NotNil(t, ptr)
		assert.Equal(t, value, *ptr)
	})

	t.Run("boolean pointer", func(t *testing.T) {
		value := true
		ptr := x.Pointer(value)
		assert.NotNil(t, ptr)
		assert.Equal(t, value, *ptr)
	})

	t.Run("struct pointer", func(t *testing.T) {
		type testStruct struct {
			field string
		}
		value := testStruct{field: "test"}
		ptr := x.Pointer(value)
		assert.NotNil(t, ptr)
		assert.Equal(t, value, *ptr)
	})
}
