package user

import (
	"errors"
)

var (
	ErrUnknownOwnerType = errors.New("unknown owner type")
	ErrUnknownUser      = errors.New("unknown user")
)
