package command

import (
	"context"
	"time"

	"github.com/mistribe/subtracker/internal/domain/authorization"
	"github.com/mistribe/subtracker/internal/domain/provider"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/langext/option"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type UpdateProviderCommand struct {
	ProviderID     types.ProviderID
	Name           string
	Description    *string
	IconUrl        *string
	Url            *string
	PricingPageUrl *string
	Labels         []types.LabelID
	UpdatedAt      option.Option[time.Time]
}

type UpdateProviderCommandHandler struct {
	providerRepository ports.ProviderRepository
	labelRepository    ports.LabelRepository
	authorization      ports.Authorization
}

func NewUpdateProviderCommandHandler(
	providerRepository ports.ProviderRepository,
	labelRepository ports.LabelRepository,
	authorization ports.Authorization) *UpdateProviderCommandHandler {
	return &UpdateProviderCommandHandler{
		providerRepository: providerRepository,
		labelRepository:    labelRepository,
		authorization:      authorization,
	}
}

func (h *UpdateProviderCommandHandler) Handle(
	ctx context.Context,
	cmd UpdateProviderCommand) result.Result[provider.Provider] {
	prov, err := h.providerRepository.GetById(ctx, cmd.ProviderID)
	if err != nil {
		return result.Fail[provider.Provider](err)
	}

	if prov == nil {
		return result.Fail[provider.Provider](provider.ErrProviderNotFound)

	}

	if err = h.authorization.Can(ctx, authorization.PermissionWrite).For(prov); err != nil {
		return result.Fail[provider.Provider](err)
	}

	return h.update(ctx, cmd, prov)
}

func (h *UpdateProviderCommandHandler) update(
	ctx context.Context,
	cmd UpdateProviderCommand,
	prov provider.Provider) result.Result[provider.Provider] {

	if err := ensureLabelExists(ctx, h.labelRepository, cmd.Labels); err != nil {
		return result.Fail[provider.Provider](err)
	}
	var updatedAt time.Time
	if cmd.UpdatedAt == nil {
		updatedAt = time.Now()
	} else {
		updatedAt = cmd.UpdatedAt.ValueOrDefault(time.Now())
	}

	prov.SetName(cmd.Name)
	prov.SetDescription(cmd.Description)
	prov.SetIconUrl(cmd.IconUrl)
	prov.SetUrl(cmd.Url)
	prov.SetPricingPageUrl(cmd.PricingPageUrl)
	prov.SetLabels(cmd.Labels)
	prov.SetUpdatedAt(updatedAt)

	validationErrors := prov.GetValidationErrors()
	if validationErrors.HasErrors() {
		return result.Fail[provider.Provider](validationErrors)
	}

	if err := h.providerRepository.Save(ctx, prov); err != nil {
		return result.Fail[provider.Provider](err)
	}

	return result.Success(prov)
}
