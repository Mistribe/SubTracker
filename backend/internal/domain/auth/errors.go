package auth

import (
	"errors"
)

var (
	ErrUnknownOwnerType = errors.New("unknown owner type")
	ErrUnknownUser      = errors.New("unknown user")
	ErrNotAuthorized    = errors.New("not authorized")
)
