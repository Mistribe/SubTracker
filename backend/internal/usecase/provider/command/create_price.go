package command

import (
	"context"
	"time"

	"github.com/google/uuid"
	"golang.org/x/text/currency"

	"github.com/mistribe/subtracker/internal/domain/provider"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/langext/result"
	"github.com/mistribe/subtracker/pkg/x"
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
	providerRepository ports.ProviderRepository
}

func NewCreatePriceCommandHandler(providerRepository ports.ProviderRepository) *CreatePriceCommandHandler {
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

	createdAt := x.ValueOrDefault(cmd.CreatedAt, time.Now())

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

	prov.Plans().Update(pln)
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
