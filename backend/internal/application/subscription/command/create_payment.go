package command

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/subscription"
	"github.com/oleexo/subtracker/pkg/langext/option"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type CreatePaymentCommand struct {
	SubscriptionId uuid.UUID
	Payment        subscription.Payment
}

type CreatePaymentCommandHandler struct {
	subscriptionRepository subscription.Repository
}

func NewCreatePaymentCommandHandler(subscriptionRepository subscription.Repository) *CreatePaymentCommandHandler {
	return &CreatePaymentCommandHandler{subscriptionRepository: subscriptionRepository}
}

func (h CreatePaymentCommandHandler) Handle(
	ctx context.Context,
	command CreatePaymentCommand) result.Result[subscription.Subscription] {
	subOpt, err := h.subscriptionRepository.Get(ctx, command.SubscriptionId)
	if err != nil {
		return result.Fail[subscription.Subscription](err)
	}

	return option.Match(subOpt,
		func(sub subscription.Subscription) result.Result[subscription.Subscription] {
			return h.createPayment(ctx, command.Payment, sub)
		}, func() result.Result[subscription.Subscription] {
			return result.Fail[subscription.Subscription](subscription.ErrSubscriptionNotFound)
		},
	)
}

func (h CreatePaymentCommandHandler) createPayment(
	ctx context.Context,
	payment subscription.Payment,
	sub subscription.Subscription) result.Result[subscription.Subscription] {

	if err := payment.Validate(); err != nil {
		return result.Fail[subscription.Subscription](err)
	}

	if existingPaymentOpt := sub.GetPaymentById(payment.Id()); existingPaymentOpt.IsSome() {
		return option.Match(existingPaymentOpt,
			func(existingPayment subscription.Payment) result.Result[subscription.Subscription] {
				if existingPayment.Equal(payment) {
					return result.Success(sub)
				}
				return result.Fail[subscription.Subscription](subscription.ErrPaymentAlreadyExists)
			},
			func() result.Result[subscription.Subscription] {
				return result.Fail[subscription.Subscription](subscription.ErrPaymentNotFound)
			})
	}

	sub.AddPayment(payment)

	sub.SetUpdatedAt(time.Now())

	if err := sub.Validate(); err != nil {
		return result.Fail[subscription.Subscription](err)
	}

	if err := h.subscriptionRepository.Save(ctx, &sub); err != nil {
		return result.Fail[subscription.Subscription](err)
	}

	return result.Success(sub)
}
