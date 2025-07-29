package command

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/provider"
	"github.com/oleexo/subtracker/pkg/ext"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type DeletePriceCommand struct {
	ProviderId uuid.UUID
	PlanId     uuid.UUID
	Id         uuid.UUID
	DeletedAt  *time.Time
}

type DeletePriceCommandHandler struct {
	providerRepository provider.Repository
}

func NewDeletePriceCommandHandler(providerRepository provider.Repository) *DeletePriceCommandHandler {
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

	deletedAt := ext.ValueOrDefault(cmd.DeletedAt, time.Now())
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
