package validationx

import (
	"fmt"
	"strings"
)

type Error interface {
	error

	Field() string
	Message() string
}

type Errors []Error

func (e Errors) Error() string {
	if len(e) == 0 {
		return ""
	}

	var messages []string
	for _, err := range e {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, "; ")
}

func (e Errors) HasErrors() bool {
	return len(e) > 0
}

type validationError struct {
	field   string
	message string
}

func (e validationError) Field() string {
	return e.field
}

func (e validationError) Message() string {
	return e.message
}

func (e validationError) Error() string {
	return fmt.Sprintf("validation error for field '%s': %s", e.field, e.message)
}

func NewError(field, message string) Error {
	return &validationError{
		field:   field,
		message: message,
	}
}
