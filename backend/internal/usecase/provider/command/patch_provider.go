package command

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/provider"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type PatchProviderCommand struct {
	Provider provider.Provider
}

type PatchProviderCommandHandler struct {
	providerRepository ports.ProviderRepository
	labelRepository    ports.LabelRepository
}

func NewPatchProviderCommandHandler(providerRepository ports.ProviderRepository,
	labelRepository ports.LabelRepository) *PatchProviderCommandHandler {
	return &PatchProviderCommandHandler{
		providerRepository: providerRepository,
		labelRepository:    labelRepository,
	}
}

func (h PatchProviderCommandHandler) Handle(ctx context.Context,
	cmd PatchProviderCommand) result.Result[provider.Provider] {
	existingProvider, err := h.providerRepository.GetById(ctx, cmd.Provider.Id())
	if err != nil {
		return result.Fail[provider.Provider](err)
	}
	if existingProvider == nil {
		return h.create(ctx, cmd.Provider)
	} else {
		return h.update(ctx, cmd.Provider, existingProvider)
	}
}

func (h PatchProviderCommandHandler) create(ctx context.Context,
	newProvider provider.Provider) result.Result[provider.Provider] {

	if err := ensureLabelExists(ctx, h.labelRepository, newProvider.Labels().Values()); err != nil {
		return result.Fail[provider.Provider](err)
	}

	if err := newProvider.GetValidationErrors(); err != nil {
		return result.Fail[provider.Provider](err)
	}

	if err := h.providerRepository.Save(ctx, newProvider); err != nil {
		return result.Fail[provider.Provider](err)
	}

	return result.Success(newProvider)
}

func (h PatchProviderCommandHandler) update(ctx context.Context, updateInfo provider.Provider,
	existingProvider provider.Provider) result.Result[provider.Provider] {
	if err := ensureLabelExists(ctx, h.labelRepository, updateInfo.Labels().Values()); err != nil {
		return result.Fail[provider.Provider](err)
	}

	existingProvider.SetName(updateInfo.Name())
	existingProvider.SetDescription(updateInfo.Description())
	existingProvider.SetIconUrl(updateInfo.IconUrl())
	existingProvider.SetUrl(updateInfo.Url())
	existingProvider.SetPricingPageUrl(updateInfo.PricingPageUrl())
	existingProvider.SetLabels(updateInfo.Labels().Values())
	existingProvider.SetOwner(updateInfo.Owner())
	h.updatePlans(updateInfo, existingProvider)

	if err := existingProvider.GetValidationErrors(); err != nil {
		return result.Fail[provider.Provider](err)
	}

	if err := h.providerRepository.Save(ctx, existingProvider); err != nil {
		return result.Fail[provider.Provider](err)
	}

	return result.Success(existingProvider)
}

func (h PatchProviderCommandHandler) updatePlans(info provider.Provider,
	existingProvider provider.Provider) {
	for pln := range existingProvider.Plans().It() {
		if !info.Plans().Contains(pln) {
			existingProvider.Plans().Remove(pln)
		}
	}

	for pln := range info.Plans().It() {
		existingPlan, ok := existingProvider.Plans().Get(pln)
		if !ok {
			existingProvider.Plans().Add(pln)
		} else {
			existingPlan.SetName(pln.Name())
			existingPlan.SetDescription(pln.Description())
			existingPlan.SetUpdatedAt(pln.UpdatedAt())
			existingProvider.Plans().Update(existingPlan)
			h.updatePrices(pln, existingPlan)
		}
	}
}

func (h PatchProviderCommandHandler) updatePrices(info provider.Plan, existingPlan provider.Plan) {
	for pr := range existingPlan.Prices().It() {
		if !info.ContainsPrice(pr.Id()) {
			existingPlan.RemovePriceById(pr.Id())
		}
	}

	for pr := range info.Prices().It() {
		existingPrice := existingPlan.GetPriceById(pr.Id())
		if existingPrice == nil {
			existingPlan.AddPrice(pr)
		} else {
			existingPrice.SetStartDate(pr.StartDate())
			existingPrice.SetEndDate(pr.EndDate())
			existingPrice.SetAmount(pr.Amount())
			existingPrice.SetCurrency(pr.Currency())
			existingPrice.SetUpdatedAt(pr.UpdatedAt())
			info.Prices().Update(existingPrice)
		}
	}
}
