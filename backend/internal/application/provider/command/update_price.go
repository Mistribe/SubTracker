package command

import (
	"context"
	"time"

	"github.com/google/uuid"
	"golang.org/x/text/currency"

	"github.com/mistribe/subtracker/internal/domain/provider"
	"github.com/mistribe/subtracker/pkg/ext"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type UpdatePriceCommand struct {
	ProviderId uuid.UUID
	PlanId     uuid.UUID

	Id        uuid.UUID
	StartDate time.Time
	EndDate   *time.Time
	Currency  currency.Unit
	Amount    float64
	UpdatedAt *time.Time
}

type UpdatePriceCommandHandler struct {
	providerRepository provider.Repository
}

func NewUpdatePriceCommandHandler(providerRepository provider.Repository) *UpdatePriceCommandHandler {
	return &UpdatePriceCommandHandler{providerRepository: providerRepository}
}

func (h UpdatePriceCommandHandler) Handle(ctx context.Context, cmd UpdatePriceCommand) result.Result[provider.Price] {
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

	price := pln.GetPriceById(cmd.Id)
	if price == nil {
		return result.Fail[provider.Price](provider.ErrPriceNotFound)
	}

	updatedAt := ext.ValueOrDefault(cmd.UpdatedAt, time.Now())

	price.SetStartDate(cmd.StartDate)
	price.SetEndDate(cmd.EndDate)
	price.SetCurrency(cmd.Currency)
	price.SetAmount(cmd.Amount)
	price.SetUpdatedAt(updatedAt)

	pln.Prices().Update(price)
	pln.SetUpdatedAt(updatedAt)
	prov.Plans().Update(pln)
	prov.SetUpdatedAt(updatedAt)

	if err := prov.GetValidationErrors(); err != nil {
		return result.Fail[provider.Price](err)
	}

	if err := h.providerRepository.Save(ctx, prov); err != nil {
		return result.Fail[provider.Price](err)
	}

	return result.Success(price)
}
