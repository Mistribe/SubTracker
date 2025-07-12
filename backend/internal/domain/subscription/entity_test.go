package subscription_test

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/subscription"
	"github.com/oleexo/subtracker/pkg/langext/option"

	"github.com/stretchr/testify/assert"
)

func createSubscription() subscription.Subscription {
	originalPayment := subscription.NewPaymentWithoutValidation(
		uuid.New(),
		42,
		time.Now().Add(-time.Hour*24*2),
		option.None[time.Time](),
		1,
		"USD",
		time.Now(),
		time.Now(),
	)
	sub := subscription.NewSubscriptionWithoutValidation(
		uuid.New(),
		"subscriptiontest",
		[]subscription.Payment{originalPayment},
		nil,
		nil,
		option.None[uuid.UUID](),
		time.Now(),
		time.Now(),
	)
	return sub
}

func TestAddPayment(t *testing.T) {
	t.Run("Add payment with a previous end date not set", func(t *testing.T) {
		sub := createSubscription()
		originalPaymentId := sub.Payments()[0].Id()
		payment := subscription.NewPaymentWithoutValidation(uuid.New(),
			42,
			time.Now().Add(-time.Hour*24),
			option.None[time.Time](),
			1,
			"USD",
			time.Now(),
			time.Now())

		sub.AddPayment(payment)

		assert.Len(t, sub.Payments(), 2)
		originalPayment := sub.GetPaymentById(originalPaymentId).Value()
		assert.NotNil(t, originalPayment)
		assert.NotNil(t, originalPayment.EndDate().Value())
	})

}
