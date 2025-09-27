package label

import (
	ex "github.com/mistribe/subtracker/pkg/x/exception"
)

var (
	ErrLabelNotFound           = ex.NewNotFound("label not found")
	ErrLabelAlreadyExists      = ex.NewAlreadyExists("label already exists")
	ErrMissingDefaultLabel     = ex.NewInvalidValue("missing default label")
	ErrCustomLabelLimitReached = ex.NewInvalidValue("custom label limit reached")
)
