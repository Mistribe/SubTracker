package subscription

import (
	"errors"
	"fmt"
)

var (
	ErrSubscriptionNotFound       = errors.New("subscription not found")
	ErrSubscriptionAlreadyExists  = errors.New("subscription already exists")
	ErrSubscriptionNameTooShort   = errors.New("subscription friendlyName should be more than one character")
	ErrSubscriptionPaymentMissing = errors.New("subscription should have at least one payment")
	ErrSubscriptionLabelExceeded  = fmt.Errorf("subscription's labels have exceeded the maximum of %d",
		MaxLabelCount)
	ErrSubscriptionFamilyMemberExceeded = fmt.Errorf("subscription's family members have exceeded the maximum of %d",
		MaxFamilyMemberPerSubscriptionCount)
	ErrPaymentCannotEndBeforeStart          = errors.New("payment end date cannot be before the starting date")
	ErrPaymentNotFound                      = errors.New("payment not found")
	ErrPaymentAlreadyExists                 = errors.New("payment already exists")
	ErrPaymentNotAlreadyExists              = errors.New("payment not already exists")
	ErrCannotHaveFamilyMembersWithoutFamily = errors.New("cannot have family members without a family")
	ErrPayerAndJointAccountConflict         = errors.New("subscription cannot have both payer and joint account payment")
	ErrNoFamilyDefined                      = errors.New("no family defined")
	ErrUnknownPayerType                     = errors.New("unknown payer type")
	ErrUnknownRecurrencyType                = errors.New("unknown recurrency type")
)
