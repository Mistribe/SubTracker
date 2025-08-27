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

type CreatePlanCommand struct {
	ProviderId uuid.UUID

	Id          *uuid.UUID
	Name        string
	Description *string
	CreatedAt   *time.Time
}

type CreatePlanCommandHandler struct {
	providerRepository ports.ProviderRepository
}

func NewCreatePlanCommandHandler(providerRepository ports.ProviderRepository) *CreatePlanCommandHandler {
	return &CreatePlanCommandHandler{providerRepository: providerRepository}
}

func (h CreatePlanCommandHandler) Handle(ctx context.Context, cmd CreatePlanCommand) result.Result[provider.Plan] {
	prov, err := h.providerRepository.GetById(ctx, cmd.ProviderId)
	if err != nil {
		return result.Fail[provider.Plan](err)
	}

	if prov == nil {
		return result.Fail[provider.Plan](provider.ErrProviderNotFound)
	}

	if cmd.Id != nil {
		if prov.ContainsPlan(*cmd.Id) {
			return result.Fail[provider.Plan](provider.ErrPlanAlreadyExists)
		}
	} else {
		newId, err := uuid.NewV7()
		if err != nil {
			return result.Fail[provider.Plan](err)
		}
		cmd.Id = &newId
	}
	createdAt := x.ValueOrDefault(cmd.CreatedAt, time.Now())

	pln := provider.NewPlan(
		*cmd.Id,
		cmd.Name,
		cmd.Description,
		nil,
		createdAt,
		createdAt,
	)

	if !prov.AddPlan(pln) {
		return result.Fail[provider.Plan](provider.ErrPlanAlreadyExists)
	}

	prov.SetUpdatedAt(createdAt)

	if err := prov.GetValidationErrors(); err != nil {
		return result.Fail[provider.Plan](err)
	}

	if err := h.providerRepository.Save(ctx, prov); err != nil {
		return result.Fail[provider.Plan](err)
	}

	return result.Success(pln)

}
