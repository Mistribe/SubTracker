package label

import "errors"

var (
	ErrLabelNotFound      = errors.New("label not found")
	ErrLabelAlreadyExists = errors.New("label already exists")
	ErrLabelNameEmpty     = errors.New("label name cannot be empty")
	ErrLabelNameTooLong   = errors.New("label name is too long (max 100 characters)")
	ErrLabelColorInvalid  = errors.New("label color must be a valid hex color (e.g., #FF0000)")
)
