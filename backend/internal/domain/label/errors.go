package label

import "errors"

var (
	ErrLabelNotFound      = errors.New("label not found")
	ErrLabelAlreadyExists = errors.New("label already exists")
)
