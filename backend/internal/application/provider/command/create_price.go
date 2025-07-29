package command

import (
	"context"
	"time"

	"github.com/google/uuid"
	"golang.org/x/text/currency"

	"github.com/oleexo/subtracker/internal/domain/provider"
	"github.com/oleexo/subtracker/pkg/ext"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type CreatePriceCommand struct {
	ProviderId uuid.UUID
	PlanId     uuid.UUID

	Id        *uuid.UUID
	StartDate time.Time
	EndDate   *time.Time
	Currency  currency.Unit
	Amount    float64
	CreatedAt *time.Time
}

type CreatePriceCommandHandler struct {
	providerRepository provider.Repository
}

func NewCreatePriceCommandHandler(providerRepository provider.Repository) *CreatePriceCommandHandler {
	return &CreatePriceCommandHandler{providerRepository: providerRepository}
}

func (h CreatePriceCommandHandler) Handle(ctx context.Context, cmd CreatePriceCommand) result.Result[provider.Price] {
	prov, err := h.providerRepository.GetById(ctx, cmd.ProviderId)
	if err != nil {
		return result.Fail[provider.Price](err)
	}

	if prov == nil {
		return result.Fail[provider.Price](provider.ErrProviderNotFound)
	}

	pln := prov.GetPlanById(cmd.PlanId)
	if pln == nil {
		return result.Fail[provider.Price](provider.ErrPlanNotFound)
	}

	if cmd.Id != nil {
		if pln.ContainsPrice(*cmd.Id) {
			return result.Fail[provider.Price](provider.ErrPriceAlreadyExists)
		}
	} else {
		newId, err := uuid.NewV7()
		if err != nil {
			return result.Fail[provider.Price](err)
		}
		cmd.Id = &newId
	}

	createdAt := ext.ValueOrDefault(cmd.CreatedAt, time.Now())

	prc := provider.NewPrice(
		*cmd.Id,
		cmd.StartDate,
		cmd.EndDate,
		cmd.Currency,
		cmd.Amount,
		createdAt,
		createdAt)

	if !pln.AddPrice(prc) {
		return result.Fail[provider.Price](provider.ErrPriceAlreadyExists)
	}

	pln.SetUpdatedAt(createdAt)
	prov.SetUpdatedAt(createdAt)

	if err := prov.GetValidationErrors(); err != nil {
		return result.Fail[provider.Price](err)
	}

	if err := h.providerRepository.Save(ctx, prov); err != nil {
		return result.Fail[provider.Price](err)
	}

	return result.Success(prc)
}
