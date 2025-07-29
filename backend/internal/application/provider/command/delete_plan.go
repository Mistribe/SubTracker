package command

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/provider"
	"github.com/oleexo/subtracker/pkg/ext"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type DeletePlanCommand struct {
	ProviderId uuid.UUID
	Id         uuid.UUID

	DeletedAt *time.Time
}

type DeletePlanCommandHandler struct {
	providerRepository provider.Repository
}

func NewDeletePlanCommandHandler(providerRepository provider.Repository) *DeletePlanCommandHandler {
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

	deletedAt := ext.ValueOrDefault(cmd.DeletedAt, time.Now())
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
