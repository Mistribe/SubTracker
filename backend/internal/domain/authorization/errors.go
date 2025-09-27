package authorization

import (
	"github.com/mistribe/subtracker/pkg/x/exception"
)

var (
	ErrUnauthorized = exception.NewUnauthorized("unauthorized")
)
