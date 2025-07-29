package command

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/provider"
	"github.com/oleexo/subtracker/pkg/ext"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type UpdatePlanCommand struct {
	ProviderId uuid.UUID

	Id          uuid.UUID
	Name        string
	Description *string
	UpdateAt    *time.Time
}

type UpdatePlanCommandHandler struct {
	providerRepository provider.Repository
}

func NewUpdatePlanCommandHandler(providerRepository provider.Repository) *UpdatePlanCommandHandler {
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

	updatedAt := ext.ValueOrDefault(cmd.UpdateAt, time.Now())

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
