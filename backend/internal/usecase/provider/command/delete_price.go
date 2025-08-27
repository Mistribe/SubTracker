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

type DeletePriceCommand struct {
	ProviderId uuid.UUID
	PlanId     uuid.UUID
	Id         uuid.UUID
	DeletedAt  *time.Time
}

type DeletePriceCommandHandler struct {
	providerRepository ports.ProviderRepository
}

func NewDeletePriceCommandHandler(providerRepository ports.ProviderRepository) *DeletePriceCommandHandler {
	return &DeletePriceCommandHandler{providerRepository: providerRepository}
}

func (h DeletePriceCommandHandler) Handle(ctx context.Context, cmd DeletePriceCommand) result.Result[bool] {
	prov, err := h.providerRepository.GetById(ctx, cmd.ProviderId)
	if err != nil {
		return result.Fail[bool](err)
	}

	if prov == nil {
		return result.Fail[bool](provider.ErrProviderNotFound)
	}

	pln := prov.GetPlanById(cmd.PlanId)
	if pln == nil {
		return result.Fail[bool](provider.ErrPlanNotFound)
	}

	if !pln.RemovePriceById(cmd.Id) {
		return result.Fail[bool](provider.ErrPriceNotFound)
	}

	deletedAt := x.ValueOrDefault(cmd.DeletedAt, time.Now())
	prov.SetUpdatedAt(deletedAt)
	pln.SetUpdatedAt(deletedAt)

	if err := prov.GetValidationErrors(); err != nil {
		return result.Fail[bool](err)
	}

	ok, err := h.providerRepository.DeletePrice(ctx, cmd.Id)
	if err != nil {
		return result.Fail[bool](err)
	}
	if !ok {
		return result.Fail[bool](provider.ErrPriceNotFound)
	}

	if err := h.providerRepository.Save(ctx, prov); err != nil {
		return result.Fail[bool](err)
	}

	return result.Success(true)
}
