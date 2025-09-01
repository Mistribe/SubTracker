package auth

import (
	ex "github.com/mistribe/subtracker/pkg/x/exception"
)

var (
	ErrUnknownOwnerType = ex.NewInvalidValue("unknown owner type")
	ErrUnknownUser      = ex.NewInvalidValue("unknown userProfile")
	ErrUnauthorized     = ex.NewUnauthorized("unauthorized")
)
