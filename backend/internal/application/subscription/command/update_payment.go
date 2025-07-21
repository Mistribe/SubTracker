package command

import (
	"context"
	"time"

	"github.com/google/uuid"
	"golang.org/x/text/currency"

	"github.com/oleexo/subtracker/internal/domain/subscription"
	"github.com/oleexo/subtracker/pkg/langext/option"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type UpdatePaymentCommand struct {
	SubscriptionId uuid.UUID
	PaymentId      uuid.UUID
	Price          float64
	StartDate      time.Time
	EndDate        option.Option[time.Time]
	Months         int
	Currency       currency.Unit
	UpdatedAt      option.Option[time.Time]
}

type UpdatePaymentCommandHandler struct {
	repository subscription.Repository
}

func NewUpdatePaymentCommandHandler(repository subscription.Repository) *UpdatePaymentCommandHandler {
	return &UpdatePaymentCommandHandler{repository: repository}
}

func (h UpdatePaymentCommandHandler) Handle(
	ctx context.Context,
	command UpdatePaymentCommand) result.Result[subscription.Subscription] {
	subOpt, err := h.repository.Get(ctx, command.SubscriptionId)
	if err != nil {
		return result.Fail[subscription.Subscription](err)
	}
	return option.Match(subOpt, func(sub subscription.Subscription) result.Result[subscription.Subscription] {
		return option.Match(sub.GetPaymentById(command.PaymentId),
			func(payment subscription.Payment) result.Result[subscription.Subscription] {
				return h.updatePayment(ctx, command, sub, payment)
			}, func() result.Result[subscription.Subscription] {
				return result.Fail[subscription.Subscription](subscription.ErrPaymentNotFound)
			})
	}, func() result.Result[subscription.Subscription] {
		return result.Fail[subscription.Subscription](subscription.ErrSubscriptionNotFound)
	})
}

func (h UpdatePaymentCommandHandler) updatePayment(
	ctx context.Context,
	command UpdatePaymentCommand,
	sub subscription.Subscription,
	payment subscription.Payment) result.Result[subscription.Subscription] {
	payment.SetPrice(command.Price)
	payment.SetStartDate(command.StartDate)
	payment.SetEndDate(command.EndDate)
	payment.SetMonths(command.Months)
	payment.SetCurrency(command.Currency)

	command.UpdatedAt.IfSome(func(updatedAt time.Time) {
		payment.SetUpdatedAt(updatedAt)
	})
	command.UpdatedAt.IfNone(func() {
		payment.SetUpdatedAt(time.Now())
	})

	if err := payment.Validate(); err != nil {
		return result.Fail[subscription.Subscription](err)
	}

	sub.UpdatePayment(payment)

	if err := sub.Validate(); err != nil {
		return result.Fail[subscription.Subscription](err)
	}

	if err := h.repository.Save(ctx, &sub); err != nil {
		return result.Fail[subscription.Subscription](err)
	}

	return result.Success(sub)
}
