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

type UpdatePlanCommand struct {
	ProviderId uuid.UUID

	Id          uuid.UUID
	Name        string
	Description *string
	UpdateAt    *time.Time
}

type UpdatePlanCommandHandler struct {
	providerRepository ports.ProviderRepository
}

func NewUpdatePlanCommandHandler(providerRepository ports.ProviderRepository) *UpdatePlanCommandHandler {
	return &UpdatePlanCommandHandler{providerRepository: providerRepository}
}

func (h UpdatePlanCommandHandler) Handle(ctx context.Context, cmd UpdatePlanCommand) result.Result[provider.Plan] {
	prov, err := h.providerRepository.GetById(ctx, cmd.ProviderId)
	if err != nil {
		return result.Fail[provider.Plan](err)
	}

	if prov == nil {
		return result.Fail[provider.Plan](provider.ErrProviderNotFound)
	}

	plan := prov.GetPlanById(cmd.Id)
	if plan == nil {
		return result.Fail[provider.Plan](provider.ErrPlanNotFound)
	}

	updatedAt := x.ValueOrDefault(cmd.UpdateAt, time.Now())

	plan.SetName(cmd.Name)
	plan.SetDescription(cmd.Description)
	plan.SetUpdatedAt(updatedAt)

	prov.SetUpdatedAt(updatedAt)

	if err := prov.GetValidationErrors(); err != nil {
		return result.Fail[provider.Plan](err)
	}

	if err := h.providerRepository.Save(ctx, prov); err != nil {
		return result.Fail[provider.Plan](err)
	}

	return result.Success(plan)
}
