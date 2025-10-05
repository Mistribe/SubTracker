package command

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/authorization"
	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type DeleteSubscriptionCommand struct {
	SubscriptionID types.SubscriptionID
}

type DeleteSubscriptionCommandHandler struct {
	subscriptionRepository ports.SubscriptionRepository
	authorization          ports.Authorization
}

func NewDeleteSubscriptionCommandHandler(subscriptionRepository ports.SubscriptionRepository,
	authorization ports.Authorization) *DeleteSubscriptionCommandHandler {
	return &DeleteSubscriptionCommandHandler{
		subscriptionRepository: subscriptionRepository,
		authorization:          authorization,
	}
}

func (h DeleteSubscriptionCommandHandler) Handle(
	ctx context.Context,
	cmd DeleteSubscriptionCommand) result.Result[bool] {
	existingSub, err := h.subscriptionRepository.GetById(ctx, cmd.SubscriptionID)
	if err != nil {
		return result.Fail[bool](err)
	}
	if existingSub == nil {
		return result.Fail[bool](family.ErrFamilyNotFound)
	}

	err = h.authorization.Can(ctx, authorization.PermissionDelete).For(existingSub)
	if err != nil {
		return result.Fail[bool](err)
	}

	ok, err := h.subscriptionRepository.Delete(ctx, cmd.SubscriptionID)
	if err != nil {
		return result.Fail[bool](err)
	}

	return result.Success(ok)
}
