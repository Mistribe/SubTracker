package provider

import (
	ex "github.com/mistribe/subtracker/pkg/x/exception"
)

var (
	ErrProviderNotFound           = ex.NewNotFound("provider not found")
	ErrPlanNotFound               = ex.NewNotFound("plan not found")
	ErrProviderAlreadyExists      = ex.NewAlreadyExists("provider already exists")
	ErrOnlyOwnerCanEdit           = ex.NewUnauthorized("only the owner can edit the provider")
	ErrPlanAlreadyExists          = ex.NewAlreadyExists("plan already exists")
	ErrPriceAlreadyExists         = ex.NewAlreadyExists("price already exists")
	ErrPriceNotFound              = ex.NewNotFound("price not found")
	ErrCustomProviderLimitReached = ex.NewInvalidOperation("custom provider limit reached")
	ErrProviderIsInUsed           = ex.NewInvalidOperation("provider is in use")
)
