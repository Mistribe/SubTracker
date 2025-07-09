package subscription

import (
	"errors"
	"fmt"
)

var (
	ErrSubscriptionNotFound             = errors.New("subscription not found")
	ErrSubscriptionAlreadyExists        = errors.New("subscription already exists")
	ErrSubscriptionNameTooShort         = errors.New("subscription name should be more than one character")
	ErrSubscriptionPaymentMissing       = errors.New("subscription should have atleast one payment")
	ErrSubscriptionLabelExcedeed        = fmt.Errorf("subscription's labels have excedeed the maximum of %d", MaxLabelCount)
	ErrSubscriptionFamilyMemberExcedeed = fmt.Errorf("subscription's family members have excedeed the maximum of %d", MaxFamilyMemberCount)
	ErrPaymentCannotEndBeforeStart      = errors.New("payment end date cannot be before the starting date")
)
