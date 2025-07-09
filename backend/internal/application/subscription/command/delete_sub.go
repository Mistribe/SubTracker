package command

import (
	"context"
	"github.com/google/uuid"
	"github.com/oleexo/subtracker/internal/application/core/option"
	"github.com/oleexo/subtracker/internal/application/core/result"
	"github.com/oleexo/subtracker/internal/domain/subscription"
)

type DeleteSubscriptionCommand struct {
	Id uuid.UUID
}

type DeleteSubscriptionCommandHandler struct {
	subscriptionRepository subscription.Repository
}

func NewDeleteSubscriptionCommandHandler(subscriptionRepository subscription.Repository) *DeleteSubscriptionCommandHandler {
	return &DeleteSubscriptionCommandHandler{subscriptionRepository: subscriptionRepository}
}

func (h DeleteSubscriptionCommandHandler) Handle(ctx context.Context, command DeleteSubscriptionCommand) result.Result[result.Unit] {
	subOpt, err := h.subscriptionRepository.Get(ctx, command.Id)
	if err != nil {
		return result.Fail[result.Unit](err)
	}
	return option.Match(subOpt, func(in subscription.Subscription) result.Result[result.Unit] {
		if err := h.subscriptionRepository.Delete(ctx, in.Id()); err != nil {
			return result.Fail[result.Unit](err)
		}
		return result.Void()
	}, func() result.Result[result.Unit] {
		return result.Fail[result.Unit](subscription.ErrSubscriptionNotFound)
	})

}
