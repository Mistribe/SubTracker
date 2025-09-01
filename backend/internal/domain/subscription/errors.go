package subscription

import (
	ex "github.com/mistribe/subtracker/pkg/x/exception"
)

var (
	ErrSubscriptionNotFound      = ex.NewNotFound("subscription not found")
	ErrSubscriptionAlreadyExists = ex.NewAlreadyExists("subscription already exists")
	ErrUnknownPayerType          = ex.NewInvalidValue("unknown payer type")
	ErrUnknownRecurrencyType     = ex.NewInvalidValue("unknown recurrency type")
)
