package auth

import (
	"errors"
)

var (
	ErrUnknownOwnerType = errors.New("unknown owner type")
	ErrUnknownUser      = errors.New("unknown userProfile")
	ErrNotAuthorized    = errors.New("not authorized")
)
