package x_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mistribe/subtracker/pkg/x"
)

func TestParseOrNewUUID(t *testing.T) {
	t.Run("When input is nil, return new UUID", func(t *testing.T) {
		var input *string = nil

		result, err := x.ParseOrNewUUID(input)

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("When input is not nil, return parsed UUID", func(t *testing.T) {
		var input = "123e4567-e89b-12d3-a456-426655440000"

		result, err := x.ParseOrNewUUID(&input)

		assert.NoError(t, err)
		assert.Equal(t, input, result.String())
	})
}

func TestParseOrNilUUID(t *testing.T) {
	t.Run("When input is nil, return nil", func(t *testing.T) {
		var input *string = nil

		result, err := x.ParseOrNilUUID(input)

		assert.NoError(t, err)
		assert.Nil(t, result)
	})

	t.Run("When input is not nil, return parsed UUID", func(t *testing.T) {
		var input = "123e4567-e89b-12d3-a456-426655440000"

		result, err := x.ParseOrNilUUID(&input)

		assert.NoError(t, err)
		assert.Equal(t, input, result.String())
	})
}
