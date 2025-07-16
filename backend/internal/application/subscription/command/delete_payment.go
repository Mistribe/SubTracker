package command

import (
	"context"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/subscription"
	"github.com/oleexo/subtracker/pkg/langext/option"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type DeletePaymentCommand struct {
	SubscriptionId uuid.UUID
	PaymentId      uuid.UUID
}

type DeletePaymentCommandHandler struct {
	repository subscription.Repository
}

func NewDeletePaymentCommandHandler(repository subscription.Repository) *DeletePaymentCommandHandler {
	return &DeletePaymentCommandHandler{repository: repository}
}

func (h DeletePaymentCommandHandler) Handle(
	ctx context.Context,
	command DeletePaymentCommand) result.Result[result.Unit] {
	subOpt, err := h.repository.Get(ctx, command.SubscriptionId)
	if err != nil {
		return result.Fail[result.Unit](err)
	}

	return option.Match(subOpt, func(sub subscription.Subscription) result.Result[result.Unit] {
		return h.deletePayment(ctx, command, sub)
	}, func() result.Result[result.Unit] {
		return result.Fail[result.Unit](subscription.ErrSubscriptionNotFound)
	})
}

func (h DeletePaymentCommandHandler) deletePayment(
	ctx context.Context, command DeletePaymentCommand,
	sub subscription.Subscription) result.Result[result.Unit] {
	sub.RemovePayment(command.PaymentId)

	if err := h.repository.Save(ctx, &sub); err != nil {
		return result.Fail[result.Unit](err)
	}

	return result.Void()
}
