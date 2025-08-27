package command

import (
	"context"

	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type DeleteSubscriptionCommand struct {
	Id uuid.UUID
}

type DeleteSubscriptionCommandHandler struct {
	subscriptionRepository ports.SubscriptionRepository
}

func NewDeleteSubscriptionCommandHandler(subscriptionRepository ports.SubscriptionRepository) *DeleteSubscriptionCommandHandler {
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
