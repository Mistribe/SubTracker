package x_test

import (
	"testing"

	"github.com/oleexo/subtracker/pkg/x"
)

func TestMakeKey(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic_lowercase",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "uppercase_to_lowercase",
			input:    "HelloWorld",
			expected: "helloworld",
		},
		{
			name:     "spaces_to_hyphens",
			input:    "hello world",
			expected: "hello-world",
		},
		{
			name:     "multiple_spaces",
			input:    "hello   world",
			expected: "hello---world",
		},
		{
			name:     "mixed_case_and_spaces",
			input:    "  HelLo  WoRLd  ",
			expected: "--hello--world--",
		},
		{
			name:     "empty_string",
			input:    "",
			expected: "",
		},
		{
			name:     "only_spaces",
			input:    "     ",
			expected: "-----",
		},
		{
			name:     "special_characters",
			input:    "hello@# world!",
			expected: "hello@#-world!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := x.MakeKey(tt.input)
			if result != tt.expected {
				t.Errorf("MakeKey(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
