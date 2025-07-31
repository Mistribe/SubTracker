package x_test

import (
	"testing"

	"github.com/oleexo/subtracker/pkg/x"

	"github.com/stretchr/testify/assert"
)

func TestMakeKey(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "String with spaces",
			input:    "test string",
			expected: "test string",
		},
		{
			name:     "String with numbers",
			input:    "123456",
			expected: "123456",
		},
		{
			name:     "Special characters",
			input:    "#$%^&**",
			expected: "#$%^&**",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			result := x.MakeKey(tc.input)
			assert.Equal(tc.expected, result)
		})
	}
}
