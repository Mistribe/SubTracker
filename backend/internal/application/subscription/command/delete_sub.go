package command

import (
	"context"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/subscription"
	"github.com/oleexo/subtracker/pkg/langext/result"
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

func (h DeleteSubscriptionCommandHandler) Handle(
	ctx context.Context,
	cmd DeleteSubscriptionCommand) result.Result[bool] {
	ok, err := h.subscriptionRepository.Delete(ctx, cmd.Id)
	if err != nil {
		return result.Fail[bool](err)
	}

	return result.Success(ok)
}
