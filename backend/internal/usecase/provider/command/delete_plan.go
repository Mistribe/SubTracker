package command

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/domain/provider"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/langext/result"
	"github.com/mistribe/subtracker/pkg/x"
)

type DeletePlanCommand struct {
	ProviderId uuid.UUID
	Id         uuid.UUID

	DeletedAt *time.Time
}

type DeletePlanCommandHandler struct {
	providerRepository ports.ProviderRepository
}

func NewDeletePlanCommandHandler(providerRepository ports.ProviderRepository) *DeletePlanCommandHandler {
	return &DeletePlanCommandHandler{providerRepository: providerRepository}
}

func (h DeletePlanCommandHandler) Handle(ctx context.Context, cmd DeletePlanCommand) result.Result[bool] {
	prov, err := h.providerRepository.GetById(ctx, cmd.ProviderId)
	if err != nil {
		return result.Fail[bool](err)
	}

	if prov == nil {
		return result.Fail[bool](provider.ErrProviderNotFound)
	}

	if !prov.RemovePlanById(cmd.Id) {
		return result.Fail[bool](provider.ErrPlanNotFound)
	}

	deletedAt := x.ValueOrDefault(cmd.DeletedAt, time.Now())
	prov.SetUpdatedAt(deletedAt)

	if err := prov.GetValidationErrors(); err != nil {
		return result.Fail[bool](err)
	}

	ok, err := h.providerRepository.DeletePlan(ctx, cmd.Id)
	if err != nil {
		return result.Fail[bool](err)
	}
	if !ok {
		return result.Fail[bool](provider.ErrPlanNotFound)
	}

	if err := h.providerRepository.Save(ctx, prov); err != nil {
		return result.Fail[bool](err)
	}

	return result.Success(true)
}
