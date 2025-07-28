package provider

import (
	"errors"
)

var (
	ErrProviderNotFound      = errors.New("provider not found")
	ErrProviderAlreadyExists = errors.New("provider already exists")
	ErrProviderInvalid       = errors.New("provider is invalid")
	ErrOnlyOwnerCanEdit      = errors.New("only the owner can edit the provider")
)
