package command

import (
	"context"

	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/domain/user"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type DeleteSubscriptionCommand struct {
	Id uuid.UUID
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
	existingSub, err := h.subscriptionRepository.GetById(ctx, cmd.Id)
	if err != nil {
		return result.Fail[bool](err)
	}
	if existingSub == nil {
		return result.Fail[bool](family.ErrFamilyNotFound)
	}

	err = h.authorization.Can(ctx, user.PermissionDelete).For(existingSub)
	if err != nil {
		return result.Fail[bool](err)
	}

	ok, err := h.subscriptionRepository.Delete(ctx, cmd.Id)
	if err != nil {
		return result.Fail[bool](err)
	}

	return result.Success(ok)
}
