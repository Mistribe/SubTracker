package provider

import (
	"errors"
)

var (
	ErrProviderNotFound      = errors.New("provider not found")
	ErrPlanNotFound          = errors.New("plan not found")
	ErrProviderAlreadyExists = errors.New("provider already exists")
	ErrProviderInvalid       = errors.New("provider is invalid")
	ErrOnlyOwnerCanEdit      = errors.New("only the owner can edit the provider")
	ErrPlanAlreadyExists     = errors.New("plan already exists")
	ErrPriceAlreadyExists    = errors.New("price already exists")
	ErrPriceNotFound         = errors.New("price not found")
)
