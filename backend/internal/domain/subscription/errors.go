package subscription

import (
	ex "github.com/mistribe/subtracker/pkg/x/exception"
)

var (
	ErrSubscriptionNotFound           = ex.NewNotFound("subscription not found")
	ErrSubscriptionAlreadyExists      = ex.NewAlreadyExists("subscription already exists")
	ErrUnknownPayerType               = ex.NewInvalidValue("unknown payer type")
	ErrUnknownRecurrencyType          = ex.NewInvalidValue("unknown recurrency type")
	ErrActiveSubscriptionLimitReached = ex.NewInvalidValue("active subscription limit reached")
	ErrPayerNeedFamily                = ex.NewInvalidValue("subscription need to be owned by a family to set a payer")
	ErrUnknownStatus                  = ex.NewInvalidValue("unknown status")
)
