package x_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mistribe/subtracker/pkg/x"
)

func TestValueOrDefault(t *testing.T) {
	t.Run("When input is nil, return default value", func(t *testing.T) {
		var input *string = nil
		var defaultValue = "default"
		result := x.ValueOrDefault(input, defaultValue)

		assert.Equal(t, result, defaultValue)
	})

	t.Run("When input is not nil, return input value", func(t *testing.T) {
		var input = "test"
		var defaultValue = "default"

		result := x.ValueOrDefault(&input, defaultValue)

		assert.Equal(t, result, input)
	})
}
