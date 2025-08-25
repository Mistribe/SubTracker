package validation_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mistribe/subtracker/pkg/x/validation"
)

func TestNewError(t *testing.T) {
	err := validation.NewError("email", "must be valid")

	assert.Equal(t, "email", err.Field(), "Field() should return the provided field")
	assert.Equal(t, "must be valid", err.Message(), "Message() should return the provided message")
	assert.Equal(t, "validation error for field 'email': must be valid", err.Error(),
		"Error() should format the error message correctly")
}

func TestErrors_Error_Empty(t *testing.T) {
	var errs validation.Errors
	assert.Equal(t, "", errs.Error(), "Error() on empty Errors should return an empty string")
}

func TestErrors_Error_Multiple(t *testing.T) {
	e1 := validation.NewError("name", "required")
	e2 := validation.NewError("age", "invalid")

	errs := validation.Errors{e1, e2}

	expected := e1.Error() + "; " + e2.Error()
	assert.Equal(t, expected, errs.Error(), "Error() should join individual error messages with '; '")
}

func TestErrors_HasErrors(t *testing.T) {
	var empty validation.Errors
	assert.False(t, empty.HasErrors(), "HasErrors() should be false for empty Errors")

	nonEmpty := validation.Errors{validation.NewError("field", "reason")}
	assert.True(t, nonEmpty.HasErrors(), "HasErrors() should be true for non-empty Errors")
}
